package okex

import (
	om "models/okex"
	"server/agent"
	compress "server/gzipcompress"
	process "server/jsonprocess"
	"server/wshb"
	"util/wclient"
)

type AgentPrice struct {
	Agent    wclient.Agent
	Subs     []interface{}
	Process  agent.MsgProcess
	Compress agent.MsgCompress
}

func NewAgentPrice(chanMsgLen uint32) wshb.AgentInstance {

	Process := process.NewJsonProcess()
	Compress := compress.NewMsgGZip()

	return &AgentPrice{
		Process:  Process,
		Compress: Compress,
		Agent:    agent.NewAgent(Compress, Process, chanMsgLen, Handler),
		Subs:     []interface{}{&om.ReqAddChannel{Event: "addChannel", Channel: "btc_forecast_price"}},
	}
}

func (a *AgentPrice) OnInit() {
	a.Agent.SetSubs(a.Subs)
}
func (a *AgentPrice) GetAgent() wclient.Agent {
	return a.Agent
}

func (a *AgentPrice) WriteMsg(msg interface{}) {
	a.Agent.WriteMsg(msg)
}
