package logic

import (
	"google.golang.org/grpc/status"
	"strings"
	"time"
	"web_game/common/cryptx"
	"web_game/common/logintype"
	"web_game/service/user/model"
	"web_game/service/user/rpc/user"
)

func (l *LoginLogic) loginByPwd(in *user.ReqLogin) (playerId int64, isReg bool, err error) {
	account := strings.TrimSpace(in.Account)
	if len(account) <= 0 {
		return 0, false, status.Error(100, "账号不能为空")
	}

	//校验密码
	password := cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.AccToken)

	res, err := l.svcCtx.AccountPwdModel.FindOne(account)
	if nil == err {
		//账号存在，检查密码是否匹配
		if res.Pwd != password {
			return 0, false, status.Error(100, "密码不正确")
		}
		return res.PlayerId, false, nil
	} else {
		//账号不存在，走创建账号流程
		newUser := model.Player{
			Name:           account,
			OriAccountType: logintype.Pwd,
			OriAccount:     account,
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

		newAccount := model.AccountPwd{
			Account:    account,
			PlayerId:   newUser.PlayerId,
			Pwd:        password,
			CreateTime: time.Now(),
			LastLogin:  time.Now(),
		}

		res, err = l.svcCtx.AccountPwdModel.Insert(&newAccount)
		if err != nil {
			return 0, false, status.Error(500, err.Error())
		}
		return newUser.PlayerId, true, nil
	}
}
