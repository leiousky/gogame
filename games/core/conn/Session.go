package conn

import (
	"games/core/cb"
	"games/core/conn/transmit"
)

const (
	ReasonPeerClosed int32 = iota + 1001 //对端关闭
	ReasonSelfClosed                     //本端关闭
	ReasonSelfExcept                     //本端异常
)

const (
	KConnected    int32 = iota + 1001 //已经连接
	KDisconnected                     //连接关闭
)

type ConnType uint8

const (
	ClientType ConnType = ConnType(0)
	ServerType ConnType = ConnType(1)
)

type Session interface {
	ID() int64
	Name() string
	IsWebsocket() bool
	Type() ConnType
	Conn() interface{}
	SetChannel(channel transmit.IChannel)
	SetContext(key int, val interface{})
	GetContext(key int) interface{}
	Close()
	Write(msg interface{})

	SetOnConnected(cb cb.OnConnected)
	SetOnMessage(cb cb.OnMessage)
	SetOnClosed(cb cb.OnClosed)
	SetOnError(cb cb.OnError)
	SetOnWritten(cb cb.OnWritten)
	SetCloseCallback(cb cb.CloseCallback)
}
