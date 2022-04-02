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

/// <summary>
/// sessions 服务端容器
/// <summary>
var sessions = conn.NewSessions()

func Get(id int64) conn.Session {
	return sessions.Get(id)
}

func Count() int {
	return sessions.Count()
}

func CloseAll() {
	sessions.CloseAll()
}

func Wait() {
	sessions.Wait()
}

func Stop() {
	sessions.Stop()
}
