package okex

type PingPang struct {
	Event string `json:"event"`
}

type ReqAddChannel struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}
