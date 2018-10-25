package okex

import (
	"encoding/json"
	"fmt"
	om "models/okex"
	"net"
	"reflect"
	"server/wshb"
	"util/log"
	"util/wclient"
)

type AgentTicker struct {
	gate *wshb.Gate
	conn *wclient.WSConn
}

func NewAgentTicker(conn *wclient.WSConn, gate *wshb.Gate) wclient.Agent {
	return &AgentTicker{conn: conn, gate: gate}
}

func (a *AgentTicker) TickerHandler(msg interface{}) error {

	var pingPang om.PingPang

	if err := json.Unmarshal(msg.([]byte), pingPang); err == nil {
		a.WriteMsg(&om.PingPang{Event: "Pong"})
		return nil
	}

	fmt.Println(string(msg.([]byte)))

	return nil

}
func (a *AgentTicker) Run() {

	var (
		err  error
		data []byte
	)

	a.WriteMsg(&om.ReqComm{Event: "addChannel", Channel: "ok_sub_futureusd_btc_ticker_this_week"})

	for {

		if data, err = a.conn.ReadMsg(); err != nil {
			log.GetLog().LogError("read message: ", err)
			break
		}

		if a.gate.Compress != nil {
			if data, err = a.gate.Compress.UnCompress(data); err != nil {
				log.GetLog().LogError("read message uncompress: ", err)
				break
			}
		}

		if err = a.TickerHandler(data); err != nil {
			log.GetLog().LogError("KlineHandler message error: ", err)
			break
		}

	}
}

func (a *AgentTicker) WriteMsg(msg interface{}) {
	var (
		data []byte
		err  error
	)

	fmt.Println("writeMsg", msg)

	if a.gate.Processor != nil {
		if data, err = a.gate.Processor.Marshal(msg); err != nil {
			log.GetLog().LogError("marshal message ", reflect.TypeOf(msg), " error:", err)
			return
		}
	}

	if err = a.conn.WriteMsg(data); err != nil {
		log.GetLog().LogError("write message ", reflect.TypeOf(msg), "error:", err)
	}
}

func (a *AgentTicker) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *AgentTicker) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *AgentTicker) OnClose() {
	a.conn.Close()
}

func (a *AgentTicker) Destroy() {
	a.conn.Destroy()
}
