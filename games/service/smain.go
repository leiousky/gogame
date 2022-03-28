package service

import (
	"fmt"
	"games/comm/utils"
	"games/core/net"
)

type sMain struct {
	entry    *Sentry
	handlers net.CmdCallbacks
	tick     uint32
}

func newMain(s *Sentry) *sMain {
	return &sMain{entry: s,
		handlers: net.CmdCallbacks{}}
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

func (s *sMain) OnRead(cmd uint32, msg interface{}, peer net.Session) {
	utils.SafeCall(func() {
		if handler, ok := s.handlers[cmd]; ok {
			handler(msg, peer)
		} else {
		}
	})
}

func (s *sMain) OnCustom(cmd uint32, msg interface{}, peer net.Session) {

}
