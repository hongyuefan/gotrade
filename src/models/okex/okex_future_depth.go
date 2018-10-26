package okex

type ReqFurtureDepth struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}

type RspFurtureDepth struct {
	Data    FurtureDepth `json:"data"`
	Channel string       `json:"channel"`
}

type FurtureDepth struct {
	TimeStamp int64       `json:"timestamp"`
	Asks      [][]float32 `json:"asks"`
	Bids      [][]float32 `json:"bids"`
}
