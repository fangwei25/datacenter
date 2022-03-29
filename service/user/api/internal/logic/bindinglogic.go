package logic

import (
	"context"

	"web_game/service/user/api/internal/svc"
	"web_game/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindingLogic(ctx context.Context, svcCtx *svc.ServiceContext) BindingLogic {
	return BindingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindingLogic) Binding(req types.ReqBinding) (resp *types.ResBindding, err error) {
	// todo: add your logic here and delete this line

	return
}
