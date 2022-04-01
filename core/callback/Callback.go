package callback

import (
	"games/comm/utils"
	"games/core/conn"
	"games/core/conn/transmit"
)

/// <summary>
/// Functor 回调函数
/// <summary>
type Functor struct {
	// f    func(args ...interface{})
	// args []interface{}
	f    func(args interface{})
	args interface{}
}

// func NewFunctor(f func(args ...interface{}), args ...interface{}) Functor {
// 	s := Functor{f: f, args: args}
// 	// s := Functor{f: f}
// 	// s.args = append(s.args, args...)
// 	return s
// }

func NewFunctor(f func(args interface{}), args interface{}) Functor {
	s := Functor{f: f, args: args}
	return s
}

func (s Functor) Call() {
	// s.f(s.args...)
	s.f(s.args)
}

/// 流协议解析
type OnProtocol func(proto string) transmit.IChannel

/// 接受新连接条件
type OnCondition func(conn interface{}) bool

type OnNewConnection func(conn interface{}, channel transmit.IChannel)

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
