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
	SU         uint   `gorm:"column:s_u"`
	RU         uint   `gorm:"column:r_u"`
	MaxDate    string `gorm:"column:max_date"`
	MaxPreview string `gorm:"column:max_preview"`
	IsTop      bool   `gorm:"column:is_top"`
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
	if len(friendIDList) == 0 {
		return &types.ChatSessionResponse{List: []types.ChatSession{}, Count: 0}, nil
	}

	chatList, count, _ := list_query.ListQuery(l.svcCtx.DB, Data{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "is_top desc, max_date desc",
		},
		Table: func() (string, any) {
			// 内层：GROUP BY 获取唯一聊天对
			inner := l.svcCtx.DB.Table("chat_models").
				Select("least(send_user_id, rev_user_id) as s_u",
					"greatest(send_user_id, rev_user_id) as r_u",
					"max(created_at) as max_date").
				Where("(send_user_id = ? or rev_user_id = ?) and id not in (select chat_id from user_chat_delete_models where user_id = ?) and ((send_user_id = ? and rev_user_id in ?) or (rev_user_id = ? and send_user_id in ?))",
					req.UserID, req.UserID, req.UserID, req.UserID, friendIDList, req.UserID, friendIDList).
				Group("least(send_user_id, rev_user_id)").
				Group("greatest(send_user_id, rev_user_id)")
			// 外层：在 grouped 结果上追加 maxPreview 和 isTop
			return "(?) as u", l.svcCtx.DB.Table("(?) as grouped", inner).
				Select("s_u", "r_u", "max_date",
					fmt.Sprintf("(select msg_preview from chat_models as c where least(c.send_user_id, c.rev_user_id) = grouped.s_u and greatest(c.send_user_id, c.rev_user_id) = grouped.r_u and c.id not in (select chat_id from user_chat_delete_models where user_id = %d) order by c.created_at desc limit 1) as max_preview", req.UserID),
					fmt.Sprintf("CASE WHEN EXISTS (SELECT 1 FROM top_user_models WHERE user_id = %d AND (top_user_id = grouped.s_u OR top_user_id = grouped.r_u)) THEN true ELSE false END AS is_top", req.UserID))
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
