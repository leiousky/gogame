package tcpserver

import (
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
/// ITCPServer TCP服务端
/// <summary>
type ITCPServer interface {
	ListenTCP(address string)
	OnConnection(peer conn.Session)
	OnMessage(peer conn.Session, msg interface{}, recvTime utils.Timestamp)
	OnWriteComplete(peer conn.Session)
	Stop()
}

/// <summary>
/// sessions 服务端容器
/// <summary>
var sessions = conn.NewSessions()

/// <summary>
/// TCPServer TCP服务端
/// <summary>
type TCPServer struct {
	name            string
	peer            conn.Session
	acceptor        tcp.Acceptor
	onCondition     cb.OnCondition
	onConnection    cb.OnConnection
	onMessage       cb.OnMessage
	onWriteComplete cb.OnWriteComplete
}

func NewTCPServer() ITCPServer {
	s := &TCPServer{
		acceptor: tcp.NewAcceptor()}
	s.acceptor.SetProtocolCallback(s.onProtocol)
	s.acceptor.SetNewConnectionCallback(s.newConnection)
	s.SetConditionCallback(s.OnCondition)
	s.SetConnectionCallback(s.OnConnection)
	s.SetMessageCallback(s.OnMessage)
	s.SetWriteCompleteCallback(s.OnWriteComplete)
	return s
}

func (s *TCPServer) SetProtocolCallback(cb cb.OnProtocol) {
	if s.acceptor == nil {
		panic(fmt.Sprintf("TCPServer SetProtocolCallback s.acceptor == nil"))
	}
	s.acceptor.SetProtocolCallback(cb)
}

func (s *TCPServer) SetConditionCallback(cb cb.OnCondition) {
	if s.acceptor == nil {
		panic(fmt.Sprintf("TCPServer SetConditionCallback s.acceptor == nil"))
	}
	s.acceptor.SetConditionCallback(cb)
}

func (s *TCPServer) SetConnectionCallback(cb cb.OnConnection) {
	s.onConnection = cb
}

func (s *TCPServer) SetMessageCallback(cb cb.OnMessage) {
	s.onMessage = cb
}

func (s *TCPServer) SetWriteCompleteCallback(cb cb.OnWriteComplete) {
	s.onWriteComplete = cb
}

func (s *TCPServer) ListenTCP(address string) {
	s.acceptor.ListenTCP(address)
}

func (s *TCPServer) Stop() {
	s.acceptor.Close()
}

func (s *TCPServer) OnCondition(c interface{}) bool {
	return true
}

func (s *TCPServer) newConnection(c interface{}, channel transmit.IChannel) {
	connID := conn.NewConnID()
	if p, ok := c.(net.Conn); ok {
		localAddr := p.LocalAddr().String()
		peerAddr := p.RemoteAddr().String()
		s.peer = tcp.NewTCPConnection(
			connID,
			s.name+fmt.Sprintf("-%v|%v#%v", localAddr, peerAddr, connID),
			c,
			conn.KServer,
			channel)
	} else if p, ok := c.(*websocket.Conn); ok {
		localAddr := p.LocalAddr().String()
		peerAddr := p.RemoteAddr().String()
		s.peer = tcp.NewTCPConnection(
			connID,
			s.name+fmt.Sprintf("-%v|%v#%v", localAddr, peerAddr, connID),
			c,
			conn.KServer,
			channel)
	} else {
		panic(fmt.Sprintf("newConnection conn error"))
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

func (s *TCPServer) onProtocol(proto string) transmit.IChannel {
	switch proto {
	case "tcp":
		return tcp_channel.NewChannel()
	case "ws":
		return ws_channel.NewChannel()
	}
	panic(fmt.Sprintf("no proto setup"))
}

func (s *TCPServer) OnConnection(peer conn.Session) {
	if peer.Connected() {
		log.Print("--- *** TCPServer - TCPServer:: OnConnected \n")
	} else {
		log.Print("--- *** TCPServer - TCPServer:: OnClosed \n")
	}
}

func (s *TCPServer) OnMessage(peer conn.Session, msg interface{}, recvTime utils.Timestamp) {
	log.Print("--- *** TCPServer - TCPServer:: OnMessage \n")
}

func (s *TCPServer) OnWriteComplete(peer conn.Session) {
	log.Print("--- *** TCPServer - TCPServer:: OnWriteComplete \n")
}

func (s *TCPServer) removeConnection(peer conn.Session) {
	sessions.Remove(peer)
	peer.(*tcp.TCPConnection).ConnectDestroyed()
}

func (s *TCPServer) onConnectionError(err error) {
	log.Print("--- *** TCPServer - TCPServer:: onConnectionError \n")
}
