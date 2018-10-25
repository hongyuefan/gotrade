package okex

type ReqTicker struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}
