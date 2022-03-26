package net

const (
	ReasonPeerClosed int32 = iota + 1001 //对端关闭
	ReasonSelfClosed                     //本端关闭
	ReasonSelfExcept                     //本端异常
)

const (
	KConnected    int32 = iota + 1001 //已经连接
	KDisconnected                     //连接关闭
)

type SesType uint8

const (
	SesClient SesType = SesType(1)
	SesServer SesType = SesType(2)
)

type Session interface {
	ID() int64
	IsWebsocket() bool
	Type() SesType
	Conn() interface{}
	SetContext(key int, val interface{})
	GetContext(key int) interface{}
	Close()
	Write(msg interface{})

	SetOnConnected(cb OnConnected)
	SetOnClosed(cb OnClosed)
	SetOnMessage(cb OnMessage)
	SetOnError(cb OnError)
	SetOnWritten(cb OnWritten)
	SetCloseCallback(cb CloseCallback)
}
