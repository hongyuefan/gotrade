package okex

import (
	"fmt"
	om "models/okex"
	"server/wshb"
)

type AgentTrade struct {
	process     wshb.MsgProcess
	compress    wshb.MsgCompress
	chanSendMsg chan interface{}
}

func NewAgentTrade(p wshb.MsgProcess, c wshb.MsgCompress, sendMsgLen int32) wshb.AgentInstance {
	return &AgentTrade{
		process:     p,
		compress:    c,
		chanSendMsg: make(chan interface{}, sendMsgLen),
	}
}

func (a *AgentTrade) Handler(msg interface{}) error {
	var (
		rsps []om.RspFurtureTrade
		err  error
	)

	if err = a.process.UnMarshal(msg.([]byte), &rsps); err != nil {
		return err
	}

	fmt.Println(rsps)

	return nil

}
func (a *AgentTrade) GetSubs() []interface{} {
	return []interface{}{&om.ReqAddChannel{Event: "addChannel", Channel: "ok_sub_futureusd_btc_trade_this_week"}}
}
func (a *AgentTrade) GetWriteMsg() chan interface{} {
	return a.chanSendMsg
}
