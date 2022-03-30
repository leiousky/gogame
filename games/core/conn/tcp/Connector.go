package tcp

import (
	"fmt"
	cb "games/core/callback"
	"games/core/conn"
	"games/core/conn/transmit"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

/// <summary>
/// Connector TCP连接器
/// <summary>
type Connector interface {
	Session() conn.Session
	Write(msg interface{})
	ConnectTCP(name, address string)
	Reconnect(d time.Duration)
	Disconnect()
	SetConnectionCallback(cb cb.OnConnection)
	SetMessageCallback(cb cb.OnMessage)
	SetWriteCompleteCallback(cb cb.OnWriteComplete)
}

/// <summary>
/// sessions 客户端容器
/// <summary>
var sessions = conn.NewSessions()

/// <summary>
/// Connector TCP连接器
/// <summary>
type connector struct {
	peer            conn.Session
	onConnection    cb.OnConnection
	onMessage       cb.OnMessage
	onWriteComplete cb.OnWriteComplete
	closeCallback   cb.CloseCallback
	errorCallback   cb.ErrorCallback
}

func NewConnector() Connector {
	s := &connector{}
	return s
}

func (s *connector) SetConnectionCallback(cb cb.OnConnection) {
	s.onConnection = cb
}

func (s *connector) SetMessageCallback(cb cb.OnMessage) {
	s.onMessage = cb
}

func (s *connector) SetWriteCompleteCallback(cb cb.OnWriteComplete) {
	s.onWriteComplete = cb
}

func (s *connector) SetCloseCallback(cb cb.CloseCallback) {
	s.closeCallback = cb
}

func (s *connector) SetErrorCallback(cb cb.ErrorCallback) {
	s.errorCallback = cb
}

func (s *connector) Session() conn.Session {
	return s.peer
}

func (s *connector) Write(msg interface{}) {
	if s.peer != nil {
		s.peer.Write(msg)
	}
}

func (s *connector) connectTCP(name, address string) int {
	c, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	s.peer = NewTCPConnection(
		conn.NewConnID(), name, c,
		conn.KClient, transmit.NewTCPChannel())
	s.peer.(*TCPConnection).SetConnectionCallback(s.onConnection)
	s.peer.(*TCPConnection).SetMessageCallback(s.onMessage)
	s.peer.(*TCPConnection).SetWriteCompleteCallback(s.onWriteComplete)
	s.peer.(*TCPConnection).SetCloseCallback(s.removeConnection)
	s.peer.(*TCPConnection).SetErrorCallback(s.onConnectionError)
	s.peer.(*TCPConnection).ConnectEstablished()
	if !sessions.Add(s.peer) {
		s.peer.Close()
	}
	return 0
}

func (s *connector) connectWS(name, address string) int {
	//ws://ip:port wss://ip:port
	vec := strings.Split(address, "//")
	if len(vec) != 2 {
		return -1
	}
	dialer := websocket.Dialer{}
	dialer.Proxy = http.ProxyFromEnvironment
	dialer.HandshakeTimeout = 3 * time.Second
	proto := strings.Trim(vec[0], ":")
	host := vec[1]
	//log.Printf("ConnectTCP %v://%v\n", proto, host)
	u := url.URL{Scheme: proto, Host: host, Path: "/"}
	c, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	s.peer = NewTCPConnection(
		conn.NewConnID(), name, c,
		conn.KClient, transmit.NewWSChannel())
	s.peer.(*TCPConnection).SetConnectionCallback(s.onConnection)
	s.peer.(*TCPConnection).SetMessageCallback(s.onMessage)
	s.peer.(*TCPConnection).SetWriteCompleteCallback(s.onWriteComplete)
	s.peer.(*TCPConnection).SetCloseCallback(s.removeConnection)
	s.peer.(*TCPConnection).SetErrorCallback(s.onConnectionError)
	s.peer.(*TCPConnection).ConnectEstablished()
	if !sessions.Add(s.peer) {
		s.peer.Close()
	}
	return 0
}

func (s *connector) ConnectTCP(name, address string) {
	if s.connectWS(name, address) == -1 {
		s.connectTCP(name, address)
	}
}

func (s *connector) removeConnection(peer conn.Session) {
	sessions.Remove(peer)
	peer.(*TCPConnection).ConnectDestroyed()
}

func (s *connector) onConnectionError(err error) {
	log.Print("--- *** connector - connector:: onConnectionError \n")
}

func (s *connector) Reconnect(d time.Duration) {
}

func (s *connector) Disconnect() {
	if s.peer != nil {
		s.peer.Close()
	}
}
