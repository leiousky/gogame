package db

import (
	"games/core/cell"
	"games/core/conn"
)

type Sentry struct {
	c    cell.IProc
	main *sMain
}

func newEntry(c cell.IProc) cell.IWorker {
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

func (s *Sentry) OnConnected(peer conn.Session) {

}

func (s *Sentry) OnClosed(peer conn.Session, reason conn.Reason) {

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
	cell.IWorkerCreator
}

func (s *SentryCreator) Create(c cell.IProc) cell.IWorker {
	return newEntry(c)
}

func NewSentryCreator() cell.IWorkerCreator {
	return &SentryCreator{}
}
