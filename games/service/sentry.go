package service

import (
	"games/core"
	"games/core/conn"
)

type Sentry struct {
	c    core.IProc
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

func (s *Sentry) OnConnected(peer conn.Session, Type conn.ConnType) {

}

func (s *Sentry) OnClosed(peer conn.Session, Type conn.ConnType) {

}

func (s *Sentry) OnRead(cmd uint32, msg interface{}, session conn.Session) {
	s.main.OnRead(cmd, msg, session)
}

func (s *Sentry) OnCustom(cmd uint32, msg interface{}, session conn.Session) {
	s.main.OnCustom(cmd, msg, session)
}

func (s *Sentry) OnTimer(timerID uint32, dt int32, args interface{}) bool {
	return true
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
