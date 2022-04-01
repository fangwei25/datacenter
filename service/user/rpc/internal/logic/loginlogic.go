package logic

import (
	"context"
	"errors"
	"google.golang.org/grpc/status"
	"web_game/common/logintype"
	"web_game/service/user/rpc/internal/svc"
	"web_game/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.ReqLogin) (*user.ResLogin, error) {
	//检查是否是支持的登录方式
	if !l.IsSupportLoginType(in.GetLoginType()) {
		return nil, errors.New("not supported login type")
	}

	if len(in.DeviceId) < 32 {
		return nil, errors.New("invalid device id")
	}

	var playerId int64
	var isReg bool
	var err error
	//根据登录方式不同，选择不同的登录校验逻辑
	switch in.LoginType {
	case logintype.Pwd:
		playerId, isReg, err = l.loginByPwd(in)
	case logintype.Mobile:
		playerId, isReg, err = l.loginByMobile(in)
	case logintype.Wechat:
		playerId, isReg, err = l.loginByWeChat(in)
	case logintype.Guest:
		playerId, isReg, err = l.loginByGuest(in)
	case logintype.Apple:
		playerId, isReg, err = l.loginByApple(in)
	case logintype.Facebook:
		playerId, isReg, err = l.loginByFacebook(in)
	default:
		err = status.Error(100, "不支持的登录方式")
	}

	if nil != err {
		return nil, err
	}

	return &user.ResLogin{
		PlayerId: playerId,
		IsReg:    isReg,
	}, nil
}

func (l *LoginLogic) IsSupportLoginType(loginTypeName string) bool {
	if v, ok := l.svcCtx.Config.LoginType[loginTypeName]; ok {
		return v
	}
	return false
}
