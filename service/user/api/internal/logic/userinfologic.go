package logic

import (
	"context"
	"encoding/json"
	"web_game/service/user/rpc/userclient"

	"web_game/service/user/api/internal/svc"
	"web_game/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserInfoLogic {
	return UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.ResUserInfo, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	res, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &userclient.ReqUserInfo{
		PlayerId: uid,
	})
	if err != nil {
		return nil, err
	}

	return &types.ResUserInfo{
		PlayerId:     res.GetPlayerId(),
		Name:         res.GetName(),
		Gender:       res.GetGender(),
		AvatorUrl:    res.GetAvatorUrl(),
		InvitationId: res.GetInvitationId(),
		Channel:      res.GetChannel(),
		VipLv:        res.GetVipLv(),
		VipExp:       res.GetVipExp(),
		Level:        res.GetLevel(),
		LevelExp:     res.GetLevelExp(),
	}, nil
}
