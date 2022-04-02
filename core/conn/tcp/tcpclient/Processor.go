package tcpclient

import (
	"errors"
	"fmt"
	"games/comm/utils"
	cb "games/core/callback"
	"games/core/conn"
	"games/core/conn/tcp"
	"games/core/conn/transmit"
	"games/core/conn/transmit/tcp_channel"
	"games/core/conn/transmit/ws_channel"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

/// <summary>
/// Processor TCP客户端
/// <summary>
type Processor struct {
	name            string
	peer            conn.Session
	c               tcp.Connector
	onConnection    cb.OnConnection
	onMessage       cb.OnMessage
	onWriteComplete cb.OnWriteComplete
}

func NewTCPClient(name string) TCPClient {
	s := &Processor{
		name: name,
		c:    tcp.NewConnector(name)}
	s.c.SetProtocolCallback(s.onProtocol)
	s.c.SetNewConnectionCallback(s.newConnection)
	s.SetConnectionCallback(s.OnConnection)
	s.SetMessageCallback(s.OnMessage)
	s.SetWriteCompleteCallback(s.OnWriteComplete)
	return s
}

func (s *Processor) EnableRetry(bv bool) {
	s.c.EnableRetry(bv)
}

func (s *Processor) Session() conn.Session {
	return s.peer
}

func (s *Processor) Write(msg interface{}) {
	if s.peer != nil {
		s.peer.Write(msg)
	}
}

func (s *Processor) SetProtocolCallback(cb cb.OnProtocol) {
	if s.c == nil {
		panic(errors.New("TCPClient SetProtocolCallback s.c == nil"))
	}
	s.c.SetProtocolCallback(cb)
}

func (s *Processor) SetConnectionCallback(cb cb.OnConnection) {
	s.onConnection = cb
}

func (s *Processor) SetMessageCallback(cb cb.OnMessage) {
	s.onMessage = cb
}

func (s *Processor) SetWriteCompleteCallback(cb cb.OnWriteComplete) {
	s.onWriteComplete = cb
}

func (s *Processor) newConnection(c interface{}, channel transmit.IChannel) {
	connID := conn.NewConnID()
	if p, ok := c.(net.Conn); ok {
		localAddr := p.LocalAddr().String()
		peerAddr := p.RemoteAddr().String()
		s.peer = tcp.NewTCPConnection(
			connID,
			s.name+fmt.Sprintf("#%v->%v#%v", localAddr, peerAddr, connID),
			c,
			conn.KClient,
			channel)
	} else if p, ok := c.(*websocket.Conn); ok {
		localAddr := p.LocalAddr().String()
		peerAddr := p.RemoteAddr().String()
		s.peer = tcp.NewTCPConnection(
			connID,
			s.name+fmt.Sprintf("#%v->%v#%v", localAddr, peerAddr, connID),
			c,
			conn.KClient,
			channel)
	} else {
		panic(errors.New("newConnection conn error"))
	}
	s.peer.(*tcp.TCPConnection).SetConnectionCallback(s.onConnection)
	s.peer.(*tcp.TCPConnection).SetMessageCallback(s.onMessage)
	s.peer.(*tcp.TCPConnection).SetWriteCompleteCallback(s.onWriteComplete)
	s.peer.(*tcp.TCPConnection).SetCloseCallback(s.removeConnection)
	s.peer.(*tcp.TCPConnection).SetErrorCallback(s.onConnectionError)
	s.peer.(*tcp.TCPConnection).ConnectEstablished()
	if !sessions.Add(s.peer) {
		s.peer.Close()
	}
}

func (s *Processor) onProtocol(proto string) transmit.IChannel {
	switch proto {
	case "tcp":
		return tcp_channel.NewChannel()
	case "ws":
		return ws_channel.NewChannel()
	}
	panic(errors.New("no proto setup"))
}

func (s *Processor) ConnectTCP(address string) {
	s.c.ConnectTCP(address)
}

func (s *Processor) OnConnection(peer conn.Session, reason conn.Reason) {
	if peer.Connected() {
		log.Print("--- *** TCPClient - TCPClient:: OnConnected \n")
	} else {
		log.Print("--- *** TCPClient - TCPClient:: OnClosed \n")
	}
}

func (s *Processor) OnMessage(peer conn.Session, msg interface{}, recvTime utils.Timestamp) {
	log.Print("--- *** TCPClient - TCPClient:: OnMessage \n")
}

func (s *Processor) OnWriteComplete(peer conn.Session) {
	log.Print("--- *** TCPClient - TCPClient:: OnWriteComplete \n")
}

func (s *Processor) removeConnection(peer conn.Session) {
	sessions.Remove(peer)
	peer.(*tcp.TCPConnection).ConnectDestroyed()
	if s.c.Retry() {
		s.c.Reconnect(time.Millisecond * 500)
	}
}

func (s *Processor) onConnectionError(err error) {
	log.Printf("--- *** TCPClient - TCPClient:: onConnectionError \n")
}

func (s *Processor) Reconnect(d time.Duration) {
	s.c.Reconnect(d)
}

func (s *Processor) Disconnect() {
	if s.peer != nil {
		s.peer.Close()
	}
}
