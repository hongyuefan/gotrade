package wshb

import (
	"time"
	"util/wclient"
)

type AgentInstance interface {
	WriteMsg(interface{})
	OnInit()
	GetAgent() wclient.Agent
}

type Gate struct {
	Addr             string
	ConnNum          int
	ConnectInterval  time.Duration
	PendingWriteNum  int
	MaxMsgLen        uint32
	HandshakeTimeout time.Duration
	AutoReconnect    bool
	Agent            AgentInstance
}

func NewGate(addr string, conNum, writeNum int, maxMsgLen uint32, conInterval, handshakeTimeout time.Duration, autoReconect bool, agentInstance AgentInstance) *Gate {

	gate := new(Gate)
	gate.Addr = addr
	gate.ConnNum = conNum
	gate.ConnectInterval = conInterval
	gate.PendingWriteNum = writeNum
	gate.MaxMsgLen = maxMsgLen
	gate.HandshakeTimeout = handshakeTimeout
	gate.AutoReconnect = autoReconect
	gate.Agent = agentInstance

	return gate
}

func (gate *Gate) Run(closeSig chan bool) {

	gate.Agent.OnInit()

	wc := new(wclient.WSClient)
	wc.Addr = gate.Addr
	wc.AutoReconnect = gate.AutoReconnect
	wc.ConnectInterval = gate.ConnectInterval
	wc.ConnNum = gate.ConnNum
	wc.HandshakeTimeout = gate.HandshakeTimeout
	wc.MaxMsgLen = gate.MaxMsgLen
	wc.PendingWriteNum = gate.PendingWriteNum
	wc.WAgent = gate.Agent.GetAgent()
	wc.Start()
	<-closeSig
	wc.Close()
}

func (gate *Gate) OnDestroy() {}
