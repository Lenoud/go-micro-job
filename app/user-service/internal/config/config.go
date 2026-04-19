package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	MySQL struct {
		DataSource string
	}
	RedisConf struct {
		Addr     string
		Password string
	}
	JWT struct {
		AccessSecret string
		AccessExpire int64
	}
}
