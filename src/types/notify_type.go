package types

type ReqNotify struct {
	ReturnCode    string `xml:"return_code"`
	ReturnMsg     string `xml:"return_msg"`
	AppId         string `xml:"appid"`
	MchId         string `xml:"mch_id"`
	Nonce         string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
	ResultCode    string `xml:"result_code"`
	ErrCode       string `xml:"err_code"`
	ErrCodeDes    string `xml:"err_code_des"`
	OpenId        string `xml:"openid"`
	IsSubscribe   string `xml:"is_subscribe"`
	TradeType     string `xml:"trade_type"`
	BankType      string `xml:"bank_type"`
	TotalFee      uint32 `xml:"total_fee"`
	CashFee       uint32 `xml:"cash_fee"`
	TransactionId string `xml:"transaction_id"`
	OutRradeNo    string `xml:"out_trade_no"`
	TimeEnd       string `xml:"time_end"`
}

type RspNotify struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}
