package types

type ReqWechatPay struct {
	AppId    string `json:"appid"`
	Amount   int    `json:"amount"`
	Rand     string `json:"rand"`
	TradeNum string `json:"tradeNum"`
	Sign     string `json:"sign"`
}
type RspWechatPay struct {
	ReturnCode string `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
	Sign       string `json:"sign"`         // 签名
	ResultCode string `json:"result_code"`  // 业务结果
	ErrCode    string `json:"err_code"`     // 错误码
	ErrCodeDes string `json:"err_code_des"` // 错误描述
	CodeUrl    string `json:"code_url"`
}
