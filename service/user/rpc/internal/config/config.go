package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}

	CacheRedis       cache.CacheConf
	Salt             string
	SupportLoginType map[string]bool //支持的登录方式
}
