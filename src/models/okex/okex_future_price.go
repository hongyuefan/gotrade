package okex

type RspFurturePrice struct {
	Data      string `json:"data"`
	Channel   string `json:"channel"`
	TimeStamp string `json:"timestamp"`
}
