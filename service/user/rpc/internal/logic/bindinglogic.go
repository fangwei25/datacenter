package logic

import (
	"context"

	"web_game/service/user/rpc/internal/svc"
	"web_game/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindingLogic {
	return &BindingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindingLogic) Binding(in *user.ReqBinding) (*user.ResBinding, error) {
	// todo: add your logic here and delete this line

	return &user.ResBinding{}, nil
}
