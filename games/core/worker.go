package core

import (
	"games/core/net"
)

/// <summary>
/// IWorker 业务处理
/// <summary>
type IWorker interface {
	OnInit(args ...interface{})
	OnConnected(peer net.Session, Type net.SesType)
	OnClosed(peer net.Session, Type net.SesType)
	OnRead(cmd uint32, msg interface{}, peer net.Session)
	OnCustom(cmd uint32, msg interface{}, peer net.Session)
	OnTimer(timerID uint32, dt int32, args interface{}) bool
}

type IWorkerCreator interface {
	Create(c IProc) IWorker
}
