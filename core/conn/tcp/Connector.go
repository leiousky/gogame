package tcp

import (
	"errors"
	cb "games/core/callback"
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
	Retry() bool
	EnableRetry(bool)
	ServerAddr() string
	ConnectTCP(address string) string
	Reconnect(d time.Duration)
	SetProtocolCallback(cb cb.OnProtocol)
	SetNewConnectionCallback(cb cb.OnNewConnection)
}

/// <summary>
/// Connector TCP连接器
/// <summary>
type connector struct {
	name            string
	address         string
	addrType        string
	retry           bool
	d               time.Duration
	onProtocol      cb.OnProtocol
	onNewConnection cb.OnNewConnection
}

func NewConnector(name string) Connector {
	s := &connector{name: name}
	return s
}

func (s *connector) Retry() bool {
	return s.retry
}

func (s *connector) EnableRetry(bv bool) {
	s.retry = bv
}

func (s *connector) ServerAddr() string {
	return s.address
}

func (s *connector) SetProtocolCallback(cb cb.OnProtocol) {
	s.onProtocol = cb
}

func (s *connector) SetNewConnectionCallback(cb cb.OnNewConnection) {
	s.onNewConnection = cb
}

func (s *connector) connectTCPTimeout(address string, d time.Duration) int {
	if s.onProtocol == nil {
		panic(errors.New("connector.connectTCP s.onProtocol == nil"))
	}
	if s.onNewConnection == nil {
		panic(errors.New("connector.connectTCP s.onNewConnection == nil"))
	}
	//log.Printf("ConnectTCP %v://%v\n", "tcp", address)
	c, err := net.DialTimeout("tcp", address, d)
	if err != nil {
		//log.Println(err)
		return 1
	}
	s.address = address
	s.addrType = "tcp"
	channel := s.onProtocol("tcp")
	s.onNewConnection(c, channel)
	return 0
}

func (s *connector) connectWSTimeout(address string, d time.Duration) int {
	if s.onProtocol == nil {
		panic(errors.New("connector.connectWS s.onProtocol == nil"))
	}
	if s.onNewConnection == nil {
		panic(errors.New("connector.connectWS s.onNewConnection == nil"))
	}
	//ws://ip:port wss://ip:port
	vec := strings.Split(address, "//")
	if len(vec) != 2 {
		return -1
	}
	dialer := websocket.Dialer{}
	dialer.Proxy = http.ProxyFromEnvironment
	dialer.HandshakeTimeout = d
	proto := strings.Trim(vec[0], ":")
	host := vec[1]
	u := url.URL{Scheme: proto, Host: host, Path: "/"}
	//log.Printf("ConnectTCP %v://%v\n", proto, host)
	c, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		//log.Println(err)
		return 1
	}
	s.address = address
	s.addrType = "ws"
	channel := s.onProtocol("ws")
	s.onNewConnection(c, channel)
	return 0
}

func (s *connector) connectTimeout(address string, d time.Duration) string {
	if s.connectWSTimeout(address, d) == -1 {
		s.connectTCPTimeout(address, d)
		return "tcp"
	}
	return "ws"
}

func (s *connector) ConnectTCP(address string) string {
	return s.connectTimeout(address, time.Second)
}

func (s *connector) reconnect() {
	log.Printf("--- *** connector[%v] - Reconnecting to %v \n", s.name, s.address)
	switch s.addrType {
	case "tcp":
		if 1 == s.connectTCPTimeout(s.address, time.Second) && s.retry {
			time.AfterFunc(s.d, s.reconnect)
		}
		break
	case "ws":
		if 1 == s.connectWSTimeout(s.address, time.Second) && s.retry {
			time.AfterFunc(s.d, s.reconnect)
		}
		break
	}
}

func (s *connector) Reconnect(d time.Duration) {
	s.d = d
	time.AfterFunc(s.d, s.reconnect)
}
