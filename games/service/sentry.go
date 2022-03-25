package service

import (
	"games/core"
	"games/core/net"
)

type Sentry struct {
	c    core.IProc
	d    core.IProc
	main *sMain
}

func newEntry(c core.IProc) core.IWorker {
	p := &Sentry{c: c}
	p.main = newMain(p)
	return p
}

func (s *Sentry) OnInit(args ...interface{}) {
	s.main.initModules(args...)
}

func (s *Sentry) OnTick() {
	s.main.OnTick()
}

func (s *Sentry) OnConnected(peer net.Session, Type net.SesType) {

}

func (s *Sentry) OnClosed(peer net.Session, Type net.SesType) {

}

func (s *Sentry) OnMessage(cmd uint32, msg interface{}, session net.Session) {
	s.main.OnMessage(cmd, msg, session)
}

func (s *Sentry) RunAfter(delay int32, args interface{}) uint32 {
	return 0
}

func (s *Sentry) RunAfterWith(delay int32, handler net.TimerCallback, args interface{}) uint32 {
	return 0
}

func (s *Sentry) RunEvery(delay, interval int32, args interface{}) uint32 {
	return 0
}
func (s *Sentry) RunEveryWith(delay, interval int32, handler net.TimerCallback, args interface{}) uint32 {
	return 0
}

func (s *Sentry) RemoveTimer(timerID uint32) {

}

func (s *Sentry) SetProc(c core.IProc) {
	s.c = c
}

func (s *Sentry) GetProc() core.IProc {
	return s.c
}

func (s *Sentry) SetDispatcher(c core.IProc) {
	s.d = c
}

func (s *Sentry) GetDispatcher() core.IProc {
	return s.d
}

func (s *Sentry) ResetDispatcher() {
	s.d = nil
}

type SentryCreator struct {
	core.IWorkerCreator
}

func (s *SentryCreator) Create(c core.IProc) core.IWorker {
	return newEntry(c)
}

func NewSentryCreator() core.IWorkerCreator {
	return &SentryCreator{}
}
