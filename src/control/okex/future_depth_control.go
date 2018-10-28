package okex

import (
	"fmt"
	om "models/okex"
	"server/wshb"
)

type AgentDepth struct {
	process     wshb.MsgProcess
	compress    wshb.MsgCompress
	chanSendMsg chan interface{}
}

func NewAgentDepth(p wshb.MsgProcess, c wshb.MsgCompress, sendMsgLen int32) wshb.AgentInstance {
	return &AgentDepth{
		process:     p,
		compress:    c,
		chanSendMsg: make(chan interface{}, sendMsgLen),
	}
}

func (a *AgentDepth) Handler(msg interface{}) error {
	var (
		depths []om.RspFurtureDepth
		err    error
	)

	if err = a.process.UnMarshal(msg.([]byte), &depths); err != nil {
		return err
	}

	fmt.Println(depths)

	return nil

}
func (a *AgentDepth) GetSubs() []interface{} {
	return []interface{}{&om.ReqAddChannel{Event: "addChannel", Channel: "ok_sub_futureusd_btc_depth_this_week"}}
}
func (a *AgentDepth) GetWriteMsg() chan interface{} {
	return a.chanSendMsg
}
