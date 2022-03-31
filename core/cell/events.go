package cell

import (
	cb "games/core/callback"
	"games/core/conn"
)

const (
	EVTConnected int8 = iota + 10 // 连接
	EVTClosed                     // 关闭
	EVTRead                       // 读
	EVTSend                       // 写
	EVTCustom                     // 自定义
	EVTLogger                     // 日志
	EVTClose                      // 延迟关闭
)

/// <summary>
/// Event 事件数据
/// <summary>
type Event struct {
	ev  int8
	obj interface{}
	ext interface{}
}

func createEvent(ev int8, obj interface{}, ext interface{}) *Event {
	return &Event{ev, obj, ext}
}

/// <summary>
/// readEvent 读事件
/// <summary>
type readEvent struct {
	cmd     uint32
	peer    conn.Session
	msg     interface{}
	handler cb.ReadCallback
}

func createReadEvent(cmd uint32, msg interface{}, peer conn.Session) *readEvent {
	ev := &readEvent{cmd: cmd, msg: msg, peer: peer}
	return ev
}

func createReadEventWith(handler cb.ReadCallback, cmd uint32, msg interface{}, peer conn.Session) *readEvent {
	ev := &readEvent{handler: handler, cmd: cmd, msg: msg, peer: peer}
	return ev
}

/// <summary>
/// customEvent 自定义事件
/// <summary>
type customEvent struct {
	cmd     uint32
	peer    conn.Session
	msg     interface{}
	handler cb.CustomCallback
}

func createCustomEvent(cmd uint32, msg interface{}, peer conn.Session) *customEvent {
	ev := &customEvent{cmd: cmd, msg: msg, peer: peer}
	return ev
}

func createCustomEventWith(handler cb.CustomCallback, cmd uint32, msg interface{}, peer conn.Session) *customEvent {
	ev := &customEvent{handler: handler, cmd: cmd, msg: msg, peer: peer}
	return ev
}
