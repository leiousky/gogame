package core

// import "games/core/net"

// type CloseCallback func(peer net.Session)

// type OnConnected func(peer net.Session)

// type OnClosed func(peer net.Session)

// type OnMessage func(msg interface{}, peer net.Session)

// type OnWritten func(msg interface{}, peer net.Session)

// type OnError func(peer net.Session, err error)

// type ReadCallback func(cmd uint32, msg interface{}, peer net.Session)

// type CustomCallback func(cmd uint32, msg interface{}, peer net.Session)

// type CmdCallback func(msg interface{}, peer net.Session)

// type CmdCallbacks map[uint32]CmdCallback

type TimerCallback func(timerID uint32, dt int32, args interface{}) bool
