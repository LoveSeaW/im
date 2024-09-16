package svc

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"im_server/core"
	"im_server/im_user/user_rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.RedisConf.Addr, c.RedisConf.Pwd, c.RedisConf.DB)
	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
		Redis:  client,
	}
}
