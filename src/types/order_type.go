package types

type ReqGetOrder struct {
	Token         string `json:"token"`
	Page          int64  `json:"page"`
	PerPage       int64  `json:"perpage"`
	Offset        int64  `json:"offset"`
	Sort          int64  `json:"sort"`
	AppId         string `json:"appid"`
	ResultCode    string `json:"result_code"`
	TransactionId string `json:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no"`
	Time          int64  `json:"time"`
}

type RspOrders struct {
	AppId         string `json:"appid"`
	ReturnCode    string `json:"return_code"`
	ReturnMsg     string `json:"return_msg"`
	Sign          string `orm:"column(sign)"`
	ResultCode    string `orm:"column(result_code)"`
	ErrCode       string `orm:"column(err_code)"`
	ErrCodeDes    string `orm:"column(err_code_des)"`
	TotalFee      uint32 `orm:"column(total_fee)"`
	TransactionId string `orm:"column(transaction_id)"`
	OutTradeNo    string `orm:"column(out_trade_no)"`
	TimeEnd       string `orm:"column(time_end)"`
}
