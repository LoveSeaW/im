package core

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitGorm(DataSource string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(DataSource), &gorm.Config{})
	if err != nil {
		panic("连接postgresql数据库失败, error=" + err.Error())
	} else {
		fmt.Println("连接postgresql数据库成功")
	}
	return db
}
