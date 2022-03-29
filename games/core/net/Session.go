package net

import "games/core/net/transmit"

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
	SesClient SesType = SesType(0)
	SesServer SesType = SesType(1)
)

type Session interface {
	ID() int64
	Name() string
	IsWebsocket() bool
	Type() SesType
	Conn() interface{}
	SetChannel(channel transmit.IChannel)
	SetContext(key int, val interface{})
	GetContext(key int) interface{}
	Close()
	Write(msg interface{})

	SetOnConnected(cb OnConnected)
	SetOnMessage(cb OnMessage)
	SetOnClosed(cb OnClosed)
	SetOnError(cb OnError)
	SetOnWritten(cb OnWritten)
	SetCloseCallback(cb CloseCallback)
}
