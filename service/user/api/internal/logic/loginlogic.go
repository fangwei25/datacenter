package logic

import (
	"context"
	"math/rand"
	"strings"
	"time"
	"web_game/common/counterx"
	"web_game/common/jwtx"
	"web_game/common/logintype"
	"web_game/common/tracex"
	"web_game/service/counter/rpc/counterclient"
	"web_game/service/trace/rpc/traceclient"
	"web_game/service/user/rpc/userclient"

	"web_game/service/user/api/internal/svc"
	"web_game/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginLogic {
	return LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req types.ReqLogin) (resp *types.ResLogin, err error) {

	//对游客账号的预处理
	if req.LoginType == logintype.Guest && len(strings.TrimSpace(req.Account)) <= 0 {
		req.Account = GenRandomGustAccount()
	}

	res, err := l.svcCtx.UserRpc.Login(l.ctx, &userclient.ReqLogin{
		AppVersion: req.AppVersion,
		Account:    req.Account,
		AccToken:   req.AccToken,
		LoginType:  req.LoginType,
		Time:       req.Time,
		DeviceId:   req.DeviceId,
		Invitation: req.Invitation,
		Channel:    req.Channel,
	})
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire

	accessToken, err := jwtx.GetToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, res.PlayerId)
	if err != nil {
		return nil, err
	}

	//上报counter，trace
	loginTracInfo := &tracex.LoginInfo{
		PlayerId: res.PlayerId,
		Account:  req.Account,
		DeviceId: req.DeviceId,
		//IpAddr:   l.ctx.Value("ip"),
	}
	if res.IsReg {
		_, _ = l.svcCtx.CounterRpc.Update(l.ctx, &counterclient.ReqUpdate{
			OwnerId:    res.PlayerId,
			EventType:  counterx.ActRegister,
			EventField: req.LoginType,
			Value:      1,
		})
		_, _ = l.svcCtx.TraceRpc.PushTrace(l.ctx, &traceclient.ReqTrace{
			TraceName: tracex.Register,
			JsonData:  loginTracInfo.ToJsonStr(),
		})
	} else {
		_, _ = l.svcCtx.CounterRpc.Update(l.ctx, &counterclient.ReqUpdate{
			OwnerId:    res.PlayerId,
			EventType:  counterx.ActLogin,
			EventField: req.LoginType,
			Value:      1,
		})
		_, _ = l.svcCtx.TraceRpc.PushTrace(l.ctx, &traceclient.ReqTrace{
			TraceName: tracex.Register,
			JsonData:  loginTracInfo.ToJsonStr(),
		})
	}

	return &types.ResLogin{
		ErrNo:        0,
		ErrMsg:       "succ",
		Jwt:          accessToken,
		GuestAccount: "",
		IsReg:        res.IsReg,
	}, nil

}

// GenRandomGustAccount 生成长度为32的随机游客名字（大写英文字母）
var randSeed = rand.New(rand.NewSource(time.Now().Unix()))

func GenRandomGustAccount() string {
	accountLength := 32
	bytes := make([]byte, accountLength)
	for i := 0; i < accountLength; i++ {
		b := randSeed.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
