package sms

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type smsInfo struct {
	Phone         string    `json:"phone"` //手机号
	SendTime      time.Time `json:"-"`     //发送时间
	CheckNum      string    `json:"code"`  //验证码
	ClientPkgName string    `json:"client_pkg_name"`
}

var smsSendTime map[string]*smsInfo

func SendSms(phone string) error {
	sms, ok := smsSendTime[phone]
	if ok {
		//发送过短信，检查发送间隔
		nonce := time.Now()
		if nonce.Sub(sms.SendTime) < 60*time.Second {
			//上次发送时间间隔不足60秒
			return errors.New("cd period")
		}
		sms.SendTime = nonce
		sms.CheckNum = genCheckNum(6)
	} else {
		//没有发送过短信，创建一个信息保存
		sms = &smsInfo{
			Phone:    phone,
			SendTime: time.Now(),
			CheckNum: genCheckNum(6),
		}
		smsSendTime[phone] = sms
	}

	//todo SendHttpPostAsync("sendSms", sms, HttpHandlerSendSMS, phone, sms.CheckNum)
	return nil
}

func GetCheckNum(phone string) string {
	sms, ok := smsSendTime[phone]
	if ok {
		return sms.CheckNum
	}
	return ""
}

func genCheckNum(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, err := fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
		if err != nil {
			return ""
		}
	}
	return sb.String()

}
