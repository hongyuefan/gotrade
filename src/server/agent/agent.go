package agent

import (
	"net"
	"reflect"
	"server/wshb"
	"time"
	"util/log"
	"util/wclient"
)

type PingPang struct {
	Event string `json:"event"`
}

type Agent struct {
	chanSign chan bool
	gate     *wshb.Gate
	instance wshb.AgentInstance
	conn     *wclient.WSConn
}

func NewAgent(conn *wclient.WSConn, gate *wshb.Gate, instance wshb.AgentInstance) wclient.Agent {

	return &Agent{
		conn:     conn,
		gate:     gate,
		instance: instance,
		chanSign: make(chan bool, 1),
	}

}

func (a *Agent) subTitles() {
	for _, sub := range a.instance.GetSubs() {
		a.WriteMsg(sub)
	}
}

func (a *Agent) WriteMsgHandler() {
	for {
		select {
		case msg := <-a.instance.GetWriteMsg():
			a.WriteMsg(msg)
		case <-a.chanSign:
			return
		}
	}

}

func (a *Agent) Run() {

	var (
		err      error
		data     []byte
		pingPong PingPang
	)

	go a.Ping()

	go a.WriteMsgHandler()

	a.subTitles()

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

		if err = a.instance.Handler(data); err != nil {
			log.GetLog().LogError("handler message error: ", err)
		}

	}
}

func (a *Agent) WriteMsg(msg interface{}) {
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

func (a *Agent) Ping() {

	ticker := time.NewTicker(time.Second * 30)

	for {
		select {
		case <-ticker.C:
			a.WriteMsg(&PingPang{Event: "ping"})
		case <-a.chanSign:
			ticker.Stop()
			return
		}
	}

}

func (a *Agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *Agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *Agent) OnClose() {
	a.conn.Close()
}

func (a *Agent) Destroy() {
	a.conn.Destroy()
}
