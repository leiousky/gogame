package tcpserver

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

	"github.com/gorilla/websocket"
)

/// <summary>
/// Processor TCP服务端
/// <summary>
type Processor struct {
	name            string
	peer            conn.Session
	acceptor        tcp.Acceptor
	onCondition     cb.OnCondition
	onConnection    cb.OnConnection
	onMessage       cb.OnMessage
	onWriteComplete cb.OnWriteComplete
}

func NewTCPServer(name string) TCPServer {
	s := &Processor{
		name:     name,
		acceptor: tcp.NewAcceptor()}
	s.acceptor.SetProtocolCallback(s.onProtocol)
	s.acceptor.SetNewConnectionCallback(s.newConnection)
	s.SetConditionCallback(s.OnCondition)
	s.SetConnectionCallback(s.OnConnection)
	s.SetMessageCallback(s.OnMessage)
	s.SetWriteCompleteCallback(s.OnWriteComplete)
	return s
}

func (s *Processor) SetProtocolCallback(cb cb.OnProtocol) {
	if s.acceptor == nil {
		panic(errors.New("TCPServer SetProtocolCallback s.acceptor == nil"))
	}
	s.acceptor.SetProtocolCallback(cb)
}

func (s *Processor) SetConditionCallback(cb cb.OnCondition) {
	if s.acceptor == nil {
		panic(errors.New("TCPServer SetConditionCallback s.acceptor == nil"))
	}
	s.acceptor.SetConditionCallback(cb)
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

func (s *Processor) ListenTCP(address string) {
	s.acceptor.ListenTCP(address)
}

func (s *Processor) Stop() {
	s.acceptor.Close()
}

func (s *Processor) OnCondition(c interface{}) bool {
	return true
}

func (s *Processor) newConnection(c interface{}, channel transmit.IChannel) {
	connID := conn.NewConnID()
	if p, ok := c.(net.Conn); ok {
		localAddr := p.LocalAddr().String()
		peerAddr := p.RemoteAddr().String()
		s.peer = tcp.NewTCPConnection(
			connID,
			s.name+fmt.Sprintf("#%v<-%v#%v", localAddr, peerAddr, connID),
			c,
			conn.KServer,
			channel)
	} else if p, ok := c.(*websocket.Conn); ok {
		localAddr := p.LocalAddr().String()
		peerAddr := p.RemoteAddr().String()
		s.peer = tcp.NewTCPConnection(
			connID,
			s.name+fmt.Sprintf("#%v<-%v#%v", localAddr, peerAddr, connID),
			c,
			conn.KServer,
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

func (s *Processor) OnConnection(peer conn.Session, reason conn.Reason) {
	if peer.Connected() {
		log.Print("--- *** TCPServer - TCPServer:: OnConnected \n")
	} else {
		log.Print("--- *** TCPServer - TCPServer:: OnClosed \n")
	}
}

func (s *Processor) OnMessage(peer conn.Session, msg interface{}, recvTime utils.Timestamp) {
	log.Print("--- *** TCPServer - TCPServer:: OnMessage \n")
}

func (s *Processor) OnWriteComplete(peer conn.Session) {
	log.Print("--- *** TCPServer - TCPServer:: OnWriteComplete \n")
}

func (s *Processor) removeConnection(peer conn.Session) {
	sessions.Remove(peer)
	peer.(*tcp.TCPConnection).ConnectDestroyed()
}

func (s *Processor) onConnectionError(err error) {
	log.Print("--- *** TCPServer - TCPServer:: onConnectionError \n")
}
