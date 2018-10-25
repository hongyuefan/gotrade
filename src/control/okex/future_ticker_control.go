package okex

import (
	"fmt"
	om "models/okex"
	"net"
	"reflect"
	"server/wshb"
	"time"
	"util/log"
	"util/wclient"
)

type AgentTicker struct {
	gate     *wshb.Gate
	conn     *wclient.WSConn
	chanSign chan bool
}

func NewAgentTicker(conn *wclient.WSConn, gate *wshb.Gate) wclient.Agent {
	return &AgentTicker{conn: conn, gate: gate, chanSign: make(chan bool, 1)}
}

func (a *AgentTicker) TickerHandler(msg interface{}) error {
	return nil

}
func (a *AgentTicker) Run() {

	var (
		err  error
		data []byte
	)

	go a.Ping()

	a.WriteMsg(&om.ReqComm{Event: "addChannel", Channel: "ok_sub_futureusd_btc_ticker_this_week"})

	for {

		if data, err = a.conn.ReadMsg(); err != nil {
			log.GetLog().LogError("read message: ", err)
			a.chanSign <- true
			break
		}

		if a.gate.Compress != nil {
			if data, err = a.gate.Compress.UnCompress(data); err != nil {
				log.GetLog().LogError("read message uncompress: ", err)
				a.chanSign <- true
				break
			}
		}

		fmt.Println(string(data))

		if err = a.TickerHandler(data); err != nil {
			log.GetLog().LogError("KlineHandler message error: ", err)
			a.chanSign <- true
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

func (a *AgentTicker) Ping() {

	ticker := time.NewTicker(time.Second * 30)

	for {
		select {
		case <-ticker.C:
			a.WriteMsg(&om.PingPang{Event: "ping"})
		case <-a.chanSign:
			ticker.Stop()
			return
		}
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
