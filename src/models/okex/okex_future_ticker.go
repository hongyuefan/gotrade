package okex

type ReqFurtureTicker struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}

type RspFurtureTicker struct {
	Data    *FurtureTicker `json:"data"`
	Channel string         `json:"channel"`
	Binary  int32          `json:"binary"`
}

type FurtureTicker struct {
	LimitHigh  string  `json:"limitHigh"`
	Vol        int32   `json:"vol"`
	Last       float32 `json:"last"`
	Sell       float32 `json:"sell"`
	Buy        float32 `json:"buy"`
	UnitAmount int32   `json:"unitAmount"`
	HoldAmount int32   `json:"hold_amount"`
	ContractId int64   `json:"contractId"`
	High       float32 `json:"high"`
	Low        float32 `json:"low"`
	LimitLow   string  `json:"limitLow"`
}
