package okex

type ReqComm struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}
