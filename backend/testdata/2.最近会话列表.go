package main

import (
	"fim_server/common/list_query"
	"fim_server/common/models"
	"fim_server/core"
	"fim_server/fim_chat/chat_models"
	"fmt"
)

func main() {
	db := core.InitGorm("host=127.0.0.1 user=root password=root dbname=fim_server_db port=5432 sslmode=disable TimeZone=Asia/Shanghai")

	var userId = 1
	type Data struct {
		SU         uint   `gorm:"column:sU"`
		RU         uint   `gorm:"column:rU"`
		MaxDate    string `gorm:"column:maxDate"`
		MaxPreview string `gorm:"column:maxPreview"`
		IsTop      bool   `gorm:"column:isTop"`
	}
	var list []Data

	db.Table("(?) as u", db.Model(&chat_models.ChatModel{}).
		Select("least(send_user_id, rev_user_id)    as sU",
			"greatest(send_user_id, rev_user_id) as rU",
			" max(created_at)   as maxDate",
			"max(msg_preview) as maxPreview").Where("send_user_id = ? or rev_user_id = ?", userId, userId).
		Group("least(send_user_id, rev_user_id)").
		Group("greatest(send_user_id, rev_user_id)")).
		Order("maxDate desc").Limit(1).Offset(0).Scan(&list)
	fmt.Println(list)

	column := fmt.Sprintf("CASE WHEN EXISTS (SELECT 1 FROM top_user_models WHERE user_id = %d AND (top_user_id = sU OR top_user_id = rU)) THEN 1 ELSE 0 END AS isTop", userId)

	chatList, count, _ := list_query.ListQuery(db, Data{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  1,
			Limit: 10,
			Sort:  "isTop desc, maxDate desc",
		},
		Table: func() (string, any) {
			return "(?) as u", db.Model(&chat_models.ChatModel{}).
				Select("least(send_user_id, rev_user_id) as sU",
					"greatest(send_user_id, rev_user_id) as rU",
					"max(created_at) as maxDate",
					"(select msg_preview from chat_models  where (send_user_id = sU and rev_user_id = rU) or (send_user_id = rU and rev_user_id = sU) order by created_at desc  limit 1) as maxPreview",
					column).
				Where("send_user_id = ? or rev_user_id = ?", userId, userId).
				Group("least(send_user_id, rev_user_id)").
				Group("greatest(send_user_id, rev_user_id)")
		},
	})

	fmt.Println(chatList, count)
	for _, data := range chatList {
		fmt.Println(data.IsTop, data.MaxPreview, data.MaxDate)
	}

}
