package wclient

type Agent interface {
	Run()
	WriteMsg(interface{})
	OnClose()
}
