package okex

import (
	"fmt"
	om "models/okex"
	"server/wshb"
)

type AgentPrice struct {
	process     wshb.MsgProcess
	compress    wshb.MsgCompress
	chanSendMsg chan interface{}
}

func NewAgentPrice(p wshb.MsgProcess, c wshb.MsgCompress, sendMsgLen int32) wshb.AgentInstance {
	return &AgentPrice{
		process:     p,
		compress:    c,
		chanSendMsg: make(chan interface{}, sendMsgLen),
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
func (a *AgentPrice) GetWriteMsg() chan interface{} {
	return a.chanSendMsg
}
