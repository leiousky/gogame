package conn

import (
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

/// <summary>
/// Session 连接会话
/// <summary>
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
}
