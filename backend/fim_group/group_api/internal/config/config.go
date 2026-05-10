package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Etcd  string
	Database struct {
		DataSource string
	}
	UserRpc  zrpc.RpcClientConf
	GroupRpc zrpc.RpcClientConf
	FileRpc  zrpc.RpcClientConf
	Redis    struct {
		Addr string
		Pwd  string
		DB   int
	}
}
