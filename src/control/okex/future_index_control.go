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

type AgentIndex struct {
	Agent    wclient.Agent
	Subs     []interface{}
	Process  agent.MsgProcess
	Compress agent.MsgCompress
}

func NewAgentIndex(chanMsgLen uint32) wshb.AgentInstance {

	Process := process.NewJsonProcess()
	Compress := compress.NewMsgGZip()

	return &AgentIndex{
		Process:  Process,
		Compress: Compress,
		Agent:    agent.NewAgent(Compress, Process, chanMsgLen),
		Subs:     []interface{}{&om.ReqAddChannel{Event: "addChannel", Channel: "ok_sub_futureusd_btc_index"}},
	}
}

func (a *AgentIndex) OnInit() {
	a.Agent.SetSubs(a.Subs)
}

func (a *AgentIndex) GetAgent() wclient.Agent {
	return a.Agent
}

func (a *AgentIndex) Handler(msg interface{}) {
	var (
		depths []om.RspFurtureIndex
		err    error
	)

	if err = a.Process.UnMarshal(msg.([]byte), &depths); err != nil {
		log.GetLog().LogError("AgentIndex handler error", err)
		return
	}

	fmt.Println(depths)

	return

}
func (a *AgentIndex) WriteMsg(msg interface{}) {
	a.Agent.WriteMsg(msg)
}
