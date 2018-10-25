package huobi

import (
	"encoding/json"
	"fmt"
	"models"
	"net"
	"reflect"
	"server/wshb"
	"util/log"
	"util/wclient"
)

type AgentKline struct {
	gate *wshb.Gate
	conn *wclient.WSConn
}

func NewAgentKline(conn *wclient.WSConn, gate *wshb.Gate) wclient.Agent {
	return &AgentKline{conn: conn, gate: gate}
}

func (a *AgentKline) KlineHandler(msg interface{}) error {

	var ping models.Ping

	if err := json.Unmarshal(msg.([]byte), ping); err != nil {
		a.WriteMsg(&models.Pong{Pong: ping.Ping})
		return nil
	}

	var klineData models.KLineData

	if err := json.Unmarshal(msg.([]byte), &klineData); err != nil {
		return err
	}

	fmt.Println(klineData)

	return nil

}
func (a *AgentKline) Run() {

	var (
		err  error
		data []byte
	)

	a.WriteMsg(&models.SubReq{Sub: "market.btcusdt.kline.1min", RreqMs: 5000})

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

		if err = a.KlineHandler(data); err != nil {
			log.GetLog().LogError("KlineHandler message error: ", err)
			break
		}

	}
}

func (a *AgentKline) WriteMsg(msg interface{}) {
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

func (a *AgentKline) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *AgentKline) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *AgentKline) OnClose() {
	a.conn.Close()
}

func (a *AgentKline) Destroy() {
	a.conn.Destroy()
}
