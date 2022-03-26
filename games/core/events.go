package core

import "games/core/net"

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
	peer    net.Session
	msg     interface{}
	handler net.ReadCallback
}

func createReadEvent(cmd uint32, msg interface{}, peer net.Session) *readEvent {
	ev := &readEvent{cmd: cmd, msg: msg, peer: peer}
	return ev
}

func createReadEventWith(handler net.ReadCallback, cmd uint32, msg interface{}, peer net.Session) *readEvent {
	ev := &readEvent{handler: handler, cmd: cmd, msg: msg, peer: peer}
	return ev
}

/// <summary>
/// customEvent 自定义事件
/// <summary>
type customEvent struct {
	cmd     uint32
	peer    net.Session
	msg     interface{}
	handler net.CustomCallback
}

func createCustomEvent(cmd uint32, msg interface{}, peer net.Session) *customEvent {
	ev := &customEvent{cmd: cmd, msg: msg, peer: peer}
	return ev
}

func createCustomEventWith(handler net.CustomCallback, cmd uint32, msg interface{}, peer net.Session) *customEvent {
	ev := &customEvent{handler: handler, cmd: cmd, msg: msg, peer: peer}
	return ev
}
