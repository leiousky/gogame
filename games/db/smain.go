package db

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
	fmt.Println("数据库初始化")
}

func (s *sMain) OnTick() {
}

func (s *sMain) OnMessage(cmd uint32, msg interface{}, peer net.Session) {
	utils.SafeCall(func() {
		if handler, ok := s.handlers[cmd]; ok {
			handler(msg, peer)
		} else {
		}
	})
}
