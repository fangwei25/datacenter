package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"web_game/service/user/model"
	"web_game/service/user/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config

	PlayerModel        model.PlayerModel
	AccountAppleModel  model.AccountAppleModel
	AccountFacebookMod model.AccountFacebookModel
	AccountGuestModel  model.AccountGuestModel
	AccountMobileModel model.AccountMobileModel
	AccountPwdModel    model.AccountPwdModel
	AccountWechatModel model.AccountWechatModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:             c,
		PlayerModel:        model.NewPlayerModel(conn, c.CacheRedis),
		AccountAppleModel:  model.NewAccountAppleModel(conn, c.CacheRedis),
		AccountFacebookMod: model.NewAccountFacebookModel(conn, c.CacheRedis),
		AccountGuestModel:  model.NewAccountGuestModel(conn, c.CacheRedis),
		AccountMobileModel: model.NewAccountMobileModel(conn, c.CacheRedis),
		AccountPwdModel:    model.NewAccountPwdModel(conn, c.CacheRedis),
		AccountWechatModel: model.NewAccountWechatModel(conn, c.CacheRedis),
	}
}
