package tcpclient

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
	"time"

	"github.com/gorilla/websocket"
)

/// <summary>
/// ITCPClient TCP客户端
/// <summary>
type ITCPClient interface {
	Session() conn.Session
	Write(msg interface{})
	ConnectTCP(address string)
	OnConnection(peer conn.Session)
	OnMessage(peer conn.Session, msg interface{}, recvTime utils.Timestamp)
	OnWriteComplete(peer conn.Session)
	Reconnect(d time.Duration)
	Disconnect()
}

/// <summary>
/// sessions 客户端容器
/// <summary>
var sessions = conn.NewSessions()

/// <summary>
/// TCPClient TCP客户端
/// <summary>
type TCPClient struct {
	name            string
	peer            conn.Session
	c               tcp.Connector
	onConnection    cb.OnConnection
	onMessage       cb.OnMessage
	onWriteComplete cb.OnWriteComplete
}

func NewTCPClient() ITCPClient {
	s := &TCPClient{
		c: tcp.NewConnector()}
	s.c.SetProtocolCallback(s.onProtocol)
	s.c.SetNewConnectionCallback(s.newConnection)
	s.SetConnectionCallback(s.OnConnection)
	s.SetMessageCallback(s.OnMessage)
	s.SetWriteCompleteCallback(s.OnWriteComplete)
	return s
}

func (s *TCPClient) Session() conn.Session {
	return s.peer
}

func (s *TCPClient) Write(msg interface{}) {
	if s.peer != nil {
		s.peer.Write(msg)
	}
}

func (s *TCPClient) SetProtocolCallback(cb cb.OnProtocol) {
	if s.c == nil {
		panic(fmt.Sprintf("TCPClient SetProtocolCallback s.c == nil"))
	}
	s.c.SetProtocolCallback(cb)
}

func (s *TCPClient) SetConnectionCallback(cb cb.OnConnection) {
	s.onConnection = cb
}

func (s *TCPClient) SetMessageCallback(cb cb.OnMessage) {
	s.onMessage = cb
}

func (s *TCPClient) SetWriteCompleteCallback(cb cb.OnWriteComplete) {
	s.onWriteComplete = cb
}

func (s *TCPClient) newConnection(c interface{}, channel transmit.IChannel) {
	connID := conn.NewConnID()
	if p, ok := c.(net.Conn); ok {
		//localAddr := p.LocalAddr().String()
		peerAddr := p.RemoteAddr().String()
		s.peer = tcp.NewTCPConnection(
			connID,
			s.name+fmt.Sprintf("-%v#%v", peerAddr, connID),
			c,
			conn.KClient,
			channel)
	} else if p, ok := c.(*websocket.Conn); ok {
		//localAddr := p.LocalAddr().String()
		peerAddr := p.RemoteAddr().String()
		s.peer = tcp.NewTCPConnection(
			connID,
			s.name+fmt.Sprintf("-%v#%v", peerAddr, connID),
			c,
			conn.KClient,
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

func (s *TCPClient) onProtocol(proto string) transmit.IChannel {
	switch proto {
	case "tcp":
		return tcp_channel.NewChannel()
	case "ws":
		return ws_channel.NewChannel()
	}
	panic(fmt.Sprintf("no proto setup"))
}

func (s *TCPClient) ConnectTCP(address string) {
	s.c.ConnectTCP(address)
}

func (s *TCPClient) OnConnection(peer conn.Session) {
	if peer.Connected() {
		log.Print("--- *** TCPClient - TCPClient:: OnConnected \n")
	} else {
		log.Print("--- *** TCPClient - TCPClient:: OnClosed \n")
	}
}

func (s *TCPClient) OnMessage(peer conn.Session, msg interface{}, recvTime utils.Timestamp) {
	log.Print("--- *** TCPClient - TCPClient:: OnMessage \n")
}

func (s *TCPClient) OnWriteComplete(peer conn.Session) {
	log.Print("--- *** TCPClient - TCPClient:: OnWriteComplete \n")
}

func (s *TCPClient) removeConnection(peer conn.Session) {
	sessions.Remove(peer)
	peer.(*tcp.TCPConnection).ConnectDestroyed()
}

func (s *TCPClient) onConnectionError(err error) {
	log.Print("--- *** TCPClient - TCPClient:: onConnectionError \n")
}

func (s *TCPClient) Reconnect(d time.Duration) {
	s.c.Reconnect(d)
}

func (s *TCPClient) Disconnect() {
	if s.peer != nil {
		s.peer.Close()
	}
}
