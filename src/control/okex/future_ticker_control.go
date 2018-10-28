package okex

import (
	"fmt"
	om "models/okex"
	"server/wshb"
)

type AgentTicker struct {
	process     wshb.MsgProcess
	compress    wshb.MsgCompress
	chanSendMsg chan interface{}
}

func NewAgentTicker(p wshb.MsgProcess, c wshb.MsgCompress, sendMsgLen int32) wshb.AgentInstance {
	return &AgentTicker{
		process:     p,
		compress:    c,
		chanSendMsg: make(chan interface{}, sendMsgLen),
	}
}

func (a *AgentTicker) Handler(msg interface{}) error {
	var (
		rsps []om.RspFurtureTicker
		err  error
	)

	if err = a.process.UnMarshal(msg.([]byte), &rsps); err != nil {
		return err
	}

	fmt.Println(rsps)

	return nil

}
func (a *AgentTicker) GetSubs() []interface{} {
	return []interface{}{&om.ReqAddChannel{Event: "addChannel", Channel: "ok_sub_futureusd_btc_ticker_this_week"}}
}
func (a *AgentTicker) GetWriteMsg() chan interface{} {
	return a.chanSendMsg
}
