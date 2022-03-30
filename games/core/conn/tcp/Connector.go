package tcp

import (
	"fmt"
	cb "games/core/callback"
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
	ConnectTCP(address string) string
	Reconnect(d time.Duration)
	SetProtocolCallback(cb cb.OnProtocol)
	SetNewConnectionCallback(cb cb.OnNewConnection)
}

/// <summary>
/// Connector TCP连接器
/// <summary>
type connector struct {
	onProtocol      cb.OnProtocol
	onNewConnection cb.OnNewConnection
}

func NewConnector() Connector {
	s := &connector{}
	return s
}

func (s *connector) SetProtocolCallback(cb cb.OnProtocol) {
	s.onProtocol = cb
}

func (s *connector) SetNewConnectionCallback(cb cb.OnNewConnection) {
	s.onNewConnection = cb
}

func (s *connector) connectTCP(address string) int {
	if s.onProtocol == nil {
		panic(fmt.Sprintf("connectTCP s.onProtocol == nil"))
	}
	if s.onNewConnection == nil {
		panic(fmt.Sprintf("connectTCP s.onNewConnection == nil"))
	}
	c, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	channel := s.onProtocol("tcp")
	s.onNewConnection(c, channel)
	return 0
}

func (s *connector) connectWS(address string) int {
	if s.onProtocol == nil {
		panic(fmt.Sprintf("connectWS s.onProtocol == nil"))
	}
	if s.onNewConnection == nil {
		panic(fmt.Sprintf("connectWS s.onNewConnection == nil"))
	}
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
	channel := s.onProtocol("ws")
	s.onNewConnection(c, channel)
	return 0
}

func (s *connector) ConnectTCP(address string) string {
	if s.connectWS(address) == -1 {
		s.connectTCP(address)
		return "tcp"
	}
	return "ws"
}

func (s *connector) Reconnect(d time.Duration) {
}
