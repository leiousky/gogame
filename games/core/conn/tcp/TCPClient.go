package tcp

import (
	"games/comm/utils"
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
	OnConnection(peer conn.Session)
	OnMessage(peer conn.Session, msg interface{}, recvTime utils.Timestamp)
	OnWriteComplete(peer conn.Session)
	Reconnect(d time.Duration)
	Disconnect()
}
