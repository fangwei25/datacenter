package logic

import (
	fb "github.com/huandu/facebook/v2"
	"google.golang.org/grpc/status"
	"time"
	"web_game/common/logintype"
	"web_game/service/user/model"
	"web_game/service/user/rpc/user"
)

func (l *LoginLogic) loginByFacebook(in *user.ReqLogin) (playerId int64, isReg bool, err error) {
	globalApp := fb.New(l.svcCtx.Config.FaceBook.AppID, l.svcCtx.Config.FaceBook.AppSecret)
	session := globalApp.Session(globalApp.AppAccessToken())
	res, err := session.Get("/debug_token", nil)
	if err != nil {
		return 0, false, err
	}

	var userInfo struct {
		Name    string
		Id      string
		Picture struct {
			Data struct {
				Height       int32
				Width        int32
				IsSilhouette bool
				Url          string
			}
		}
	}

	res, err = session.Get("/me", fb.Params{"fields": "id,name,picture,gender"})
	_ = res.Decode(&userInfo)

	result, err := l.svcCtx.AccountFacebookMod.FindOne(in.GetAccount())
	if nil == err {
		//账号存在
		return result.PlayerId, false, nil
	} else {
		//账号不存在，走创建账号流程
		newUser := model.Player{
			Name:           userInfo.Name,
			AvatorUrl:      userInfo.Picture.Data.Url,
			OriAccountType: logintype.Facebook,
			OriAccount:     userInfo.Id,
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

		newAccount := model.AccountFacebook{
			Account:       userInfo.Id,
			IdentityToken: in.GetAccToken(),
			PlayerId:      newUser.PlayerId,
			CreateTime:    time.Now(),
			LastLogin:     time.Now(),
		}

		res, err = l.svcCtx.AccountFacebookMod.Insert(&newAccount)
		if err != nil {
			return 0, false, status.Error(500, err.Error())
		}
		return newUser.PlayerId, true, nil
	}
}
