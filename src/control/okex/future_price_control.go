package okex

import (
	"fmt"
	om "models/okex"
	"server/wshb"
)

type AgentPrice struct {
	process  wshb.MsgProcess
	compress wshb.MsgCompress
}

func NewAgentPrice(p wshb.MsgProcess, c wshb.MsgCompress) wshb.AgentInstance {
	return &AgentPrice{
		process:  p,
		compress: c,
	}
}

func (a *AgentPrice) Handler(msg interface{}) error {
	var (
		rsps []om.RspFurturePrice
		err  error
	)

	if err = a.process.UnMarshal(msg.([]byte), &rsps); err != nil {
		return err
	}

	fmt.Println(rsps)

	return nil

}
func (a *AgentPrice) GetSubs() []interface{} {
	return []interface{}{&om.ReqAddChannel{Event: "addChannel", Channel: "btc_forecast_price"}}
}
