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

	player, err := l.svcCtx.PlayerModel.FindOne(in.GetPlayerId())
	if err != nil {
		return nil, err
	}
	return &user.ResUserInfo{
		PlayerId:     player.PlayerId,
		Name:         player.Name,
		Gender:       int32(player.Gender),
		AvatorUrl:    player.AvatorUrl,
		InvitationId: int32(player.InvitationId),
		Channel:      player.Channel,
		VipLv:        int32(player.VipLv),
		VipExp:       player.VipExp,
		Level:        int32(player.Level),
		LevelExp:     player.LevelExp,
	}, nil
}
