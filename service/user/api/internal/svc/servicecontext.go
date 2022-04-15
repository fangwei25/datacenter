package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"web_game/service/counter/rpc/counterclient"
	"web_game/service/user/api/internal/config"
	"web_game/service/user/rpc/userclient"
)

type ServiceContext struct {
	Config     config.Config
	UserRpc    userclient.User
	CounterRpc counterclient.Counter
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		CounterRpc: counterclient.NewCounter(zrpc.MustNewClient(c.CounterRpc)),
	}
}
