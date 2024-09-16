package svc

import (
	"gorm.io/gorm"
	"im_server/core"
	"im_server/im_settings/settings_api/internal/config"
	"im_server/im_settings/settings_api/internal/middleware"
	"net/http"
)

type ServiceContext struct {
	Config          config.Config
	DB              *gorm.DB
	AdminMiddleware func(next http.HandlerFunc) http.HandlerFunc
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	return &ServiceContext{
		Config:          c,
		DB:              mysqlDb,
		AdminMiddleware: middleware.NewAdminMiddleware().Handle,
	}
}
