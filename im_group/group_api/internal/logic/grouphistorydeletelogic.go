package logic

import (
	"context"
	"errors"
	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"
	"im_server/im_group/group_models"
	"im_server/utils/set"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupHistoryDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupHistoryDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupHistoryDeleteLogic {
	return &GroupHistoryDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupHistoryDelete 用户删除群聊天记录
func (l *GroupHistoryDeleteLogic) GroupHistoryDelete(req *types.GroupHistoryDeleteRequest) (resp *types.GroupHistoryDeleteResponse, err error) {
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("用户不是群成员")
	}
	// 去查我删除了哪些聊天记录
	var msgIDList []uint
	l.svcCtx.DB.Model(group_models.GroupUserMsgDeleteModel{}).
		Where("group_id = ? and user_id = ?", req.ID, req.UserID).
		Select("msg_id").Scan(&msgIDList)

	// 和我传的聊天记录做一个交集
	// 1,2,3,4  2,3,6
	addMsgIDList := set.Difference(req.MsgIDList, msgIDList)
	logx.Infof("删除聊天记录的id列表 %v", addMsgIDList)
	if len(addMsgIDList) == 0 {
		return
	}

	// 用户传过来的消息id，消息不一定存在
	var msgIDFindList []uint
	l.svcCtx.DB.Model(group_models.GroupMsgModel{}).
		Where("id in ?", addMsgIDList).
		Select("id").Scan(&msgIDFindList)

	if len(msgIDFindList) != len(addMsgIDList) {
		return nil, errors.New("消息一致性错误")
	}

	var list []group_models.GroupUserMsgDeleteModel
	for _, i2 := range addMsgIDList {
		list = append(list, group_models.GroupUserMsgDeleteModel{
			MsgID:   i2,
			UserID:  req.UserID,
			GroupID: req.ID,
		})
	}
	err = l.svcCtx.DB.Create(&list).Error
	if err != nil {
		return
	}
	return
}
