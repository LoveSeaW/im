package logic

import (
	"context"
	"errors"
	"fim_server/common/list_query"
	"fim_server/common/models"
	"fim_server/fim_chat/chat_api/internal/svc"
	"fim_server/fim_chat/chat_api/internal/types"
	"fim_server/fim_user/user_rpc/types/user_rpc"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChatSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatSessionLogic {
	return &ChatSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type Data struct {
	SU         uint   `gorm:"column:sU"`
	RU         uint   `gorm:"column:rU"`
	MaxDate    string `gorm:"column:maxDate"`
	MaxPreview string `gorm:"column:maxPreview"`
	IsTop      bool   `gorm:"column:isTop"`
}

func (l *ChatSessionLogic) ChatSession(req *types.ChatSessionRequest) (resp *types.ChatSessionResponse, err error) {

	var friendIDList []uint
	friendRes, err := l.svcCtx.UserRpc.FriendList(l.ctx, &user_rpc.FriendListRequest{
		User: uint32(req.UserID),
	})
	if err != nil {
		logx.Error(err)
		return nil, errors.New("用户服务错误")
	}
	for _, info := range friendRes.FriendList {
		friendIDList = append(friendIDList, uint(info.UserId))
	}

	chatList, count, _ := list_query.ListQuery(l.svcCtx.DB, Data{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "isTop desc, maxDate desc",
		},
		Table: func() (string, any) {
			// 内层：GROUP BY 获取唯一聊天对
			inner := l.svcCtx.DB.Table("chat_models").
				Select("least(send_user_id, rev_user_id) as sU",
					"greatest(send_user_id, rev_user_id) as rU",
					"max(created_at) as maxDate").
				Where("(send_user_id = ? or rev_user_id = ?) and id not in (select chat_id from user_chat_delete_models where user_id = ?) and ((send_user_id = ? and rev_user_id in ?) or (rev_user_id = ? and send_user_id in ?))",
					req.UserID, req.UserID, req.UserID, req.UserID, friendIDList, req.UserID, friendIDList).
				Group("least(send_user_id, rev_user_id)").
				Group("greatest(send_user_id, rev_user_id)")
			// 外层：在 grouped 结果上追加 maxPreview 和 isTop
			return "(?) as u", l.svcCtx.DB.Table("(?) as grouped", inner).
				Select("sU", "rU", "maxDate",
					fmt.Sprintf("(select msg_preview from chat_models as c where least(c.send_user_id, c.rev_user_id) = grouped.sU and greatest(c.send_user_id, c.rev_user_id) = grouped.rU and c.id not in (select chat_id from user_chat_delete_models where user_id = %d) order by c.created_at desc limit 1) as maxPreview", req.UserID),
					fmt.Sprintf("CASE WHEN EXISTS (SELECT 1 FROM top_user_models WHERE user_id = %d AND (top_user_id = grouped.sU OR top_user_id = grouped.rU)) THEN 1 ELSE 0 END AS isTop", req.UserID))
		},
	})

	var userIDList []uint32
	for _, data := range chatList {
		if data.RU != req.UserID {
			userIDList = append(userIDList, uint32(data.RU))
		}
		if data.SU != req.UserID {
			userIDList = append(userIDList, uint32(data.SU))
		}
		if data.SU == req.UserID && req.UserID == data.RU {
			// 自己和自己聊
			userIDList = append(userIDList, uint32(req.UserID))
		}
	}
	response, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})
	if err != nil {
		logx.Error(err)
		return nil, errors.New("用户服务错误")
	}

	userOnlineRes, err := l.svcCtx.UserRpc.UserOnlineList(l.ctx, &user_rpc.UserOnlineListRequest{})
	if err != nil {
		logx.Error(err)
		return nil, errors.New("用户服务错误")
	}
	var onlineUserMap = map[uint]bool{}
	for _, u := range userOnlineRes.UserIdList {
		onlineUserMap[uint(u)] = true
	}

	var list = make([]types.ChatSession, 0)
	for _, data := range chatList {
		s := types.ChatSession{
			CreatedAt:  data.MaxDate,
			MsgPreview: data.MaxPreview,
			IsTop:      data.IsTop,
		}
		if data.RU != req.UserID {
			s.UserID = data.RU

			s.Avatar = response.UserInfo[uint32(s.UserID)].Avatar
			s.Nickname = response.UserInfo[uint32(s.UserID)].NickName
		}
		if data.SU != req.UserID {
			s.UserID = data.SU
			s.Avatar = response.UserInfo[uint32(s.UserID)].Avatar
			s.Nickname = response.UserInfo[uint32(s.UserID)].NickName
		}
		if data.SU == req.UserID && data.RU == req.UserID {
			s.UserID = data.SU
			s.Avatar = response.UserInfo[uint32(s.UserID)].Avatar
			s.Nickname = response.UserInfo[uint32(s.UserID)].NickName
		}
		s.IsOnlone = onlineUserMap[s.UserID]

		list = append(list, s)
	}

	return &types.ChatSessionResponse{List: list, Count: count}, nil
}
