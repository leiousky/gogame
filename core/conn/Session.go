package conn

import (
	"games/core/conn/transmit"
)

type Reason uint8

const (
	KPeerClosed Reason = Reason(0) //对端关闭
	KSelfClosed Reason = Reason(1) //本端关闭
	KSelfExcept Reason = Reason(2) //本端异常
)

type State uint8

const (
	KDisconnected State = State(0)
	KConnected    State = State(1)
)

type Type uint8

const (
	KClient Type = Type(0)
	KServer Type = Type(1)
)

/// <summary>
/// Session 连接会话
/// <summary>
type Session interface {
	ID() int64
	Name() string
	IsWebsocket() bool
	Type() Type
	Connected() bool
	Conn() interface{}
	LocalAddr() string
	RemoteAddr() string
	SetChannel(channel transmit.IChannel)
	SetContext(key int, val interface{})
	GetContext(key int) interface{}
	Close()
	Write(msg interface{})
}
