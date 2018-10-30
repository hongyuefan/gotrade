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

type AgentTicker struct {
	Agent    wclient.Agent
	Subs     []interface{}
	Process  agent.MsgProcess
	Compress agent.MsgCompress
}

func NewAgentTicker(chanMsgLen uint32) wshb.AgentInstance {

	Process := process.NewJsonProcess()
	Compress := compress.NewMsgGZip()

	return &AgentTicker{
		Process:  Process,
		Compress: Compress,
		Agent:    agent.NewAgent(Compress, Process, chanMsgLen),
		Subs:     []interface{}{&om.ReqAddChannel{Event: "addChannel", Channel: "ok_sub_futureusd_btc_ticker_this_week"}},
	}
}

func (a *AgentTicker) OnInit() {
	a.Agent.SetSubs(a.Subs)
}
func (a *AgentTicker) GetAgent() wclient.Agent {
	return a.Agent
}
func (a *AgentTicker) Handler(msg interface{}) {
	var (
		depths []om.RspFurtureTicker
		err    error
	)

	if err = a.Process.UnMarshal(msg.([]byte), &depths); err != nil {
		log.GetLog().LogError("AgentTicker handler error", err)
		return
	}

	fmt.Println(depths)

	return

}
func (a *AgentTicker) WriteMsg(msg interface{}) {
	a.Agent.WriteMsg(msg)
}
