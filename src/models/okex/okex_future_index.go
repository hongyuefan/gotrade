package okex

type ReqFurtureIndex struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}

type RspFurtureIndex struct {
	Data    FurtureIndex `json:"data"`
	Channel string       `json:"channel"`
}

type FurtureIndex struct {
	TimeStamp   string `json:"timestamp"`
	FutureIndex string `json:"futureIndex"`
}
