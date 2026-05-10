package logic

import (
	"context"
	"fim_server/common/list_query"
	"fim_server/common/models"
	"fim_server/fim_group/group_api/internal/svc"
	"fim_server/fim_group/group_api/internal/types"
	"fim_server/fim_group/group_models"
	"fim_server/fim_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupValidListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupValidListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupValidListLogic {
	return &GroupValidListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupValidListLogic) GroupValidList(req *types.GroupValidListRequest) (resp *types.GroupValidListResponse, err error) {
	var groupIDList []uint
	l.svcCtx.DB.Model(group_models.GroupMemberModel{}).
		Where("user_id = ? and (role = 1 or role = 2)", req.UserID).
		Select("group_id").Scan(&groupIDList)

	groupMap := map[uint]bool{}
	for _, id := range groupIDList {
		groupMap[id] = true
	}

	where := l.svcCtx.DB.Where("user_id = ?", req.UserID)
	if len(groupIDList) > 0 {
		where = where.Or("group_id in ?", groupIDList)
	}

	groups, count, err := list_query.ListQuery(l.svcCtx.DB, group_models.GroupVerifyModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Preload: []string{"GroupModel"},
		Where:   where,
	})
	if err != nil {
		return nil, err
	}

	var userIDList []uint32
	for _, group := range groups {
		userIDList = append(userIDList, uint32(group.UserID))
	}

	userList, err1 := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})

	resp = new(types.GroupValidListResponse)
	resp.Count = int(count)
	for _, groupVerify := range groups {
		info := types.GroupValidInfoResponse{
			ID:                 groupVerify.ID,
			GrouID:             groupVerify.GroupID,
			UserID:             groupVerify.UserID,
			Status:             groupVerify.Status,
			AdditionalMessages: groupVerify.AdditionalMessages,
			Title:              groupVerify.GroupModel.Title,
			CreatedAt:          groupVerify.CreatedAt.String(),
			Type:               groupVerify.Type,
			Avatar:             groupVerify.GroupModel.Avatar,
			Flag:               "send",
		}
		if groupVerify.VerificationQuestion != nil {
			info.VerificationQuestion = &types.VerificationQuestion{
				Problem1: groupVerify.VerificationQuestion.Problem1,
				Problem2: groupVerify.VerificationQuestion.Problem2,
				Problem3: groupVerify.VerificationQuestion.Problem3,
				Answer1:  groupVerify.VerificationQuestion.Answer1,
				Answer2:  groupVerify.VerificationQuestion.Answer2,
				Answer3:  groupVerify.VerificationQuestion.Answer3,
			}
		}

		if groupMap[groupVerify.GroupID] {
			info.Flag = "rev"
		}

		if err1 == nil {
			info.UserNickname = userList.UserInfo[uint32(info.UserID)].NickName
			info.UserAvatar = userList.UserInfo[uint32(info.UserID)].Avatar
		}

		resp.List = append(resp.List, info)
	}

	return resp, nil
}
