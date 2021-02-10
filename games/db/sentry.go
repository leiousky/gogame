package db

import "games/core"

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

func (s *Sentry) OnMessage(cmd uint32, msg interface{}, session core.Session) {
	s.main.OnMessage(cmd, msg, session)
}

//func (s *Sentry) SetProc(c core.IProc) {
//	s.c = c
//}

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
