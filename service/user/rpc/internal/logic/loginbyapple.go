package logic

import (
	"google.golang.org/grpc/status"
	"time"
	"web_game/common/logintype"
	"web_game/service/user/model"
	"web_game/service/user/rpc/user"

	"github.com/Timothylock/go-signin-with-apple/apple"
)

func (l *LoginLogic) loginByApple(in *user.ReqLogin) (playerId int64, isReg bool, err error) {
	// Get the email
	claim, err := apple.GetClaims(in.GetAccToken())
	if err != nil {
		//access token 校验失败
		return 0, false, status.Error(200, "授权失败")
	}
	email := (*claim)["email"].(string)
	openId := (*claim)["sub"].(string)
	if openId != in.GetAccount() {
		//提供的openId和解析出来的openId不一致
		return 0, false, status.Error(200, "授权失败，openId不一致")
	}

	result, err := l.svcCtx.AccountAppleModel.FindOne(in.GetAccount())
	if nil == err {
		//账号存在
		return result.PlayerId, false, nil
	} else {
		//账号不存在，走创建账号流程
		newUser := model.Player{
			Name:           email,
			OriAccountType: logintype.Apple,
			OriAccount:     openId,
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

		newAccount := model.AccountApple{
			Account:       openId,
			Email:         email,
			IdentityToken: in.GetAccToken(),

			PlayerId:   newUser.PlayerId,
			CreateTime: time.Now(),
			LastLogin:  time.Now(),
		}

		res, err = l.svcCtx.AccountAppleModel.Insert(&newAccount)
		if err != nil {
			return 0, false, status.Error(500, err.Error())
		}
		return newUser.PlayerId, true, nil
	}
}
