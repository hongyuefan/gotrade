package okex

import (
	om "models/okex"
	"server/agent"
	compress "server/gzipcompress"
	process "server/jsonprocess"
	"server/wshb"

	"util/wclient"
)

type AgentTrade struct {
	Agent    wclient.Agent
	Subs     []interface{}
	Process  agent.MsgProcess
	Compress agent.MsgCompress
}

func NewAgentTrade(chanMsgLen uint32) wshb.AgentInstance {

	Process := process.NewJsonProcess()
	Compress := compress.NewMsgGZip()

	return &AgentTrade{
		Process:  Process,
		Compress: Compress,
		Agent:    agent.NewAgent(Compress, Process, chanMsgLen, Handler),
		Subs:     []interface{}{&om.ReqAddChannel{Event: "addChannel", Channel: "ok_sub_futureusd_btc_trade_this_week"}},
	}
}

func (a *AgentTrade) OnInit() {
	a.Agent.SetSubs(a.Subs)
}
func (a *AgentTrade) GetAgent() wclient.Agent {
	return a.Agent
}
func (a *AgentTrade) WriteMsg(msg interface{}) {
	a.Agent.WriteMsg(msg)
}
