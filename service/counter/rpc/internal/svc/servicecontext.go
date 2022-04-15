package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"web_game/service/counter/rpc/internal/config"
)

type ServiceContext struct {
	Config   config.Config
	RedisObj *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		RedisObj: redis.New(c.RedisCfg.Host, redis.WithPass(c.RedisCfg.Pass)),
	}
}
