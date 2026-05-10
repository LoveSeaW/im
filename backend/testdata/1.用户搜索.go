package main

import (
	"fim_server/common/list_query"
	"fim_server/common/models"
	"fim_server/core"
	"fim_server/fim_user/user_models"
	"fmt"
)

func main() {
	db := core.InitGorm("host=127.0.0.1 user=root password=root dbname=fim_server_db port=5432 sslmode=disable TimeZone=Asia/Shanghai")

	key := "3"
	friends, count, err := list_query.ListQuery(db, user_models.UserConfModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			//Page:  req.Page,
			//Limit: req.Limit,
		},
		Preload: []string{"UserModel"},
		Joins:   "left join user_models um on um.id = user_conf_models.user_id",
		Where:   db.Where("(user_conf_models.search_user <> 0 or user_conf_models.search_user is not null)   and (user_conf_models.search_user = 1 and um.id = ?)   or (user_conf_models.search_user = 2 and (    um.id = ? or um.nickname ILIKE ? ))", key, key, fmt.Sprintf("%%%s%%", key)),
	})
	fmt.Println(err)
	fmt.Println(count)
	for _, friend := range friends {
		fmt.Println(friend.UserID, friend.UserModel.Nickname)
	}
}
