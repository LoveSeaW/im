package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"im_server/common/models/ctype"
	"im_server/common/response"
	"im_server/common/service/redis_service"
	"im_server/im_file/file_rpc/files"
	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"
	"im_server/im_group/group_models"
	"im_server/im_user/user_rpc/types/user_rpc"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type UserWsInfo struct {
	UserInfo    ctype.UserInfo             // 用户信息
	WsClientMap map[string]*websocket.Conn // 这个用户管理的所有ws客户端
}

var UserOnlineWsMap = map[uint]*UserWsInfo{}

type ChatRequest struct {
	GroupID uint      `json:"groupID"` // 群id
	Msg     ctype.Msg `json:"msg"`     // 消息
}

type ChatResponse struct {
	GroupID        uint          `json:"groupID"`
	UserID         uint          `json:"userID"`
	UserNickname   string        `json:"userNickname"`
	UserAvatar     string        `json:"userAvatar"`
	Msg            ctype.Msg     `json:"msg"`
	ID             uint          `json:"id"`
	MsgType        ctype.MsgType `json:"msgType"`
	CreatedAt      time.Time     `json:"createdAt"`
	IsMe           bool          `json:"isMe"`
	MemberNickname string        `json:"memberNickname"` // 群好友备注
	MsgPreview     string        `json:"msgPreview"`
}

func groupChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupChatRequest
		if err := httpx.ParseHeaders(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		var upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// 鉴权 true表示放行，false表示拦截
				return true
			},
		}

		conn, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}

		addr := conn.RemoteAddr().String()
		logx.Infof("用户建立ws连接 %s", addr)
		defer func() {
			conn.Close()

			userWsInfo, ok := UserOnlineWsMap[req.UserID]
			if ok {
				// 删除的退出的那个ws信息
				delete(userWsInfo.WsClientMap, addr)
			}
			if userWsInfo != nil && len(userWsInfo.WsClientMap) == 0 {
				// 全退完了
				delete(UserOnlineWsMap, req.UserID)
			}
		}()

		baseInfoResponse, err := svcCtx.UserRpc.UserBaseInfo(context.Background(), &user_rpc.UserBaseInfoRequest{
			UserId: uint32(req.UserID),
		})
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}

		userInfo := ctype.UserInfo{
			ID:       req.UserID,
			NickName: baseInfoResponse.NickName,
			Avatar:   baseInfoResponse.Avatar,
		}

		userWsInfo, ok := UserOnlineWsMap[req.UserID]
		if !ok {
			userWsInfo = &UserWsInfo{
				UserInfo: userInfo,
				WsClientMap: map[string]*websocket.Conn{
					addr: conn,
				},
			}
			// 代表这个用户第一次来
			UserOnlineWsMap[req.UserID] = userWsInfo
		}
		_, ok1 := userWsInfo.WsClientMap[addr]
		if !ok1 {
			// 代表这个用户二开及以上
			UserOnlineWsMap[req.UserID].WsClientMap[addr] = conn
		}

		for {
			// 消息类型，消息，错误
			_, p, err1 := conn.ReadMessage()
			if err1 != nil {
				// 用户断开聊天
				fmt.Println(err1)
				break
			}

			var request ChatRequest
			err = json.Unmarshal(p, &request)
			if err != nil {
				logx.Error(err)
				SendTipErrMsg(conn, "参数解析失败")
				continue
			}

			// 校验消息
			msgValidateErr := request.Msg.Validate()
			if msgValidateErr != nil {
				SendTipErrMsg(conn, msgValidateErr.Error())
				continue
			}

			// 判断自己是不是这个群的成员
			var member group_models.GroupMemberModel
			err = svcCtx.DB.Preload("GroupModel").Take(&member, "group_id = ? and user_id = ?", request.GroupID, req.UserID).Error
			if err != nil {
				// 自己不是群的成员
				SendTipErrMsg(conn, "你还不是这个群的成员呢")
				continue
			}

			if member.GroupModel.IsProhibition && member.Role == 3 {
				// 开启了全员禁言  判断当前用户的角色，如果是普通用户，就不能聊天
				SendTipErrMsg(conn, "当前群正在全员禁言中")
				continue
			}

			// 我是不是被禁言了
			if member.GetProhibitionTime(svcCtx.Redis, svcCtx.DB) != nil {
				SendTipErrMsg(conn, "当前用户正在禁言中")
				continue
			}

			switch request.Msg.Type {
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
			case ctype.WithdrawMsgType: // 撤回消息
				// 校验
				withdrawMsg := request.Msg.WithdrawMsg
				if withdrawMsg == nil {
					SendTipErrMsg(conn, "撤回消息的格式错误")
					continue
				}
				if withdrawMsg.MsgID == 0 {
					SendTipErrMsg(conn, "撤回消息id为空")
					continue
				}
				// 去找消息
				var groupMsg group_models.GroupMsgModel
				err = svcCtx.DB.Take(&groupMsg, "group_id = ? and id = ?", request.GroupID, withdrawMsg.MsgID).Error
				if err != nil {
					SendTipErrMsg(conn, "原消息不存在")
					continue
				}
				// 原消息不能是撤回消息
				if groupMsg.MsgType == ctype.WithdrawMsgType {
					SendTipErrMsg(conn, "该消息已撤回")
					continue
				}
				// 要去拿我在这个群的角色

				// 自己是普通用户
				if member.Role == 3 {
					// 如果是自己撤自己的
					if req.UserID != groupMsg.SendUserID {
						SendTipErrMsg(conn, "普通用户只能撤回自己的消息")
						continue
					}
					// 要判断时间是不是大于了2分钟
					now := time.Now()
					if now.Sub(groupMsg.CreatedAt) > 2*time.Minute {
						SendTipErrMsg(conn, "只能撤回两分钟以内的消息")
						continue
					}
				}

				// 查这个消息的用户，在这个群里的角色
				var msgUserRole int8 = 3
				err = svcCtx.DB.Model(group_models.GroupMemberModel{}).
					Where("group_id = ? and user_id = ?", request.GroupID, groupMsg.SendUserID).
					Select("role").
					Scan(&msgUserRole).Error
				// 这里有可能查不到  原因是这个消息的用户退群了，那么也是可以撤回的

				// 如果是管理员撤回  它能撤自己和用户的，没有时间限制
				if member.Role == 2 {
					// 不能撤群主和别的管理员
					if msgUserRole == 1 || (msgUserRole == 2 && groupMsg.SendUserID != req.UserID) {
						SendTipErrMsg(conn, "管理员只能撤回自己或者普通用户的消息")
						continue
					}
				}
				// 如果是群主，那就能撤管理员和用户的

				// 代表消息可以撤回了
				// 修改原消息
				var content = "撤回了一条消息"
				content = "你" + content
				// 前端可以判断，这个消息如果不是isMe，就可以把你替换成对方的昵称

				originMsg := groupMsg.Msg
				originMsg.WithdrawMsg = nil // 这里可能会出现循环引用，所以拷贝了这个值，并且把撤回消息置空了

				svcCtx.DB.Model(&groupMsg).Updates(group_models.GroupMsgModel{
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
				var msgModel group_models.GroupMsgModel
				err = svcCtx.DB.Take(&msgModel, "group_id = ? and id = ?", request.GroupID, request.Msg.ReplyMsg.MsgID).Error
				if err != nil {
					SendTipErrMsg(conn, "消息不存在")
					continue
				}

				// 不能回复撤回消息
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipErrMsg(conn, "该消息已撤回")
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
				var msgModel group_models.GroupMsgModel
				err = svcCtx.DB.Take(&msgModel, "group_id = ? and id = ?", request.GroupID, request.Msg.QuoteMsg.MsgID).Error
				if err != nil {
					SendTipErrMsg(conn, "消息不存在")
					continue
				}

				// 不能回复撤回消息
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipErrMsg(conn, "该消息已撤回")
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
			}

			msgID := insertMsg(svcCtx.DB, conn, member, request.Msg)
			// 遍历这个用户列表，去找ws的客户端
			sendGroupOnlineUserMsg(
				svcCtx.DB,
				member,
				request.Msg,
				msgID,
			)
		}
	}
}

func insertMsg(db *gorm.DB, conn *websocket.Conn, member group_models.GroupMemberModel, msg ctype.Msg) uint {
	switch msg.Type {
	case ctype.WithdrawMsgType:
		fmt.Println("撤回消息自己是不入库的")
		return 0
	}
	groupMsg := group_models.GroupMsgModel{
		GroupID:       member.GroupID,
		SendUserID:    member.UserID,
		GroupMemberID: member.ID,
		MsgType:       msg.Type,
		Msg:           msg,
	}
	groupMsg.MsgPreview = groupMsg.MsgPreviewMethod()
	err := db.Create(&groupMsg).Error
	if err != nil {
		logx.Error(err)
		SendTipErrMsg(conn, "消息保存失败")
		return 0
	}
	return groupMsg.ID
}

// 给这个群的用户发消息
func sendGroupOnlineUserMsg(db *gorm.DB, member group_models.GroupMemberModel, msg ctype.Msg, msgID uint) {

	// 查在线的用户列表
	userOnlineIDList := getOnlineUserIDList()
	// 查这个群的成员 并且在线
	var groupMemberOnlineIDList []uint
	db.Model(group_models.GroupMemberModel{}).
		Where("group_id = ? and user_id in ?", member.GroupID, userOnlineIDList).
		Select("user_id").Scan(&groupMemberOnlineIDList)

	// 构造响应
	var chatResponse = ChatResponse{
		GroupID:        member.GroupID,
		UserID:         member.UserID,
		Msg:            msg,
		ID:             msgID,
		MsgType:        msg.Type,
		CreatedAt:      time.Now(),
		MemberNickname: member.MemberNickname,
		MsgPreview:     msg.MsgPreview(),
	}

	wsInfo, ok := UserOnlineWsMap[member.UserID]
	if ok {
		chatResponse.UserNickname = wsInfo.UserInfo.NickName
		chatResponse.UserAvatar = wsInfo.UserInfo.Avatar
	}

	for _, u := range groupMemberOnlineIDList {
		wsUserInfo, ok2 := UserOnlineWsMap[u]
		if !ok2 {
			continue
		}
		chatResponse.IsMe = false
		// 判断isMe
		if wsUserInfo.UserInfo.ID == member.UserID {
			chatResponse.IsMe = true
		}

		byteData, _ := json.Marshal(chatResponse)
		for _, w2 := range wsUserInfo.WsClientMap {
			w2.WriteMessage(websocket.TextMessage, byteData)
		}
	}
}

func getOnlineUserIDList() (userOnlineIDList []uint) {
	for u, _ := range UserOnlineWsMap {
		userOnlineIDList = append(userOnlineIDList, u)
	}
	return
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
