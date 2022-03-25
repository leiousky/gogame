package core

import "games/core/net"

/// <summary>
/// IWorker 业务处理
/// <summary>
type IWorker interface {
	OnInit(args ...interface{})
	OnConnected(peer net.Session, Type net.SesType)
	OnClosed(peer net.Session, Type net.SesType)
	OnMessage(cmd uint32, msg interface{}, peer net.Session)

	RunAfter(delay int32, args interface{}) uint32
	RunAfterWith(delay int32, handler net.TimerCallback, args interface{}) uint32
	RunEvery(delay, interval int32, args interface{}) uint32
	RunEveryWith(delay, interval int32, handler net.TimerCallback, args interface{}) uint32
	RemoveTimer(timerID uint32)

	SetProc(c IProc)
	GetProc() IProc

	SetDispatcher(c IProc)
	GetDispatcher() IProc
	ResetDispatcher()
}

type IWorkerCreator interface {
	Create(c IProc) IWorker
}
