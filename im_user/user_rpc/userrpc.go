package main

import (
	"flag"
	"fmt"
	"im_server/common/zrpc_interceptor"
	"im_server/im_user/user_rpc/internal/config"
	"im_server/im_user/user_rpc/internal/server"
	"im_server/im_user/user_rpc/internal/svc"
	"im_server/im_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/userrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user_rpc.RegisterUsersServer(grpcServer, server.NewUsersServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(zrpc_interceptor.ServerUnaryInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}