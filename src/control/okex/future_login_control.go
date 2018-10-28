package okex

import (
	"fmt"
	om "models/okex"
	"server/wshb"
)

type AgentLogin struct {
	process     wshb.MsgProcess
	compress    wshb.MsgCompress
	chanSendMsg chan interface{}
}

func NewAgentLogin(p wshb.MsgProcess, c wshb.MsgCompress, sendMsgLen int32) wshb.AgentInstance {
	return &AgentLogin{
		process:     p,
		compress:    c,
		chanSendMsg: make(chan interface{}, sendMsgLen),
	}
}

func (a *AgentLogin) Handler(msg interface{}) error {
	var (
		rsps []om.RspFurtureLogin
		err  error
	)

	if err = a.process.UnMarshal(msg.([]byte), &rsps); err != nil {
		return err
	}

	fmt.Println(rsps)

	return nil

}
func (a *AgentLogin) GetSubs() []interface{} {
	return []interface{}{}
}
func (a *AgentLogin) GetWriteMsg() chan interface{} {
	return a.chanSendMsg
}
