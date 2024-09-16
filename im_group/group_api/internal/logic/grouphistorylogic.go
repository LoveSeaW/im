package logic

import (
	"context"
	"errors"
	"im_server/common/list_query"
	"im_server/common/models"
	"im_server/common/models/ctype"
	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"
	"im_server/im_group/group_models"
	"im_server/im_user/user_rpc/types/user_rpc"
	"im_server/utils"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupHistoryLogic {
	return &GroupHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type HistoryResponse struct {
	GroupID        uint          `json:"groupID"`
	UserID         uint          `json:"userID"`
	UserNickname   string        `json:"userNickname"`
	UserAvatar     string        `json:"userAvatar"`
	Msg            ctype.Msg     `json:"msg"`
	MsgPreview     string        `json:"msgPreview"`
	ID             uint          `json:"id"`
	MsgType        ctype.MsgType `json:"msgType"`
	CreatedAt      string        `json:"createdAt"`
	IsMe           bool          `json:"isMe"`
	MemberNickname string        `json:"memberNickname"` // 群好友备注
	ShowDate       bool          `json:"showDate"`
}

type HistoryListResponse struct {
	List  []HistoryResponse `json:"list"`
	Count int               `json:"count"`
}

func (l *GroupHistoryLogic) GroupHistory(req *types.GroupHistoryRequest) (resp *HistoryListResponse, err error) {
	// 谁能调这个接口 必须得是这个群的成员
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("该用户不是群成员")
	}
	// 去查我删除了哪些聊天记录
	var msgIDList []uint
	l.svcCtx.DB.Model(group_models.GroupUserMsgDeleteModel{}).
		Where("group_id = ? and user_id = ?", req.ID, req.UserID).
		Select("msg_id").Scan(&msgIDList)
	var query = l.svcCtx.DB.Where("")
	if len(msgIDList) > 0 {
		query.Where("id not in ?", msgIDList)
	}

	groupMsgList, count, err := list_query.ListQuery(l.svcCtx.DB, group_models.GroupMsgModel{GroupID: req.ID}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc",
		},
		Where:   query,
		Preload: []string{"GroupMemberModel"},
	})

	var userIDList []uint32
	for _, model := range groupMsgList {
		userIDList = append(userIDList, uint32(model.SendUserID))
	}
	userIDList = utils.DeduplicationList(userIDList)
	userListResponse, err1 := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})

	utils.ReverseAny(groupMsgList)

	var list = make([]HistoryResponse, 0)
	for index, model := range groupMsgList {
		info := HistoryResponse{
			GroupID:   model.GroupID,
			UserID:    model.SendUserID,
			Msg:       model.Msg,
			ID:        model.ID,
			MsgType:   model.MsgType,
			CreatedAt: model.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		// 这一条消息与上一条消息的时间差值
		if index == 0 {
			info.ShowDate = true
		} else {
			sub := model.CreatedAt.Sub(groupMsgList[index-1].CreatedAt)
			if sub > time.Hour {
				info.ShowDate = true
			}
		}
		if model.GroupMemberModel != nil {
			info.MemberNickname = model.GroupMemberModel.MemberNickname
		}
		if err1 == nil {
			info.UserNickname = userListResponse.UserInfo[uint32(info.UserID)].NickName
			info.UserAvatar = userListResponse.UserInfo[uint32(info.UserID)].Avatar
		}
		if req.UserID == info.UserID {
			info.IsMe = true
		}
		list = append(list, info)
	}

	resp = new(HistoryListResponse)
	resp.List = list
	resp.Count = int(count)
	return
}
