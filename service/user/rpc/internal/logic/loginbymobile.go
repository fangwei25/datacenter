package logic

import (
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"strings"
	"time"
	"web_game/common/logintype"
	"web_game/common/sms"
	"web_game/service/user/model"
	"web_game/service/user/rpc/user"
)

func (l *LoginLogic) loginByMobile(in *user.ReqLogin) (playerId int64, isReg bool, err error) {
	account := strings.TrimSpace(in.Account)
	if len(account) <= 0 {
		return 0, false, status.Error(100, "账号不能为空")
	}

	//校验 验证码
	checkNum := in.GetAccToken()
	if checkNum == "" || sms.GetCheckNum(account) != checkNum {
		if superSms, ok := l.svcCtx.Config.SuperSmsCode[checkNum]; !ok || !superSms {
			return 0, false, status.Error(200, "短信验证码不正确")
		} else {
			logx.Info("super sms code triggered, account: ", account, ", code: ", checkNum)
		}
	}

	res, err := l.svcCtx.AccountPwdModel.FindOne(account)
	if nil == err {
		//账号存在
		return res.PlayerId, false, nil
	} else {
		//账号不存在，走创建账号流程
		newUser := model.Player{
			Name:           account,
			OriAccountType: logintype.Mobile,
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

		newAccount := model.AccountMobile{
			Account:    account,
			PlayerId:   newUser.PlayerId,
			CreateTime: time.Now(),
			LastLogin:  time.Now(),
		}

		res, err = l.svcCtx.AccountMobileModel.Insert(&newAccount)
		if err != nil {
			return 0, false, status.Error(500, err.Error())
		}
		return newUser.PlayerId, true, nil
	}
}
