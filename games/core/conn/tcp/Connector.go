package tcp

import (
	"fmt"
	"games/core/cb"
	"games/core/conn"
	"games/core/conn/transmit"
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
	SetOnConnected(cb cb.OnConnected)
	SetOnMessage(cb cb.OnMessage)
	SetOnClosed(cb cb.OnClosed)
	SetOnWritten(cb cb.OnWritten)
	SetOnError(cb cb.OnError)
}

/// <summary>
/// sessions 客户端容器
/// <summary>
var sessions = conn.NewSessions()

/// <summary>
/// Connector TCP连接器
/// <summary>
type connector struct {
	peer        conn.Session
	onConnected cb.OnConnected
	onMessage   cb.OnMessage
	onClosed    cb.OnClosed
	onWritten   cb.OnWritten
	onError     cb.OnError
}

func NewConnector() Connector {
	s := &connector{}
	return s
}

func (s *connector) SetOnConnected(cb cb.OnConnected) {
	s.onConnected = cb
}

func (s *connector) SetOnMessage(cb cb.OnMessage) {
	s.onMessage = cb
}

func (s *connector) SetOnClosed(cb cb.OnClosed) {
	s.onClosed = cb
}

func (s *connector) SetOnWritten(cb cb.OnWritten) {
	s.onWritten = cb
}

func (s *connector) SetOnError(cb cb.OnError) {
	s.onError = cb
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
		conn.ClientType, transmit.NewTCPChannel())
	s.peer.(*TCPConnection).SetOnConnected(s.onConnected)
	s.peer.(*TCPConnection).SetOnMessage(s.onMessage)
	s.peer.(*TCPConnection).SetOnClosed(s.onClosed)
	s.peer.(*TCPConnection).SetOnWritten(s.onWritten)
	s.peer.(*TCPConnection).SetOnError(s.onError)
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
		conn.ClientType, transmit.NewWSChannel())
	s.peer.(*TCPConnection).SetOnConnected(s.onConnected)
	s.peer.(*TCPConnection).SetOnMessage(s.onMessage)
	s.peer.(*TCPConnection).SetOnClosed(s.onClosed)
	s.peer.(*TCPConnection).SetOnWritten(s.onWritten)
	s.peer.(*TCPConnection).SetOnError(s.onError)
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

func (s *connector) Reconnect(d time.Duration) {
}

func (s *connector) Disconnect() {
	if s.peer != nil {
		s.peer.Close()
	}
}
