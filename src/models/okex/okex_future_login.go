package okex

type ReqFurtureLogin struct {
	Event  string     `json:"event"`
	Params Parameters `json:"parameters"`
}

type Parameters struct {
	ApiKey     string `json:"api_key"`
	Sign       string `json:"sign"`
	PassPhrase string `json:"passphrase"`
	TimeStamp  string `json:"timestamp"`
}

type RspFurtureLogin struct {
	Data    FurtureLogin `json:"data"`
	Channel string       `json:"channel"`
}

type FurtureLogin struct {
	Result bool `json:"result"`
}
