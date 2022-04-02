package tcpserver

import (
	"games/comm/utils"
	"games/core/conn"
)

/// <summary>
/// TCPServer TCP服务端
/// <summary>
type TCPServer interface {
	ListenTCP(address string)
	OnConnection(peer conn.Session, reason conn.Reason)
	OnMessage(peer conn.Session, msg interface{}, recvTime utils.Timestamp)
	OnWriteComplete(peer conn.Session)
	Stop()
}
