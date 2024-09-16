package svc

import (
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"im_server/common/zrpc_interceptor"
	"im_server/core"
	"im_server/im_file/file_rpc/files"
	"im_server/im_file/file_rpc/types/file_rpc"
	"im_server/im_group/group_api/internal/config"
	"im_server/im_group/group_api/internal/middleware"
	"im_server/im_group/group_rpc/groups"
	"im_server/im_group/group_rpc/types/group_rpc"
	"im_server/im_user/user_rpc/types/user_rpc"
	"im_server/im_user/user_rpc/users"
	"net/http"
)

type ServiceContext struct {
	Config          config.Config
	DB              *gorm.DB
	UserRpc         user_rpc.UsersClient
	GroupRpc        group_rpc.GroupsClient
	FileRpc         file_rpc.FilesClient
	Redis           *redis.Client
	AdminMiddleware func(next http.HandlerFunc) http.HandlerFunc
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Pwd, c.Redis.DB)
	return &ServiceContext{
		Config:          c,
		DB:              mysqlDb,
		UserRpc:         users.NewUsers(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		GroupRpc:        groups.NewGroups(zrpc.MustNewClient(c.GroupRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		FileRpc:         files.NewFiles(zrpc.MustNewClient(c.FileRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		Redis:           client,
		AdminMiddleware: middleware.NewAdminMiddleware().Handle,
	}
}
