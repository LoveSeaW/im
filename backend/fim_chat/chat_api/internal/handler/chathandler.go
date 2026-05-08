package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fim_server/common/models/ctype"
	"fim_server/common/response"
	"fim_server/common/service/redis_service"
	"fim_server/fim_chat/chat_api/internal/svc"
	"fim_server/fim_chat/chat_api/internal/types"
	"fim_server/fim_chat/chat_models"
	"fim_server/fim_file/file_rpc/files"
	"fim_server/fim_user/user_models"
	"fim_server/fim_user/user_rpc/types/user_rpc"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// 连接全局配置
const (
	maxConnections    = 50000                  // 最大连接数
	messageQueueSize  = 100                    // 单用户消息队列大小
	writeTimeout      = 100 * time.Millisecond // 写超时
	heartbeatInterval = 30 * time.Second       // 心跳间隔
)

var (
	activeConnections int32
)

// 改进的连接结构
type UserWsInfo struct {
	UserInfo    user_models.UserModel
	MsgChan     chan []byte         // 带缓冲的消息通道
	WsClientMap *sync.Map           // 并发安全的连接映射
	CurrentConn *websocket.Conn     // 最新连接
	Limiter     *limit.TokenLimiter // 限流器(100qps + 突发10)
}

// 全局并发安全存储
var (
	UserOnlineWsMap sync.Map // key: uint(userID), value: *UserWsInfo
	VideoCallMap    sync.Map // key: "user1_user2", value: time.Time
)

func chatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req types.ChatRequest
		if err := httpx.ParseHeaders(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 1. 连接数检查
		if atomic.LoadInt32(&activeConnections) >= maxConnections {
			http.Error(w, "server overload", http.StatusServiceUnavailable)
			return
		}
		atomic.AddInt32(&activeConnections, 1)
		defer atomic.AddInt32(&activeConnections, -1)

		// 2. 协议升级
		upgrader := websocket.Upgrader{
			HandshakeTimeout: 5 * time.Second,
			CheckOrigin:      func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Errorf("WebSocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		// 3. 获取用户信息
		// 调用户服务，获取当前用户信息
		res, err := svcCtx.UserRpc.UserInfo(context.Background(), &user_rpc.UserInfoRequest{
			UserId: uint32(req.UserID),
		})
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		var userInfo user_models.UserModel
		err = json.Unmarshal(res.Data, &userInfo)
		if err != nil {
			SendTipErrMsg(conn, "用户信息获取失败")
			return
		}

		// 4. 连接管理
		addr := conn.RemoteAddr().String()
		userWsInfo := registerConnection(req.UserID, addr, conn, svcCtx)
		//userWsInfo := initUserWsInfo(userInfo, conn, addr)
		defer cleanupConnection(req.UserID, addr, svcCtx)

		// 5. 启动消息处理协程
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go handleIncomingMessages(ctx, conn, userWsInfo, svcCtx, req.UserID, userInfo)
		go handleOutgoingMessages(ctx, conn, userWsInfo)

		// 6. 心跳检测
		heartbeatTicker := time.NewTicker(heartbeatInterval)
		defer heartbeatTicker.Stop()

		for {
			select {
			case <-heartbeatTicker.C:
				if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeTimeout)); err != nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

// 注册连接（集成go-zero限流器）
func registerConnection(userID uint, addr string, conn *websocket.Conn, svcCtx *svc.ServiceContext) *UserWsInfo {
	redisConf := redis.RedisConf{
		Host: svcCtx.Config.Redis.Addr,
		Pass: svcCtx.Config.Redis.Pwd,
	}

	newRedis, err := redis.NewRedis(redisConf)
	if err != nil {
		return nil
	}
	// 每个用户独立的限流器（每秒100请求，突发10）
	limiter := limit.NewTokenLimiter(1, 100, newRedis, fmt.Sprintf("ws:limit:%d", userID))

	userWs := &UserWsInfo{
		CurrentConn: conn,
		Limiter:     limiter,
	}

	// 加载或存储用户连接信息
	actual, loaded := UserOnlineWsMap.LoadOrStore(userID, userWs)
	if loaded {
		userWs = actual.(*UserWsInfo)
	}

	userWs.WsClientMap.Store(addr, conn)
	svcCtx.Redis.HSet("online_users", fmt.Sprintf("%d", userID), time.Now().Unix())
	return userWs
}

// 清理连接资源
func cleanupConnection(userID uint, addr string, svcCtx *svc.ServiceContext) {
	if val, ok := UserOnlineWsMap.Load(userID); ok {
		userWsInfo := val.(*UserWsInfo)
		userWsInfo.WsClientMap.Delete(addr)

		// 检查是否无活跃连接
		count := 0
		userWsInfo.WsClientMap.Range(func(_, _ interface{}) bool {
			count++
			return true
		})

		if count == 0 {
			UserOnlineWsMap.Delete(userID)
			svcCtx.Redis.HDel("online", fmt.Sprintf("%d", userID))
		}
	}
}

// 处理入站消息
func handleIncomingMessages(ctx context.Context, conn *websocket.Conn, userWsInfo *UserWsInfo, svcCtx *svc.ServiceContext, userID uint, userInfo user_models.UserModel) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, p, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					logx.Errorf("Read error: %v", err)
				}
				return
			}

			// 限流检查
			if !userWsInfo.Limiter.Allow() {
				SendTipErrMsg(conn, "消息发送频率过高")
				continue
			}

			var request ChatRequest
			if err := json.Unmarshal(p, &request); err != nil {
				SendTipErrMsg(conn, "消息格式错误")
				continue
			}

			// 消息处理逻辑（原业务代码）
			if err := processMessage(request, userID, userWsInfo, svcCtx, userInfo); err != nil {
				SendTipErrMsg(conn, err.Error())
			}
		}
	}
}

func processMessage(req ChatRequest, userID uint, userWsInfo *UserWsInfo, svcCtx *svc.ServiceContext, userInfo user_models.UserModel) error {
	// 遍历在线的用户， 和当前这个人是好友的，就给他发送好友在线

	// 先把所有在线的用户id取出来，以及待确认的用户id，然后传到用户rpc服务中
	// [1,2,3]  3
	// 在rpc服务中，去判断哪些用户是好友关系

	//if userInfo.UserConfModel.FriendOnline {
	// 如果用户开启了好友上线提醒
	// 查一下自己的好友是不是上线了
	friendRes, err := svcCtx.UserRpc.FriendList(context.Background(), &user_rpc.FriendListRequest{
		User: uint32(userID),
	})
	// 3 [3,4,5]
	if err != nil {
		logx.Error(err)
		return err
	}
	logx.Infof("用户上线：%s 用户id: %d", userInfo.Nickname, userID)

	for _, info := range friendRes.FriendList {
		// 因为自己可以加自己为好友，自己上线没必要给自己发
		if uint(info.UserId) == userID {
			continue
		}
		friendAny, ok := UserOnlineWsMap.Load(info.UserId)
		friend := friendAny.(*UserWsInfo)
		if ok {
			text := fmt.Sprintf("好友 %s 上线了", userInfo.Nickname)
			// 判断用户是否开了好友上线提醒
			if friend.UserInfo.UserConfModel.FriendOnline {
				// 好友上线了
				//friend.Conn.WriteMessage(websocket.TextMessage, []byte(text))

				resp := ChatResponse{
					Msg: ctype.Msg{
						Type: ctype.FriendOnlineMsgType,
						FriendOnlineMsg: &ctype.FriendOnlineMsg{
							Nickname: userInfo.Nickname,
							Avatar:   userInfo.Avatar,
							Content:  text,
							FriendID: userInfo.ID,
						},
					},
					CreatedAt: time.Now(),
				}
				byteData, _ := json.Marshal(resp)
				sendWsMapMsg(friend.WsClientMap, byteData)
			}
		}
	}
	// 查一下自己的好友列表，返回用户id列表，看看在不在这个UserWsMap中，在的话，就给自己发个好友上线的消息
	conn := userWsInfo.CurrentConn
	//}
	for {
		// 消息类型，消息，错误
		_, p, err1 := conn.ReadMessage()
		if err1 != nil {
			// 用户断开聊天
			fmt.Println(err1)
			break
		}
		// 目前这里做不到实时更新
		// 要做到实时更新，把用户的这些配置放到缓存里面去
		// 用户聊天之前就向缓存里面去拿用户的相关配置信息 拿不到的情况下，去调用户rpc方法，然后缓存到缓存里面
		// 在后台，把用户的配置更新之后，让这条缓存失效即可
		if userInfo.UserConfModel.CurtailChat {
			SendTipErrMsg(conn, "当前用户被限制聊天")
			continue
		}

		var request ChatRequest
		err2 := json.Unmarshal(p, &request)
		if err2 != nil {
			// 用户乱发消息
			logx.Error(err2)
			SendTipErrMsg(conn, "参数解析失败")
			continue
		}
		if request.RevUserID != userID {
			// 判断你聊天的这个人是不是你的好友
			isFriendRes, err := svcCtx.UserRpc.IsFriend(context.Background(), &user_rpc.IsFriendRequest{
				User1: uint32(userID),
				User2: uint32(request.RevUserID),
			})
			if err != nil {
				// 用户乱发消息
				logx.Error(err2)
				SendTipErrMsg(conn, "用户服务错误")
				continue
			}

			if !isFriendRes.IsFriend {
				SendTipErrMsg(conn, "你们还不是好友呢")
				continue
			}
		}
		// 判断type  1 - 14
		if !(request.Msg.Type >= 1 && request.Msg.Type <= 14) {
			SendTipErrMsg(conn, "消息类型错误")
			continue
		}

		// 校验消息
		msgValidateErr := request.Msg.Validate()
		if msgValidateErr != nil {
			SendTipErrMsg(conn, msgValidateErr.Error())
			continue
		}

		// 判断是否是文件类型
		switch request.Msg.Type {
		case ctype.TextMsgType:

		case ctype.FileMsgType:
			// 如果是文件类型，那么就要去请求文件rpc服务
			nameList := strings.Split(request.Msg.FileMsg.Src, "/")
			if len(nameList) == 0 {
				SendTipErrMsg(conn, "请上传文件")
				continue
			}
			fileID := nameList[len(nameList)-1]
			fileResponse, err3 := svcCtx.FileRpc.FileInfo(context.Background(), &files.FileInfoRequest{
				FileId: fileID,
			})
			if err3 != nil {
				logx.Error(err3)
				SendTipErrMsg(conn, err3.Error())
				continue
			}
			request.Msg.FileMsg.Title = fileResponse.FileName
			request.Msg.FileMsg.Size = fileResponse.FileSize
			request.Msg.FileMsg.Type = fileResponse.FileType
		case ctype.WithdrawMsgType:
			// 撤回消息的消息id是必填的
			if request.Msg.WithdrawMsg == nil {
				SendTipErrMsg(conn, "撤回消息id必填")
				continue
			}
			if request.Msg.WithdrawMsg.MsgID == 0 {
				SendTipErrMsg(conn, "撤回消息id必填")
				continue
			}

			// 自己只能撤回自己的
			// 找这个消息是谁发的
			var msgModel chat_models.ChatModel
			err = svcCtx.DB.Take(&msgModel, request.Msg.WithdrawMsg.MsgID).Error
			if err != nil {
				SendTipErrMsg(conn, "消息不存在")
				continue
			}

			// 已经是撤回消息的，不能再撤回了
			if msgModel.MsgType == ctype.WithdrawMsgType {
				SendTipErrMsg(conn, "撤回消息不能再撤回了")
				continue
			}

			// 判断是不是自己发的
			if msgModel.SendUserID != userID {
				SendTipErrMsg(conn, "只能撤回自己的消息")
				continue
			}

			// 判断消息的时间，小于2分钟的才能撤回
			now := time.Now()
			subTime := now.Sub(msgModel.CreatedAt)
			if subTime >= time.Minute*2 {
				SendTipErrMsg(conn, "只能撤回两分钟以内的消息哦~")
				continue
			}
			// 撤回逻辑
			// 收到撤回请求之后，服务端这边把原消息类型修改为撤回消息类型，并且记录原消息
			// 然后通知前端的收发双方，重新拉取聊天记录

			var content = "撤回了一条消息"
			if userInfo.UserConfModel.RecallMessage != nil {
				content = "撤回了一条消息," + *userInfo.UserConfModel.RecallMessage
			}
			// 前端可以判断，这个消息如果不是isMe，就可以把你替换成对方的昵称

			originMsg := msgModel.Msg
			originMsg.WithdrawMsg = nil // 这里可能会出现循环引用，所以拷贝了这个值，并且把撤回消息置空了

			svcCtx.DB.Model(&msgModel).Updates(chat_models.ChatModel{
				MsgPreview: "[撤回消息] - " + content,
				MsgType:    ctype.WithdrawMsgType,
				Msg: ctype.Msg{
					Type: ctype.WithdrawMsgType,
					WithdrawMsg: &ctype.WithdrawMsg{
						Content:   content,
						MsgID:     request.Msg.WithdrawMsg.MsgID,
						OriginMsg: &originMsg,
					},
				},
			})
		case ctype.ReplyMsgType:
			// 回复消息
			// 先校验
			if request.Msg.ReplyMsg == nil || request.Msg.ReplyMsg.MsgID == 0 {
				SendTipErrMsg(conn, "回复消息id必填")
				continue
			}

			// 找这个原消息
			var msgModel chat_models.ChatModel
			err = svcCtx.DB.Take(&msgModel, request.Msg.ReplyMsg.MsgID).Error
			if err != nil {
				SendTipErrMsg(conn, "消息不存在")
				continue
			}

			// 不能回复撤回消息
			if msgModel.MsgType == ctype.WithdrawMsgType {
				SendTipErrMsg(conn, "该消息已撤回")
				continue
			}

			// 回复的这个消息，必须是你自己或者当前和你聊天这个人发出来的

			// 原消息必须是 当前你要和对方聊的  原消息就会有一个 发送人id和接收人id，  我们聊天也会有一个发送人id和接收人id
			// 因为回复消息可以回复自己的，也可以回复别人的
			// 如果回复只能回复别人的？那么条件怎么写?
			if !((msgModel.SendUserID == userID && msgModel.RevUserID == request.RevUserID) ||
				(msgModel.SendUserID == request.RevUserID && msgModel.RevUserID == userID)) {
				SendTipErrMsg(conn, "只能回复自己或者对方的消息")
				continue
			}

			userBaseInfo, err5 := redis_service.GetUserBaseInfo(svcCtx.Redis, svcCtx.UserRpc, msgModel.SendUserID)
			if err5 != nil {
				logx.Error(err5)
				SendTipErrMsg(conn, err5.Error())
				continue
			}

			request.Msg.ReplyMsg.Msg = &msgModel.Msg
			request.Msg.ReplyMsg.UserID = msgModel.SendUserID
			request.Msg.ReplyMsg.UserNickName = userBaseInfo.NickName
			request.Msg.ReplyMsg.OriginMsgDate = msgModel.CreatedAt
			request.Msg.ReplyMsg.ReplyMsgPreview = msgModel.MsgPreviewMethod()

		case ctype.QuoteMsgType:
			// 回复消息
			// 先校验
			if request.Msg.QuoteMsg == nil || request.Msg.QuoteMsg.MsgID == 0 {
				SendTipErrMsg(conn, "引用消息id必填")
				continue
			}

			// 找这个原消息
			var msgModel chat_models.ChatModel
			err = svcCtx.DB.Take(&msgModel, request.Msg.QuoteMsg.MsgID).Error
			if err != nil {
				SendTipErrMsg(conn, "消息不存在")
				continue
			}

			// 不能回复撤回消息
			if msgModel.MsgType == ctype.WithdrawMsgType {
				SendTipErrMsg(conn, "该消息已撤回")
				continue
			}

			// 回复的这个消息，必须是你自己或者当前和你聊天这个人发出来的

			if !((msgModel.SendUserID == userID && msgModel.RevUserID == request.RevUserID) ||
				(msgModel.SendUserID == request.RevUserID && msgModel.RevUserID == userID)) {
				SendTipErrMsg(conn, "只能回复自己或者对方的消息")
				continue
			}

			userBaseInfo, err5 := redis_service.GetUserBaseInfo(svcCtx.Redis, svcCtx.UserRpc, msgModel.SendUserID)
			if err5 != nil {
				logx.Error(err5)
				SendTipErrMsg(conn, err5.Error())
				continue
			}

			request.Msg.QuoteMsg.Msg = &msgModel.Msg
			request.Msg.QuoteMsg.UserID = msgModel.SendUserID
			request.Msg.QuoteMsg.UserNickName = userBaseInfo.NickName
			request.Msg.QuoteMsg.OriginMsgDate = msgModel.CreatedAt
			request.Msg.QuoteMsg.QuoteMsgPreview = msgModel.MsgPreviewMethod()
		case ctype.VideoCallMsgType:
			data := request.Msg.VideoCallMsg
			// 先判断对方是否在线
			_, ok2 := UserOnlineWsMap.Load(request.RevUserID)
			if !ok2 {
				SendTipErrMsg(conn, "对方不在线")
				continue
			}

			key := fmt.Sprintf("%d_%d", userInfo.ID, request.RevUserID)

			switch data.Flag {
			case 0:
				// 给自己的页面展示一个等待对方接听的一个弹框
				conn.WriteJSON(ChatResponse{
					Msg: ctype.Msg{
						Type: ctype.VideoCallMsgType,
						VideoCallMsg: &ctype.VideoCallMsg{
							Flag: 1,
						},
					},
				})
				// 给对方的页面展示一个等待接听的一个弹框
				sendRevUserMsg(request.RevUserID, userID, ctype.Msg{
					Type: ctype.VideoCallMsgType,
					VideoCallMsg: &ctype.VideoCallMsg{
						Flag: 2,
					},
				})
			case 1: // 自己挂断
				sendRevUserMsg(request.RevUserID, userID, ctype.Msg{
					Type: ctype.VideoCallMsgType,
					VideoCallMsg: &ctype.VideoCallMsg{
						Flag: 3,
						Msg:  "发起者已挂断",
					},
				})
			case 2: // 对方挂断
				// 对方点击挂断，那么它的目标就是revUserID，也就是上面的conn
				sendRevUserMsg(request.RevUserID, userID, ctype.Msg{
					Type: ctype.VideoCallMsgType,
					VideoCallMsg: &ctype.VideoCallMsg{
						Flag: 4,
						Msg:  "用户拒绝了你的视频通话",
					},
				})
			case 3: // 对方接受
				// 让发送者准备去发offer
				sendRevUserMsg(request.RevUserID, userID, ctype.Msg{
					Type: ctype.VideoCallMsgType,
					VideoCallMsg: &ctype.VideoCallMsg{
						Flag: 5, // 让发送者准备去发offer
						Type: "create_offer",
					},
				})
			case 4: //我方正常挂断
				// 算你们的通话时长
				// 从发offer开始，算一个开始时间，到这里算一个结束时间，就是视频通话的时间
				startTimeAny, ok3 := VideoCallMap.Load(key)
				startTime := startTimeAny.(time.Time)
				var sendUserID = userID
				var revUserID = request.RevUserID
				fmt.Println("key1", key, sendUserID, revUserID)
				var endReason int8
				if !ok3 {
					// 先按照我方挂断
					// 1 -> 2  1-2
					// 2-1 1-2
					key = fmt.Sprintf("%d_%d", request.RevUserID, userInfo.ID)
					_startTimeAny, ok4 := VideoCallMap.Load(key)
					_startTime := _startTimeAny.(time.Time)
					if !ok4 {
						SendTipErrMsg(conn, "消息起始时间错误")
						continue
					}
					sendUserID = request.RevUserID
					revUserID = userInfo.ID
					endReason = 1 // 接收方挂断
					startTime = _startTime
				}
				fmt.Println("key2", key, sendUserID, revUserID)
				subTime := time.Now().Sub(startTime)
				fmt.Printf("用户正常挂断， 视频通话时长为 %s\n", subTime)
				request.Msg.VideoCallMsg.StartTime = startTime
				request.Msg.VideoCallMsg.EndTime = time.Now()
				request.Msg.VideoCallMsg.EndReason = endReason
				sendRevUserMsg(request.RevUserID, userID, ctype.Msg{
					Type: ctype.VideoCallMsgType,
					VideoCallMsg: &ctype.VideoCallMsg{
						Flag: 6, // 对方挂断了
					},
				})
				msgID := InsertMsgByChat(svcCtx.DB, revUserID, sendUserID, request.Msg)
				// 看看目标用户在不在线  给发送双方都要发消息
				SendMsgByUser(svcCtx, revUserID, sendUserID, request.Msg, msgID)

				VideoCallMap.Delete(key)

				continue

			case 5: // 对方挂断
				key = fmt.Sprintf("%d_%d", request.RevUserID, userInfo.ID)
				startTimeAny, ok3 := VideoCallMap.Load(key)
				startTime := startTimeAny.(time.Time)
				if !ok3 {
					fmt.Println("没有获取到起始时间")
					continue
				}
				subTime := time.Now().Sub(startTime)
				fmt.Printf("对方正常挂断， 视频通话时长为 %s\n", subTime)
			}

			switch data.Type {
			case "offer": // offer
				sendRevUserMsg(request.RevUserID, userID, ctype.Msg{
					Type: ctype.VideoCallMsgType,
					VideoCallMsg: &ctype.VideoCallMsg{
						Type: "offer",
						Data: data.Data,
					},
				})
				VideoCallMap.Store(key, time.Now())
				fmt.Println("offer", key)
			case "answer": // 应答
				sendRevUserMsg(request.RevUserID, userID, ctype.Msg{
					Type: ctype.VideoCallMsgType,
					VideoCallMsg: &ctype.VideoCallMsg{
						Type: "answer",
						Data: data.Data,
					},
				})
			case "offer_ice":
				sendRevUserMsg(request.RevUserID, userID, ctype.Msg{
					Type: ctype.VideoCallMsgType,
					VideoCallMsg: &ctype.VideoCallMsg{
						Type: "offer_ice",
						Data: data.Data,
					},
				})
			case "answer_ice":
				sendRevUserMsg(request.RevUserID, userID, ctype.Msg{
					Type: ctype.VideoCallMsgType,
					VideoCallMsg: &ctype.VideoCallMsg{
						Type: "answer_ice",
						Data: data.Data,
					},
				})
			}
			// 自己这方可以挂断

			// 对方也可以挂断

			// 如果对方开了多个浏览器，只用找其中的一个，找第一个
			continue
		}

		// 先入库
		msgID := InsertMsgByChat(svcCtx.DB, request.RevUserID, userID, request.Msg)
		// 看看目标用户在不在线  给发送双方都要发消息
		SendMsgByUser(svcCtx, request.RevUserID, userID, request.Msg, msgID)

	}
	return nil
}

// 处理出站消息（异步非阻塞）
func handleOutgoingMessages(ctx context.Context, conn *websocket.Conn, userWsInfo *UserWsInfo) {
	ticker := time.NewTicker(50 * time.Millisecond) // 批量发送间隔
	defer ticker.Stop()

	var buffer [][]byte
	for {
		select {
		case msg := <-userWsInfo.MsgChan:
			buffer = append(buffer, msg)
			if len(buffer) >= 20 { // 达到批量大小立即发送
				sendBatch(conn, buffer)
				buffer = nil
			}

		case <-ticker.C:
			if len(buffer) > 0 {
				sendBatch(conn, buffer)
				buffer = nil
			}

		case <-ctx.Done():
			return
		}
	}
}

// 批量发送优化
func sendBatch(conn *websocket.Conn, messages [][]byte) {
	conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err := conn.WriteMessage(websocket.TextMessage, bytes.Join(messages, []byte("\n"))); err != nil {
		logx.Errorf("Batch write failed: %v", err)
	}
}

// 改进的消息发送方法
func sendRevUserMsg(revUserID uint, sendUserID uint, msg ctype.Msg) {
	if val, ok := UserOnlineWsMap.Load(revUserID); ok {
		userWsInfo := val.(*UserWsInfo)
		resp := buildChatResponse(revUserID, sendUserID, msg)

		// 异步非阻塞发送
		select {
		case userWsInfo.MsgChan <- resp:
		default:
			logx.Errorf("User %d message queue full", revUserID)
		}
	}
}

func buildChatResponse(revUserID uint, sendUserID uint, msg ctype.Msg) []byte {
	userResAny, ok := UserOnlineWsMap.Load(revUserID)
	if !ok {
		return nil
	}
	sendUserAny, ok1 := UserOnlineWsMap.Load(sendUserID)
	sendUser := sendUserAny.(*UserWsInfo)
	var sendUserInfo ctype.UserInfo
	if ok1 {
		sendUserInfo = ctype.UserInfo{
			ID:       sendUser.UserInfo.ID,
			NickName: sendUser.UserInfo.Nickname,
			Avatar:   sendUser.UserInfo.Avatar,
		}
	}
	userRes := userResAny.(*UserWsInfo)
	respChat := ChatResponse{
		SendUser: sendUserInfo,
		RevUser: ctype.UserInfo{
			ID:       userRes.UserInfo.ID,
			NickName: userRes.UserInfo.Nickname,
			Avatar:   userRes.UserInfo.Avatar,
		},
		MsgPreview: msg.MsgPreview(),
		Msg:        msg,
		CreatedAt:  time.Now(),
	}

	data, _ := json.Marshal(respChat)
	return data
}

type ChatRequest struct {
	RevUserID uint      `json:"revUserID"` // 给谁发
	Msg       ctype.Msg `json:"msg"`
}

// 给一组的ws对象发消息
func sendWsMapMsg(wsMap *sync.Map, byteData []byte) {
	wsMap.Range(func(key, value interface{}) bool {
		conn := value.(*websocket.Conn)
		conn.WriteMessage(websocket.TextMessage, byteData)
		return true
	})
}

type ChatResponse struct {
	ID         uint           `json:"id"`
	IsMe       bool           `json:"isMe"`
	RevUser    ctype.UserInfo `json:"revUser"`
	SendUser   ctype.UserInfo `json:"sendUser"`
	Msg        ctype.Msg      `json:"msg"`
	CreatedAt  time.Time      `json:"created_at"`
	MsgPreview string         `json:"msgPreview"`
}

// SendTipErrMsg 发送错误提示的消息
func SendTipErrMsg(conn *websocket.Conn, msg string) {
	resp := ChatResponse{
		Msg: ctype.Msg{
			Type: ctype.TipMsgType,
			TipMsg: &ctype.TipMsg{
				Status:  "error",
				Content: msg,
			},
		},
		CreatedAt: time.Now(),
	}
	byteData, _ := json.Marshal(resp)
	conn.WriteMessage(websocket.TextMessage, byteData)

}

// InsertMsgByChat 消息入库
func InsertMsgByChat(db *gorm.DB, revUserId uint, sendUserID uint, msg ctype.Msg) (msgID uint) {
	switch msg.Type {
	case ctype.WithdrawMsgType:
		fmt.Println("撤回消息自己是不入库的")
		return
	}
	chatModel := chat_models.ChatModel{
		SendUserID: sendUserID,
		RevUserID:  revUserId,
		MsgType:    msg.Type,
		Msg:        msg,
	}
	chatModel.MsgPreview = chatModel.MsgPreviewMethod()
	err := db.Create(&chatModel).Error
	if err != nil {
		logx.Error(err)
		sendUserAny, ok := UserOnlineWsMap.Load(sendUserID)
		if !ok {
			return
		}
		sendUser := sendUserAny.(*UserWsInfo)
		SendTipErrMsg(sendUser.CurrentConn, "消息保存失败")
	}
	return chatModel.ID
}

// SendMsgByUser 发消息 给谁发 谁发的
func SendMsgByUser(svcCtx *svc.ServiceContext, revUserId uint, sendUserID uint, msg ctype.Msg, msgID uint) {

	revUserAny, ok1 := UserOnlineWsMap.Load(revUserId)
	sendUserAny, ok2 := UserOnlineWsMap.Load(sendUserID)
	resp := ChatResponse{
		ID:         msgID,
		Msg:        msg,
		MsgPreview: msg.MsgPreview(),
		CreatedAt:  time.Now(),
	}
	revUser := revUserAny.(*UserWsInfo)
	sendUser := sendUserAny.(*UserWsInfo)
	if ok1 && ok2 && sendUserID == revUserId {
		// 自己给自己发
		resp.RevUser = ctype.UserInfo{
			ID:       revUserId,
			NickName: revUser.UserInfo.Nickname,
			Avatar:   revUser.UserInfo.Avatar,
		}
		resp.SendUser = ctype.UserInfo{
			ID:       sendUserID,
			NickName: sendUser.UserInfo.Nickname,
			Avatar:   sendUser.UserInfo.Avatar,
		}
		resp.IsMe = true
		byteData, _ := json.Marshal(resp)
		//revUser.Conn.WriteMessage(websocket.TextMessage, byteData)
		sendWsMapMsg(revUser.WsClientMap, byteData)
		return
	}

	// 在线的情况下，我是可以拿到对方的用户信息的
	// 对方不在线的情况下，我只能通过调用户的rpc方法，去获取用户基本信息

	// 不管怎么样，都要给发送者回传消息的
	// 如果接受者不在线，那么我就要去拿接受者的用户信息

	if !ok1 {
		userBaseInfo, err := redis_service.GetUserBaseInfo(svcCtx.Redis, svcCtx.UserRpc, revUserId)
		if err != nil {
			logx.Error(err)
			return
		}
		resp.RevUser = ctype.UserInfo{
			ID:       revUserId,
			NickName: userBaseInfo.NickName,
			Avatar:   userBaseInfo.Avatar,
		}
	} else {
		resp.RevUser = ctype.UserInfo{
			ID:       revUserId,
			NickName: revUser.UserInfo.Nickname,
			Avatar:   revUser.UserInfo.Avatar,
		}
	}

	// 发送者在线
	resp.SendUser = ctype.UserInfo{
		ID:       sendUserID,
		NickName: sendUser.UserInfo.Nickname,
		Avatar:   sendUser.UserInfo.Avatar,
	}
	resp.IsMe = true
	byteData, _ := json.Marshal(resp)

	//sendUser.Conn.WriteMessage(websocket.TextMessage, byteData)
	sendWsMapMsg(sendUser.WsClientMap, byteData)

	if ok1 {
		// 接收者在线
		resp.IsMe = false
		byteData, _ = json.Marshal(resp)
		//revUser.Conn.WriteMessage(websocket.TextMessage, byteData)
		sendWsMapMsg(revUser.WsClientMap, byteData)
	}
}
