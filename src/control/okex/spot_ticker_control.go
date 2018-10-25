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

type AgentSpotTicker struct {
	gate *wshb.Gate
	conn *wclient.WSConn
}

func NewAgentSpotTicker(conn *wclient.WSConn, gate *wshb.Gate) wclient.Agent {
	return &AgentSpotTicker{conn: conn, gate: gate}
}

func (a *AgentSpotTicker) TickerHandler(msg interface{}) error {

	var pingPang om.PingPang

	if err := json.Unmarshal(msg.([]byte), pingPang); err != nil {
		a.WriteMsg(&om.PingPang{Event: "Pong"})
		return nil
	}

	return nil

}
func (a *AgentSpotTicker) Run() {

	var (
		err  error
		data []byte
	)

	a.WriteMsg(&om.ReqComm{Event: "addChannel", Channel: "ok_sub_spot_bch_btc_ticker"})

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

		fmt.Println(string(data), a.RemoteAddr())

		if err = a.TickerHandler(data); err != nil {
			log.GetLog().LogError("KlineHandler message error: ", err)
			break
		}

	}
}

func (a *AgentSpotTicker) WriteMsg(msg interface{}) {
	var (
		data []byte
		err  error
	)
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

func (a *AgentSpotTicker) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *AgentSpotTicker) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *AgentSpotTicker) OnClose() {
	a.conn.Close()
}

func (a *AgentSpotTicker) Destroy() {
	a.conn.Destroy()
}
