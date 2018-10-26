package okex

type ReqFurtureTicker struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}

type RspFurtureTicker struct {
	Data    FurtureTicker `json:"data"`
	Channel string        `json:"channel"`
}

type FurtureTicker struct {
	LimitHigh  string `json:"limitHigh"`
	Vol        string `json:"vol"`
	Last       string `json:"last"`
	Sell       string `json:"sell"`
	Buy        string `json:"buy"`
	UnitAmount string `json:"unitAmount"`
	HoldAmount string `json:"hold_amount"`
	ContractId int64  `json:"contractId"`
	High       string `json:"high"`
	Low        string `json:"low"`
	LimitLow   string `json:"limitLow"`
}
