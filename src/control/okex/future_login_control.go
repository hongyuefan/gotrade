package okex

import (
	"fmt"
	om "models/okex"
	"server/agent"
	compress "server/gzipcompress"
	process "server/jsonprocess"
	"server/wshb"
	"util/log"
	"util/wclient"
)

type AgentLogin struct {
	Agent    wclient.Agent
	Subs     []interface{}
	Process  agent.MsgProcess
	Compress agent.MsgCompress
}

func NewAgentLogin(chanMsgLen uint32) wshb.AgentInstance {

	Process := process.NewJsonProcess()
	Compress := compress.NewMsgGZip()

	return &AgentLogin{
		Process:  Process,
		Compress: Compress,
		Agent:    agent.NewAgent(Compress, Process, chanMsgLen),
		Subs:     []interface{}{},
	}
}

func (a *AgentLogin) OnInit() {
	a.Agent.SetSubs(a.Subs)
}

func (a *AgentLogin) GetAgent() wclient.Agent {
	return a.Agent
}
func (a *AgentLogin) Handler(msg interface{}) {
	var (
		depths []om.RspFurtureLogin
		err    error
	)

	if err = a.Process.UnMarshal(msg.([]byte), &depths); err != nil {
		log.GetLog().LogError("AgentLogin handler error", err)
		return
	}

	fmt.Println(depths)

	return

}
func (a *AgentLogin) WriteMsg(msg interface{}) {
	a.Agent.WriteMsg(msg)
}
