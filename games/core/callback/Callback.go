package callback

import (
	"games/comm/utils"
	"games/core/conn"
)

type OnConnection func(peer conn.Session)

type OnMessage func(peer conn.Session, msg interface{}, recvTime utils.Timestamp)

type OnWriteComplete func(peer conn.Session)

type CloseCallback func(peer conn.Session)

type ErrorCallback func(err error)

type ReadCallback func(cmd uint32, msg interface{}, peer conn.Session)

type CustomCallback func(cmd uint32, msg interface{}, peer conn.Session)

type CmdCallback func(msg interface{}, peer conn.Session)

type CmdCallbacks map[uint32]CmdCallback

type TimerCallback func(timerID uint32, dt int32, args interface{}) bool
