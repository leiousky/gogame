package tcp

import (
	"fmt"
	"games/comm/utils"
	cb "games/core/callback"
	"games/core/conn"
	"games/core/conn/transmit"
	"games/core/msq"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

/// <summary>
/// TCPConnection TCP连接会话
/// <summary>
type TCPConnection struct {
	id              int64
	name            string
	conn            interface{}
	context         map[int]interface{}
	connType        conn.Type
	msq             msq.MsgQueue
	channel         transmit.IChannel
	Wg              sync.WaitGroup
	closing         int64
	state           conn.State
	onConnection    cb.OnConnection
	onMessage       cb.OnMessage
	onWriteComplete cb.OnWriteComplete
	closeCallback   cb.CloseCallback
	errorCallback   cb.ErrorCallback
}

func NewTCPConnection(id int64, name string, c interface{}, connType conn.Type) conn.Session {
	peer := &TCPConnection{
		id:       id,
		name:     name,
		conn:     c,
		state:    conn.KDisconnected,
		connType: connType,
		context:  map[int]interface{}{},
		msq:      msq.NewBlockVecMsq()}
	return peer
}

func (s *TCPConnection) ID() int64 {
	return s.id
}

func (s *TCPConnection) Name() string {
	return s.name
}

func (s *TCPConnection) setState(state conn.State) {
	s.state = state
}

func (s *TCPConnection) Connected() bool {
	return s.state == conn.KConnected
}

func (s *TCPConnection) LocalAddr() string {
	if s.conn == nil {
		panic(fmt.Sprintf("s.conn == nil"))
	}
	if c, ok := s.conn.(net.Conn); ok {
		return c.LocalAddr().String()
	} else if c, ok := s.conn.(*websocket.Conn); ok {
		return c.LocalAddr().String()
	}
	return ""
}

func (s *TCPConnection) RemoteAddr() string {
	if s.conn == nil {
		panic(fmt.Sprintf("s.conn == nil"))
	}
	if c, ok := s.conn.(net.Conn); ok {
		return c.RemoteAddr().String()
	} else if c, ok := s.conn.(*websocket.Conn); ok {
		return c.RemoteAddr().String()
	}
	return ""
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

func (s *TCPConnection) Type() conn.Type {
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

func (s *TCPConnection) SetConnectionCallback(cb cb.OnConnection) {
	s.onConnection = cb
}

func (s *TCPConnection) SetMessageCallback(cb cb.OnMessage) {
	s.onMessage = cb
}

func (s *TCPConnection) SetWriteCompleteCallback(cb cb.OnWriteComplete) {
	s.onWriteComplete = cb
}

func (s *TCPConnection) SetCloseCallback(cb cb.CloseCallback) {
	s.closeCallback = cb
}

func (s *TCPConnection) SetErrorCallback(cb cb.ErrorCallback) {
	s.errorCallback = cb
}

func (s *TCPConnection) ConnectEstablished() {
	s.Wg.Add(1)
	go s.readLoop()
	go s.writeLoop()
	s.setState(conn.KConnected)
	if s.onConnection != nil {
		s.onConnection(s)
	}
}

func (s *TCPConnection) ConnectDestroyed() {
	if s.id == 0 {
		panic("connID == 0")
	}
	s.setState(conn.KDisconnected)
	if s.onConnection != nil {
		s.onConnection(s)
	}
}

/// 读协程
func (s *TCPConnection) readLoop() {
	utils.CheckPanic()
	for {
		msg, err := s.channel.OnRecv(s.conn)
		if err != nil {
			//log.Println("readLoop: ", err)
			// if !IsEOFOrReadError(err) {
			// 	if s.errorCallback != nil {
			// 		s.errorCallback(err)
			// 	}
			// }
			break
		}
		if msg == nil {
			log.Fatalln("readLoop: msg == nil")
		}
		if s.onMessage != nil {
			s.onMessage(s, msg, utils.TimeNow())
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
			err := s.channel.OnSend(s.conn, msg)
			if err != nil {
				log.Println("writeLoop: ", err)
				if !transmit.IsEOFOrWriteError(err) {
					if s.errorCallback != nil {
						s.errorCallback(err)
					}
				}
				//break
			}
			if s.onWriteComplete != nil {
				s.onWriteComplete(s)
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
