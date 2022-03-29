package tcp

import (
	"games/core/conn"
	"time"
)

/// <summary>
/// TCPClient TCP客户端
/// <summary>
type TCPClient interface {
	Session() conn.Session
	Write(msg interface{})
	ConnectTCP(name, address string)
	OnConnected(peer conn.Session)
	OnMessage(msg interface{}, peer conn.Session)
	OnClosed(peer conn.Session)
	Reconnect(d time.Duration)
	Disconnect()
}
