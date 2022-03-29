package net

import (
	"games/comm/utils"
	"games/core/msq"
	"games/core/net/transmit"

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
	id            int64
	name          string
	conn          interface{}
	context       map[int]interface{}
	connType      SesType
	msq           msq.MsgQueue
	channel       transmit.IChannel
	Wg            sync.WaitGroup
	closing       int64
	onConnected   OnConnected
	onMessage     OnMessage
	onClosed      OnClosed
	onWritten     OnWritten
	onError       OnError
	closeCallback CloseCallback
}

func NewTCPConnection(id int64, name string, conn interface{}, connType SesType, channel transmit.IChannel) Session {
	peer := &TCPConnection{
		id:       id,
		name:     name,
		conn:     conn,
		connType: connType,
		context:  map[int]interface{}{},
		msq:      msq.NewBlockVecMsq(),
		channel:  channel}
	return peer
}

func (s *TCPConnection) ID() int64 {
	return s.connID
}

func (s *TCPConnection) Name() string {
	return s.name
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

func (s *TCPConnection) Type() SesType {
	return s.connType
}

func (s *TCPConnection) SetChannel(channel transmit.IChannel) {
	s.channel = channel
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
		msg, err := s.channel.OnRecvMessage(s.conn)
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
		s.msq.Push(nil)
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
		msgs, exit := s.msq.Pick()
		for _, msg := range msgs {
			err := s.channel.OnSendMessage(s.conn, msg)
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

func (s *TCPConnection) Write(msg interface{}) {
	s.msq.Push(msg)
}

func (s *TCPConnection) Close() {
	//本端关闭连接
	if 0 == atomic.SwapInt64(&s.closing, 1) && s.conn != nil {
		//通知写退出
		s.msq.Push(nil)
	}
}

func (s *TCPConnection) close() {
	if s.conn == nil {
		return
	}
	if c, ok := s.conn.(net.Conn); ok {
		c.Close()
	} else if c, ok := s.conn.(*websocket.Conn); ok {
		c.Close()
	}
}
