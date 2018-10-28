package okex

import (
	"fmt"
	om "models/okex"
	"server/wshb"
)

type AgentIndex struct {
	process     wshb.MsgProcess
	compress    wshb.MsgCompress
	chanSendMsg chan interface{}
}

func NewAgentIndex(p wshb.MsgProcess, c wshb.MsgCompress, sendMsgLen int32) wshb.AgentInstance {
	return &AgentIndex{
		process:     p,
		compress:    c,
		chanSendMsg: make(chan interface{}, sendMsgLen),
	}
}

func (a *AgentIndex) Handler(msg interface{}) error {
	var (
		indexs []om.RspFurtureIndex
		err    error
	)

	if err = a.process.UnMarshal(msg.([]byte), &indexs); err != nil {
		return err
	}

	fmt.Println(indexs)

	return nil

}
func (a *AgentIndex) GetSubs() []interface{} {
	return []interface{}{&om.ReqAddChannel{Event: "addChannel", Channel: "ok_sub_futureusd_btc_index"}}
}
func (a *AgentIndex) GetWriteMsg() chan interface{} {
	return a.chanSendMsg
}
