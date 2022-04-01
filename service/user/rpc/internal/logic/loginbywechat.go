package logic

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"google.golang.org/grpc/status"
	"time"
	"web_game/common/logintype"
	"web_game/service/user/model"
	"web_game/service/user/rpc/user"
)

func (l *LoginLogic) loginByWeChat(in *user.ReqLogin) (playerId int64, isReg bool, err error) {
	wc := wechat.NewWechat()
	cfg := &offConfig.Config{
		AppID:     l.svcCtx.Config.Wechat.AppID,
		AppSecret: l.svcCtx.Config.Wechat.AppSecret,
		Token:     "xxx",
		Cache:     cache.NewMemory(),
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	oa := officialAccount.GetOauth()
	res, err := oa.CheckAccessToken(in.GetAccToken(), in.GetAccount())
	if err != nil {
		//access token 校验失败
		return 0, false, status.Error(500, err.Error())
	} else if !res {
		return 0, false, status.Error(100, "微信授权失败")
	}

	userInfo, err := oa.GetUserInfo(in.GetAccToken(), in.GetAccount(), "") //注册的时候用这个接口
	if err != nil {
		return 0, false, status.Error(500, err.Error())
	}

	result, err := l.svcCtx.AccountWechatModel.FindOne(in.GetAccount())
	if nil == err {
		//账号存在
		return result.PlayerId, false, nil
	} else {
		//账号不存在，走创建账号流程
		newUser := model.Player{
			Name:           userInfo.Nickname,
			AvatorUrl:      userInfo.HeadImgURL,
			Gender:         int64(userInfo.Sex),
			OriAccountType: logintype.Wechat,
			OriAccount:     userInfo.OpenID,
			CreateTime:     time.Now(),
			UpdateTime:     time.Now(),
		}
		res, err := l.svcCtx.PlayerModel.Insert(&newUser)
		if err != nil {
			return 0, false, status.Error(500, err.Error())
		}
		newUser.PlayerId, err = res.LastInsertId()
		if err != nil {
			return 0, false, status.Error(500, err.Error())
		}

		newAccount := model.AccountWechat{
			Account:     userInfo.OpenID,
			UnionId:     userInfo.Unionid,
			AccessToken: in.GetAccToken(),

			PlayerId:   newUser.PlayerId,
			CreateTime: time.Now(),
			LastLogin:  time.Now(),
		}

		res, err = l.svcCtx.AccountWechatModel.Insert(&newAccount)
		if err != nil {
			return 0, false, status.Error(500, err.Error())
		}
		return newUser.PlayerId, true, nil
	}
}
