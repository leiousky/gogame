package core

type Session interface{}

//
type CloseCallback func(peer Session)

//
type OnConnected func(peer Session)

//
type OnClosed func(peer Session)

//
type OnMessage func(msg interface{}, peer Session)

//
type OnWritten func(msg interface{}, peer Session)

//
type OnError func(peer Session, err error)

//
type ReadCallback func(cmd uint32, msg interface{}, peer Session)

//
type CustomCallback func(cmd uint32, msg interface{}, peer Session)

//
type CmdCallback func(msg interface{}, peer Session)

//
type CmdCallbacks map[uint32]CmdCallback

//
type TimerCallback func(timerID uint32, dt int32, args interface{}) bool
