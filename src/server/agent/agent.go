package agent

import (
	"net"
	"reflect"
	"time"
	"util/log"
	"util/wclient"
)

type PingPang struct {
	Event string `json:"event"`
}

type MsgProcess interface {
	UnMarshal([]byte, interface{}) error
	Marshal(interface{}) ([]byte, error)
}
type MsgCompress interface {
	Compress([]byte) ([]byte, error)
	UnCompress([]byte) ([]byte, error)
}

type FuncHandler func(interface{})

type Agent struct {
	chanSign chan bool
	conn     *wclient.WSConn
	compress MsgCompress
	process  MsgProcess
	chanByte chan []byte
	subs     []interface{}
	handler  FuncHandler
}

func NewAgent(compress MsgCompress, process MsgProcess, msgChanLen uint32, funcHandler FuncHandler) wclient.Agent {
	return &Agent{
		chanSign: make(chan bool, 1),
		compress: compress,
		process:  process,
		chanByte: make(chan []byte, msgChanLen),
		handler:  funcHandler,
	}
}

func (a *Agent) SetSubs(subs []interface{}) {
	a.subs = append(a.subs, subs)
}

func (a *Agent) SetCon(wsCon *wclient.WSConn) {
	a.conn = wsCon
}

func (a *Agent) sendSubs() {
	for _, sub := range a.subs {
		a.WriteMsg(sub)
	}
}

func (a *Agent) Run() {

	var (
		err  error
		data []byte
	)

	go a.Ping()

	a.sendSubs()

	for {

		if data, err = a.conn.ReadMsg(); err != nil {
			log.GetLog().LogError("read message: ", err)
			a.chanSign <- true
			break
		}

		if a.compress != nil {
			if data, err = a.compress.UnCompress(data); err != nil {
				log.GetLog().LogError("uncompress error ", err)
				a.chanSign <- true
				break
			}
		}
		a.handler(data)
	}
}

func (a *Agent) ChanMsg() chan []byte {
	return a.chanByte
}

func (a *Agent) WriteMsg(msg interface{}) {
	var (
		data []byte
		err  error
	)

	if a.process != nil {
		if data, err = a.process.Marshal(msg); err != nil {
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
