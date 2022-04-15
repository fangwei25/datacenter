package tracex

import "encoding/json"

// LoginInfo 转换工具
type LoginInfo struct {
	PlayerId int64  `json:"player_id"` //玩家ID
	Account  string `json:"account"`   //玩家登录账号 JWT登录则为生成JWT时的授权账号
	Platform string `json:"platform"`  //登录方式
	DeviceId string `json:"device_id"` //设备ID
	IpAddr   string `json:"ip_addr"`   //登录IP
}

func (l LoginInfo) ToJsonStr() string {
	jsonByte, err := json.Marshal(l)
	if nil != err {
		return ""
	}
	return string(jsonByte)
}
