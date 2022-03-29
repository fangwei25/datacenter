package logic

import (
	"context"

	"web_game/service/user/rpc/internal/svc"
	"web_game/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user.ReqUserInfo) (*user.ResUserInfo, error) {
	// todo: add your logic here and delete this line

	return &user.ResUserInfo{}, nil
}
