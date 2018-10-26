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

type AgentDepth struct {
	gate     *wshb.Gate
	conn     *wclient.WSConn
	chanSign chan bool
}

func NewAgentDepth(conn *wclient.WSConn, gate *wshb.Gate) wclient.Agent {
	return &AgentDepth{conn: conn, gate: gate, chanSign: make(chan bool, 1)}
}

func (a *AgentDepth) DepthHandler(msg interface{}) error {
	var (
		depths []om.RspFurtureDepth
		err    error
	)

	if err = a.gate.Processor.UnMarshal(msg.([]byte), &depths); err != nil {
		return err
	}

	fmt.Println(depths)

	return nil

}
func (a *AgentDepth) Run() {

	var (
		err      error
		data     []byte
		pingPong om.PingPang
	)

	go a.Ping()

	a.WriteMsg(&om.ReqFurtureTicker{Event: "addChannel", Channel: "ok_sub_futureusd_btc_depth_this_week"})

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

		if err = a.DepthHandler(data); err != nil {
			log.GetLog().LogError("KlineHandler message error: ", err)
		}

	}
}

func (a *AgentDepth) WriteMsg(msg interface{}) {
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

func (a *AgentDepth) Ping() {

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

func (a *AgentDepth) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *AgentDepth) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *AgentDepth) OnClose() {
	a.conn.Close()
}

func (a *AgentDepth) Destroy() {
	a.conn.Destroy()
}
