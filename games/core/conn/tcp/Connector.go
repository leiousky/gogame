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
/// sessions 客户端会话容器
/// <summary>
var sessions = conn.NewSessions()

/// <summary>
/// Connector 客户端连接器
/// <summary>
type Connector struct {
	peer        conn.Session
	onConnected cb.OnConnected
	onMessage   cb.OnMessage
	onClosed    cb.OnClosed
	onWritten   cb.OnWritten
	onError     cb.OnError
}

func NewConnector() *Connector {
	s := &Connector{}
	return s
}

func (s *Connector) SetOnConnected(cb cb.OnConnected) {
	s.onConnected = cb
}

func (s *Connector) SetOnMessage(cb cb.OnMessage) {
	s.onMessage = cb
}

func (s *Connector) SetOnClosed(cb cb.OnClosed) {
	s.onClosed = cb
}

func (s *Connector) SetOnWritten(cb cb.OnWritten) {
	s.onWritten = cb
}

func (s *Connector) SetOnError(cb cb.OnError) {
	s.onError = cb
}

func (s *Connector) Session() conn.Session {
	return s.peer
}

func (s *Connector) Write(msg interface{}) {
	if s.peer != nil {
		s.peer.Write(msg)
	}
}

func (s *Connector) connectTCP(name, address string) int {
	c, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	s.peer = NewTCPConnection(
		conn.NewConnID(), name, c,
		conn.ClientType, transmit.NewTCPChannel())
	s.peer.SetOnConnected(s.onConnected)
	s.peer.SetOnMessage(s.onMessage)
	s.peer.SetOnClosed(s.onClosed)
	s.peer.SetOnWritten(s.onWritten)
	s.peer.SetOnError(s.onError)
	if !sessions.Add(s.peer) {
		s.peer.Close()
	}
	return 0
}

func (s *Connector) connectWS(name, address string) int {
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
	s.peer.SetOnConnected(s.onConnected)
	s.peer.SetOnMessage(s.onMessage)
	s.peer.SetOnClosed(s.onClosed)
	s.peer.SetOnWritten(s.onWritten)
	s.peer.SetOnError(s.onError)
	if !sessions.Add(s.peer) {
		s.peer.Close()
	}
	return 0
}

func (s *Connector) ConnectTCP(name, address string) {
	if s.connectWS(name, address) == -1 {
		s.connectTCP(name, address)
	}
}

func (s *Connector) Reconnect(d time.Duration) {
}

func (s *Connector) Disconnect() {
	if s.peer != nil {
		s.peer.Close()
	}
}
