package net

import (
	"games/comm/utils"
	"games/core"
	"games/core/msq"
	"games/core/transmit"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

/// <summary>
/// TCPConnection TCP/WS连接会话
/// <summary>
type TCPConnection struct {
	SesID         int64
	conn          interface{}
	context       map[int]interface{}
	closing       int64
	Wg            sync.WaitGroup
	WMsq          msq.MsgQueue      //写队列
	RMsq          msq.MsgQueue      //读队列
	Channel       transmit.IChannel //消息传输协议
	slot          core.ISlot        //处理单元cell
	onConnected   OnConnected
	onClosed      OnClosed
	onMessage     OnMessage
	onWritten     OnWritten
	onError       OnError
	closeCallback CloseCallback
	//reason        int32
	//sta           int32
}

func newTCPConnection(conn interface{}, channel transmit.IChannel) Session {
	peer := &TCPConnection{
		SesID:   createSessionID(),
		conn:    conn,
		context: map[int]interface{}{},
		WMsq:    msq.NewBlockVecMsq(),
		Channel: channel}
	return peer
}

func (s *TCPConnection) ID() int64 {
	return s.SesID
}

func (s *TCPConnection) IsWebsocket() bool {
	if s.conn == nil {
		return false
	}
	if _, ok := s.conn.(net.Conn); ok {
		return false
	} else if _, ok := s.conn.(*websocket.Conn); ok {
		return true
	}
	return false
}

func (s *TCPConnection) Conn() interface{} {
	return s.conn
}

func (s *TCPConnection) SetContext(key int, val interface{}) {
	if val != nil {
		s.context[key] = val
	} else if _, ok := s.context[key]; ok {
		delete(s.context, key)
	}
}

func (s *TCPConnection) GetContext(key int) interface{} {
	if val, ok := s.context[key]; ok {
		return val
	}
	return nil
}

func (s *TCPConnection) SetOnConnected(cb OnConnected) {
	s.onConnected = cb
}

func (s *TCPConnection) SetOnClosed(cb OnClosed) {
	s.onClosed = cb
}

func (s *TCPConnection) SetOnMessage(cb OnMessage) {
	s.onMessage = cb
}

func (s *TCPConnection) SetOnError(cb OnError) {
	s.onError = cb
}

func (s *TCPConnection) SetOnWritten(cb OnWritten) {
	s.onWritten = cb
}

func (s *TCPConnection) SetCloseCallback(cb CloseCallback) {
	s.closeCallback = cb
}

/// 读协程
func (s *TCPConnection) readLoop() {
	utils.CheckPanic()
	for {
		msg, err := s.Channel.OnRecvMessage(s.conn)
		if err != nil {
			//log.Println("readLoop: ", err)
			// if !IsEOFOrReadError(err) {
			// 	if s.onError != nil {
			// 		s.onError(s, err)
			// 	}
			// }
			break
		}
		if msg == nil {
			log.Fatalln("readLoop: msg == nil")
		}
		if s.onMessage != nil {
			s.onMessage(msg, s)
		}
	}
	//对端关闭连接
	if 0 == atomic.LoadInt64(&s.closing) {
		//通知写退出
		s.WMsq.Push(nil)
	}
	//等待写退出
	s.Wg.Wait()
	if s.closeCallback != nil {
		s.closeCallback(s)
	}
	s.conn = nil
}

/// 写协程
func (s *TCPConnection) writeLoop() {
	utils.CheckPanic()
	for {
		msgs, exit := s.WMsq.Pick()
		for _, msg := range msgs {
			err := s.Channel.OnSendMessage(s.conn, msg)
			if err != nil {
				log.Println("writeLoop: ", err)
				if !transmit.IsEOFOrWriteError(err) {
					if s.onError != nil {
						s.onError(s, err)
					}
				}
				//break
			}
			if s.onWritten != nil {
				s.onWritten(msg, s)
			}
		}
		if exit {
			break
		}
	}
	//唤醒阻塞读
	s.close()
	s.Wg.Done()
}

/// 写
func (s *TCPConnection) Write(msg interface{}) {
	s.WMsq.Push(msg)
}

/// 关闭
func (s *TCPConnection) Close() {
	//本端关闭连接
	if 0 == atomic.SwapInt64(&s.closing, 1) && s.conn != nil {
		//通知写退出
		s.WMsq.Push(nil)
	}
}

/// 关闭连接对象
func (s *TCPConnection) close() {
	if s.conn == nil {
		return
	}
	if _, ok := s.conn.(net.Conn); ok {
		s.conn.(net.Conn).Close()
	} else if _, ok := s.conn.(*websocket.Conn); ok {
		s.conn.(*websocket.Conn).Close()
	}
}
