package service

import (
	"fmt"
	"games/comm/utils"
	"games/core"
)

type sMain struct {
	entry    *Sentry
	handlers core.CmdCallbacks
	tick     uint32
}

func newMain(s *Sentry) *sMain {
	return &sMain{entry: s,
		handlers: core.CmdCallbacks{}}
}

func (s *sMain) initModules(args ...interface{}) {
	//CreateTimer
	//CancelTimer
	//CreatCronFunc
	//DelCronFunc
	//timer := s.entry.GetProc().GetTimer().(*timer.SafeTimerScheduel)
	//s.tick, _ = timer.CreatCronFunc("@every 1s", func() {
	//	fmt.Println("OnTick ...")
	//})
}

func (s *sMain) OnTick() {
	fmt.Println("sMain::OnTick 机器人入桌检查...")
}

func (s *sMain) OnMessage(cmd uint32, msg interface{}, peer core.Session) {
	utils.SafeCall(func() {
		if handler, ok := s.handlers[cmd]; ok {
			handler(msg, peer)
		} else {
		}
	})
}
