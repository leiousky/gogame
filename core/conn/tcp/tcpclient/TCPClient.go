package tcpclient

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
	ConnectTCP(address string)
	OnConnection(peer conn.Session, reason conn.Reason)
	OnMessage(peer conn.Session, msg interface{}, recvTime utils.Timestamp)
	OnWriteComplete(peer conn.Session)
	Reconnect(d time.Duration)
	Disconnect()
}

/// <summary>
/// sessions 客户端容器
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
