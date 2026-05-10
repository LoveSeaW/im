package logic

import (
	"context"
	"fim_server/common/list_query"
	"fim_server/common/models"
	"fim_server/fim_group/group_api/internal/svc"
	"fim_server/fim_group/group_api/internal/types"
	"fim_server/fim_group/group_models"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type GroupSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSessionLogic {
	return &GroupSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type SessionData struct {
	GroupID       uint   `gorm:"column:g_id"`
	NewMsgDate    string `gorm:"column:new_msg_date"`
	NewMsgPreview string `gorm:"column:new_msg_preview"`
	IsTop         bool   `gorm:"column:is_top"`
}

func (l *GroupSessionLogic) GroupSession(req *types.GroupSessionRequest) (resp *types.GroupSessionListResponse, err error) {
	// 先查我有哪些群
	var userGroupIDList []uint
	l.svcCtx.DB.Model(group_models.GroupMemberModel{}).
		Where("user_id = ?", req.UserID).
		Select("group_id").Scan(&userGroupIDList)
	if len(userGroupIDList) == 0 {
		return &types.GroupSessionListResponse{List: []types.GroupSessionResponse{}, Count: 0}, nil
	}

	// 查哪些聊天记录是被删掉的
	var msgDeleteIDList []uint
	l.svcCtx.DB.Model(group_models.GroupUserMsgDeleteModel{}).Where("group_id in ?", userGroupIDList).Select("msg_id").Scan(&msgDeleteIDList)

	query := l.svcCtx.DB.Where("group_id in (?)", userGroupIDList)
	if len(msgDeleteIDList) > 0 {
		query = query.Where("id not in ?", msgDeleteIDList)
	}

	sessionList, count, _ := list_query.ListQuery(l.svcCtx.DB, SessionData{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "is_top desc, new_msg_date desc",
		},
		Debug: true,
		Table: func() (string, any) {
			// 内层：GROUP BY 获取每个群的最后消息时间
			inner := l.svcCtx.DB.Model(&group_models.GroupMsgModel{}).
				Select("group_id as g_id",
					"max(created_at) as new_msg_date").
				Where(query).
				Group("group_id")
			// 外层：在 grouped 结果上追加 newMsgPreview 和 isTop
			return "(?) as u", l.svcCtx.DB.Table("(?) as grouped", inner).
				Select("g_id", "new_msg_date",
					"(select msg_preview from group_msg_models as g where g.group_id = grouped.g_id order by g.created_at desc limit 1) as new_msg_preview",
					fmt.Sprintf("CASE WHEN EXISTS (SELECT 1 FROM group_user_top_models WHERE user_id = %d AND group_user_top_models.group_id = grouped.g_id) THEN true ELSE false END AS is_top", req.UserID))
		},
	})

	var groupIDList []uint
	for _, data := range sessionList {
		groupIDList = append(groupIDList, data.GroupID)
	}
	var groupListModel []group_models.GroupModel
	l.svcCtx.DB.Find(&groupListModel, "id in ?", groupIDList)
	var groupMap = map[uint]group_models.GroupModel{}
	for _, model := range groupListModel {
		groupMap[model.ID] = model
	}

	resp = new(types.GroupSessionListResponse)
	for _, data := range sessionList {
		resp.List = append(resp.List, types.GroupSessionResponse{
			GroupID:       data.GroupID,
			Title:         groupMap[data.GroupID].Title,
			Avatar:        groupMap[data.GroupID].Avatar,
			NewMsgDate:    data.NewMsgDate,
			NewMsgPreview: data.NewMsgPreview,
			IsTop:         data.IsTop,
		})
	}
	resp.Count = int(count)
	return
}
