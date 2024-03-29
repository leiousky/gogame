package tcp

import (
	"errors"
	"games/comm/utils"
	cb "games/core/callback"
	"games/core/conn"
	"games/core/conn/def"
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
	reason          conn.Reason
	onConnection    cb.OnConnection
	onMessage       cb.OnMessage
	onWriteComplete cb.OnWriteComplete
	closeCallback   cb.CloseCallback
	errorCallback   cb.ErrorCallback
}

func NewTCPConnection(id int64, name string, c interface{}, connType conn.Type, channel transmit.IChannel) conn.Session {
	peer := &TCPConnection{
		id:       id,
		name:     name,
		conn:     c,
		state:    conn.KDisconnected,
		reason:   conn.KNoError,
		connType: connType,
		context:  map[int]interface{}{},
		msq:      msq.NewBlockVecMsq(),
		channel:  channel}
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

func (s *TCPConnection) setReason(reason conn.Reason) {
	s.reason = reason
}

func (s *TCPConnection) Connected() bool {
	return s.state == conn.KConnected
}

func (s *TCPConnection) LocalAddr() string {
	if s.conn == nil {
		panic(errors.New("s.conn == nil"))
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
		panic(errors.New("s.conn == nil"))
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
		panic(errors.New("s.conn == nil"))
	}
	if _, ok := s.conn.(net.Conn); ok {
		return false
	} else if _, ok := s.conn.(*websocket.Conn); ok {
		return true
	}
	panic(errors.New("s.conn error"))
}

func (s *TCPConnection) Conn() interface{} {
	return s.conn
}

func (s *TCPConnection) Type() conn.Type {
	return s.connType
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
	s.setState(conn.KConnected)
	go s.readLoop()
	go s.writeLoop()
	if s.onConnection != nil {
		s.onConnection(s, s.reason)
	}
}

func (s *TCPConnection) ConnectDestroyed() {
	if s.id == 0 {
		panic(errors.New("connID == 0"))
	}
	s.setState(conn.KDisconnected)
	if s.onConnection != nil {
		s.onConnection(s, s.reason)
	}
}

/// 读协程
func (s *TCPConnection) readLoop() {
	utils.CheckPanic()
	for {
		msg, err, reason := s.channel.OnRecv(s.conn)
		if err != nil {
			//log.Println("readLoop: ", err)
			// if !IsEOFOrReadError(err) {
			// 	if s.errorCallback != nil {
			// 		s.errorCallback(err)
			// 	}
			// }
			if reason == def.KExcept {
				//本端异常关闭
				s.setReason(conn.KSelfExcept)
				//通知写退出
				s.msq.Push(nil)
				break
			} else if 0 == atomic.LoadInt64(&s.closing) {
				//对端关闭连接
				s.setReason(conn.KPeerClosed)
				//通知写退出
				s.msq.Push(nil)
				break
			} else if reason == def.KClosed {
				//本端正常关闭
				s.setReason(conn.KSelfClosed)
				break
			} else if s.errorCallback != nil {
				s.errorCallback(err)
			}
		} else if msg == nil {
			panic(errors.New("readLoop: msg == nil"))
		} else if s.onMessage != nil {
			s.onMessage(s, msg, utils.TimeNow())
		}
	}
	//等待写退出
	s.Wg.Wait()
	// if s.closeCallback != nil {
	// 	s.closeCallback(s)
	// }
	// s.conn = nil
}

/// 写协程
func (s *TCPConnection) writeLoop() {
	utils.CheckPanic()
	for {
		msgs, exit := s.msq.Pick()
		for _, msg := range msgs {
			err, _ := s.channel.OnSend(s.conn, msg)
			if err != nil {
				log.Println("writeLoop: ", err)
				if !transmit.IsEOFOrWriteError(err) {
					if s.errorCallback != nil {
						s.errorCallback(err)
					}
					//本端异常关闭
					//s.setReason(conn.KSelfExcept)
				}
				//break
			} else if s.onWriteComplete != nil {
				s.onWriteComplete(s)
			}
		}
		if exit {
			break
		}
	}
	//TCPServer.removeConnection ->
	//Sessions.Remove() ->
	//TCPConnection.ConnectDestroyed ->
	//TCPConnection.close
	if s.closeCallback != nil {
		s.closeCallback(s)
	}
	//唤醒阻塞读
	s.close()
	s.Wg.Done()
}

/// 写数据
func (s *TCPConnection) Write(msg interface{}) {
	s.msq.Push(msg)
}

/// 关闭连接
func (s *TCPConnection) Close() {
	if 0 == atomic.SwapInt64(&s.closing, 1) && s.conn != nil {
		//本端正常关闭
		s.setReason(conn.KSelfClosed)
		//通知写退出
		s.msq.Push(nil)
	}
}

/// 关闭执行流程
/// TCPConnection.closeCallback ->
/// TCPServer.removeConnection ->
/// Sessions.Remove() ->
/// TCPConnection.ConnectDestroyed ->
/// TCPConnection.close
func (s *TCPConnection) close() {
	if s.conn == nil {
		return
	}
	log.Printf("TCPConnection.close => %v", s.name)
	if c, ok := s.conn.(net.Conn); ok {
		c.Close()
	} else if c, ok := s.conn.(*websocket.Conn); ok {
		c.Close()
	}
}
