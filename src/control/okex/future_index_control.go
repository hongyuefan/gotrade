package okex

import (
	"fmt"
	om "models/okex"
	"server/wshb"
)

type AgentIndex struct {
	process  wshb.MsgProcess
	compress wshb.MsgCompress
}

func NewAgentIndex(p wshb.MsgProcess, c wshb.MsgCompress) wshb.AgentInstance {
	return &AgentIndex{
		process:  p,
		compress: c,
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
