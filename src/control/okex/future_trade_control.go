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

type AgentTrade struct {
	gate     *wshb.Gate
	conn     *wclient.WSConn
	chanSign chan bool
}

func NewAgentTrade(conn *wclient.WSConn, gate *wshb.Gate) wclient.Agent {
	return &AgentTrade{conn: conn, gate: gate, chanSign: make(chan bool, 1)}
}

func (a *AgentTrade) TradeHandler(msg interface{}) error {
	var (
		tickers []om.RspFurtureTrade
		err     error
	)

	if err = a.gate.Processor.UnMarshal(msg.([]byte), &tickers); err != nil {
		return err
	}

	fmt.Println(tickers)

	return nil

}
func (a *AgentTrade) Run() {

	var (
		err      error
		data     []byte
		pingPong om.PingPang
	)

	go a.Ping()

	a.WriteMsg(&om.ReqFurtureTrade{Event: "addChannel", Channel: "ok_sub_futureusd_btc_trade_this_week"})

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

		if err = a.gate.Processor.UnMarshal(data, &pingPong); err == nil {
			continue
		}

		if err = a.TradeHandler(data); err != nil {
			log.GetLog().LogError("TradeHandler message error: ", err)
		}

	}
}

func (a *AgentTrade) WriteMsg(msg interface{}) {
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

func (a *AgentTrade) Ping() {

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

func (a *AgentTrade) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *AgentTrade) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *AgentTrade) OnClose() {
	a.conn.Close()
}

func (a *AgentTrade) Destroy() {
	a.conn.Destroy()
}
