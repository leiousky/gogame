package cell

import (
	"games/core/conn"
)

/// <summary>
/// IWorker 业务处理
/// <summary>
type IWorker interface {
	OnInit(args ...interface{})
	OnConnected(peer conn.Session)
	OnClosed(peer conn.Session, reason conn.Reason)
	OnRead(cmd uint32, msg interface{}, peer conn.Session)
	OnCustom(cmd uint32, msg interface{}, peer conn.Session)
	OnTimer(timerID uint32, dt int32, args interface{}) bool
}

/// <summary>
/// IWorkerCreator 业务工厂
/// <summary>
type IWorkerCreator interface {
	Create(c IProc) IWorker
}
