package okex

type ReqFurtureTrade struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}

type RspFurtureTrade struct {
	Data    [][]string `json:"data"`
	Channel string     `json:"channel"`
}
