package tcp

import (
	"games/core/conn"
	"time"
)

/// <summary>
/// TCPClient
/// <summary>
type TCPClient interface {
	Session() conn.Session
	Write(msg interface{})
	ConnectTCP(name, address string)
	OnConnected(peer interface{})
	OnMessage(msg interface{}, peer interface{})
	OnClosed(peer interface{})
	Reconnect(d time.Duration)
	Disconnect()
}
