package logic

import (
	"context"
	"errors"

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
	loginTypeName := GetLoginTypeName(in.LoginType)
	if !l.IsSupportLoginType(loginTypeName) {
		return nil, errors.New("not supported login type")
	}

	if len(in.DeviceId) < 32 {
		return nil, errors.New("invalid device id")
	}

	//根据登录方式不同，选择不同的登录校验逻辑
	switch in.LoginType {
	case LoginTypeJWT:
		return l.loginByJWT(in)
	case LoginTypePwd:
		return l.loginByPwd(in)
	case LoginTypeMobile:
		return l.loginByMobile(in)
	case LoginTypeWechat:
		return l.loginByWeChat(in)
	case LoginTypeGuest:
		return l.loginByGuest(in)
	case LoginTypeApple:
		return l.loginByApple(in)
	case LoginTypeFacebook:
		return l.loginByFacebook(in)
	default:
		return nil, errors.New("not supported login type")
	}
}

func (l *LoginLogic) loginByJWT(in *user.ReqLogin) (*user.ResLogin, error) {
	return nil, errors.New("")
}
func (l *LoginLogic) loginByPwd(in *user.ReqLogin) (*user.ResLogin, error) {
	return nil, errors.New("")
}
func (l *LoginLogic) loginByMobile(in *user.ReqLogin) (*user.ResLogin, error) {
	return nil, errors.New("")
}
func (l *LoginLogic) loginByWeChat(in *user.ReqLogin) (*user.ResLogin, error) {
	return nil, errors.New("")
}
func (l *LoginLogic) loginByGuest(in *user.ReqLogin) (*user.ResLogin, error) {
	return nil, errors.New("")
}
func (l *LoginLogic) loginByApple(in *user.ReqLogin) (*user.ResLogin, error) {
	return nil, errors.New("")
}
func (l *LoginLogic) loginByFacebook(in *user.ReqLogin) (*user.ResLogin, error) {
	return nil, errors.New("")
}

func GetLoginTypeName(loginType int32) string {
	name, ok := LoginTypeName[loginType]
	if !ok {
		return "unknown"
	}
	return name
}

func (l *LoginLogic) IsSupportLoginType(loginTypeName string) bool {
	if v, ok := l.svcCtx.Config.SupportLoginType[loginTypeName]; ok {
		return v
	}
	return false
}
