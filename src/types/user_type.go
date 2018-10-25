package types

type ReqUserRegAndLogin struct {
	UserName     string `json:"username"`
	PassWord     string `json:"password"`
	VerifyCode   string `json:"verifycode"`
	VerifyCodeId string `json:"verifycode_id"`
}

type RspUserRegAndLogin struct {
	Token    string `json:"token"`
	UserName string `json:"username"`
	AppId    string `json:"appid"`
	AppKey   string `json:"appkey"`
}

type ReqUpdatePassword struct {
	Token   string `json:"token"`
	OldPass string `json:"old_password"`
	NewPass string `json:"new_password"`
}
