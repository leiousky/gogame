package db

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
	fmt.Println("数据库初始化")
}

func (s *sMain) OnTick() {
}

func (s *sMain) OnMessage(cmd uint32, msg interface{}, peer core.Session) {
	utils.SafeCall(func() {
		if handler, ok := s.handlers[cmd]; ok {
			handler(msg, peer)
		} else {
		}
	})
}
