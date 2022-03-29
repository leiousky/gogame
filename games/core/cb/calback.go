package cb

import "games/core/conn"

type CloseCallback func(peer conn.Session)

type OnConnected func(peer conn.Session)

type OnClosed func(peer conn.Session)

type OnMessage func(msg interface{}, peer conn.Session)

type OnWritten func(msg interface{}, peer conn.Session)

type OnError func(peer conn.Session, err error)

type ReadCallback func(cmd uint32, msg interface{}, peer conn.Session)

type CustomCallback func(cmd uint32, msg interface{}, peer conn.Session)

type CmdCallback func(msg interface{}, peer conn.Session)

type CmdCallbacks map[uint32]CmdCallback

type TimerCallback func(timerID uint32, dt int32, args interface{}) bool
