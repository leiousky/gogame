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
)

/// <summary>
/// ITCPServer TCP服务端
/// <summary>
type ITCPServer interface {
	Start()
	Stop()
	OnConnection(peer conn.Session)
	OnMessage(peer conn.Session, msg interface{}, recvTime utils.Timestamp)
	OnWriteComplete(peer conn.Session)
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
	c               tcp.Acceptor
	onCondition     cb.OnCondition
	onConnection    cb.OnConnection
	onMessage       cb.OnMessage
	onWriteComplete cb.OnWriteComplete
}

func NewTCPServer() ITCPServer {
	s := &TCPServer{
		c: tcp.NewAcceptor()}
	s.c.SetProtocolCallback(s.onProtocol)
	s.c.SetNewConnectionCallback(s.newConnection)
	s.SetConditionCallback(s.OnCondition)
	s.SetConnectionCallback(s.OnConnection)
	s.SetMessageCallback(s.OnMessage)
	s.SetWriteCompleteCallback(s.OnWriteComplete)
	return s
}
func (s *TCPServer) SetProtocolCallback(cb cb.OnProtocol) {
	if s.c == nil {
		panic(fmt.Sprintf("TCPServer SetProtocolCallback s.c == nil"))
	}
	s.c.SetProtocolCallback(cb)
}

func (s *TCPServer) SetConditionCallback(cb cb.OnCondition) {
	if s.c == nil {
		panic(fmt.Sprintf("TCPServer SetConditionCallback s.c == nil"))
	}
	s.c.SetConditionCallback(cb)
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

func (s *TCPServer) Start() {

}

func (s *TCPServer) Stop() {

}

func (s *TCPServer) OnCondition(c interface{}) bool {
	return true
}

func (s *TCPServer) newConnection(c interface{}, channel transmit.IChannel) {
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
