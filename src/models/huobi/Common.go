package huobi

type SubReq struct {
	Sub    string `json:"sub"`
	Id     string `json:"id"`
	RreqMs int64  `json:"freq-ms"`
}

type Req struct {
	Id  string `json:"id"`
	Req string `json:"req"`
}
