package wclient

type Agent interface {
	Run()
	WriteMsg(interface{})
	OnClose()
	SetCon(*WSConn)
	ChanMsg() chan []byte
	SetSubs([]interface{})
}
