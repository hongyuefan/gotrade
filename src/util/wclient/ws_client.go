package wclient

import (
	"sync"
	"time"
	"util/log"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	sync.Mutex
	Addr             string
	ConnNum          int
	ConnectInterval  time.Duration
	PendingWriteNum  int
	MaxMsgLen        uint32
	HandshakeTimeout time.Duration
	AutoReconnect    bool
	NewAgent         func(*WSConn) Agent
	dialer           websocket.Dialer
	conns            WebsocketConnSet
	wg               sync.WaitGroup
	closeFlag        bool
}

func (client *WSClient) Start() {

	client.init()

	for i := 0; i < client.ConnNum; i++ {
		client.wg.Add(1)
		go client.connect()
	}
}

func (client *WSClient) init() {
	client.Lock()
	defer client.Unlock()

	if client.ConnNum <= 0 {
		client.ConnNum = 1
		log.GetLog().LogDebug("invalid ConnNum, reset to", client.ConnNum)
	}
	if client.ConnectInterval <= 0 {
		client.ConnectInterval = 3 * time.Second
		log.GetLog().LogDebug("invalid ConnectInterval, reset to ", client.ConnectInterval)
	}
	if client.PendingWriteNum <= 0 {
		client.PendingWriteNum = 100
		log.GetLog().LogDebug("invalid PendingWriteNum, reset to", client.PendingWriteNum)
	}
	if client.MaxMsgLen <= 0 {
		client.MaxMsgLen = 4096
		log.GetLog().LogDebug("invalid MaxMsgLen, reset to", client.MaxMsgLen)
	}
	if client.HandshakeTimeout <= 0 {
		client.HandshakeTimeout = 10 * time.Second
		log.GetLog().LogDebug("invalid HandshakeTimeout, reset to", client.HandshakeTimeout)
	}
	if client.NewAgent == nil {
		log.GetLog().LogDebug("NewAgent must not be nil")
	}
	if client.conns != nil {
		log.GetLog().LogDebug("client is running")
	}

	client.conns = make(WebsocketConnSet)
	client.closeFlag = false
	client.dialer = websocket.Dialer{
		HandshakeTimeout: client.HandshakeTimeout,
	}
}

func (client *WSClient) dial() *websocket.Conn {
	for {
		conn, _, err := client.dialer.Dial(client.Addr, nil)
		if err == nil || client.closeFlag {
			return conn
		}

		log.GetLog().LogDebug("connect to", client.Addr, "error:", err)
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *WSClient) connect() {
	defer client.wg.Done()

reconnect:
	conn := client.dial()
	if conn == nil {
		return
	}
	conn.SetReadLimit(int64(client.MaxMsgLen))

	client.Lock()
	if client.closeFlag {
		client.Unlock()
		conn.Close()
		return
	}
	client.conns[conn] = struct{}{}
	client.Unlock()

	wsConn := newWSConn(conn, client.PendingWriteNum, client.MaxMsgLen)
	agent := client.NewAgent(wsConn)
	agent.Run()

	// cleanup
	wsConn.Close()
	client.Lock()
	delete(client.conns, conn)
	client.Unlock()
	agent.OnClose()

	if client.AutoReconnect {
		time.Sleep(client.ConnectInterval)
		goto reconnect
	}
}

func (client *WSClient) Close() {
	client.Lock()
	client.closeFlag = true
	for conn := range client.conns {
		conn.Close()
	}
	client.conns = nil
	client.Unlock()

	client.wg.Wait()
}
