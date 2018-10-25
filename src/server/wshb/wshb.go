package wshb

import (
	"time"

	"util/wclient"
)

type MsgProcess interface {
	UnMarshal([]byte, interface{}) error
	Marshal(interface{}) ([]byte, error)
	Route(interface{}, wclient.Agent) error
}

type MsgCompress interface {
	Compress([]byte) ([]byte, error)
	UnCompress([]byte) ([]byte, error)
}

type Gate struct {
	Addr             string
	ConnNum          int
	ConnectInterval  time.Duration
	PendingWriteNum  int
	MaxMsgLen        uint32
	HandshakeTimeout time.Duration
	AutoReconnect    bool
	NewAgent         func(*wclient.WSConn) wclient.Agent

	Processor MsgProcess
	Compress  MsgCompress
}

type FuncNewAgent func(*wclient.WSConn, *Gate) wclient.Agent

func NewGate(addr string, conNum, writeNum int, maxMsgLen uint32, conInterval, handshakeTimeout time.Duration, autoReconect bool, process MsgProcess, compress MsgCompress) *Gate {

	gate := new(Gate)
	gate.Compress = compress
	gate.Processor = process
	gate.Addr = addr
	gate.ConnNum = conNum
	gate.ConnectInterval = conInterval
	gate.PendingWriteNum = writeNum
	gate.MaxMsgLen = maxMsgLen
	gate.HandshakeTimeout = handshakeTimeout
	gate.AutoReconnect = autoReconect

	return gate
}

func (gate *Gate) Run(closeSig chan bool, funcAgent FuncNewAgent) {

	wc := new(wclient.WSClient)
	wc.Addr = gate.Addr
	wc.AutoReconnect = gate.AutoReconnect
	wc.ConnectInterval = gate.ConnectInterval
	wc.ConnNum = gate.ConnNum
	wc.HandshakeTimeout = gate.HandshakeTimeout
	wc.MaxMsgLen = gate.MaxMsgLen
	wc.PendingWriteNum = gate.PendingWriteNum

	wc.NewAgent = func(conn *wclient.WSConn) wclient.Agent {
		return funcAgent(conn, gate)
	}

	wc.Start()
	<-closeSig
	wc.Close()
}

func (gate *Gate) OnDestroy() {}
