package okex

import (
	"fmt"
	om "models/okex"
	"server/agent"
	compress "server/gzipcompress"
	process "server/jsonprocess"
	"server/wshb"
	"util/wclient"
)

type AgentDepth struct {
	Agent    wclient.Agent
	Subs     []interface{}
	Process  agent.MsgProcess
	Compress agent.MsgCompress
}

func NewAgentDepth(chanMsgLen uint32) wshb.AgentInstance {

	Process := process.NewJsonProcess()
	Compress := compress.NewMsgGZip()

	return &AgentDepth{
		Process:  Process,
		Compress: Compress,
		Agent:    agent.NewAgent(Compress, Process, chanMsgLen, Handler),
		Subs:     []interface{}{&om.ReqAddChannel{Event: "addChannel", Channel: "ok_sub_futureusd_btc_depth_this_week"}},
	}
}

func (a *AgentDepth) OnInit() {
	a.Agent.SetSubs(a.Subs)
}

func (a *AgentDepth) GetAgent() wclient.Agent {
	return a.Agent
}

func Handler(msg interface{}) {

	//	var (
	//		depths []om.RspFurtureDepth
	//		err    error
	//	)

	//	if err = a.Process.UnMarshal(msg.([]byte), &depths); err != nil {
	//		log.GetLog().LogError("agentDepth handler error", err)
	//		return
	//	}

	fmt.Println(string(msg.([]byte)))

	return
}
func (a *AgentDepth) WriteMsg(msg interface{}) {
	a.Agent.WriteMsg(msg)
}
