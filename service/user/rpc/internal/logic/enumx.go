package logic

//登录注册方式
const (
	LoginTypeJWT      int32 = 0 //JWT登录
	LoginTypePwd      int32 = 1 //账号密码登录
	LoginTypeWechat   int32 = 2 //微信登录
	LoginTypeMobile   int32 = 3 //手机号登录
	LoginTypeGuest    int32 = 4 //游客登录
	LoginTypeApple    int32 = 5 //苹果账户登录
	LoginTypeFacebook int32 = 6 //facebook登录
)

var LoginTypeName = map[int32]string{
	LoginTypeJWT:      "JWT",
	LoginTypePwd:      "Password",
	LoginTypeWechat:   "Wechat",
	LoginTypeMobile:   "Mobile",
	LoginTypeGuest:    "Guest",
	LoginTypeApple:    "Apple",
	LoginTypeFacebook: "Facebook",
}
