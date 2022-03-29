package net

import (
	"fmt"
	"games/core/net/transmit"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

/// <summary>
/// ITCPClient TCP/WS
/// <summary>
type ITCPClient interface {
	Session() Session
	Write(msg interface{})
	ConnectTCP(name, address string)
	OnConnected(peer Session)
	OnMessage(msg interface{}, peer Session)
	OnClosed(peer Session)
	Reconnect(d time.Duration)
	Disconnect()
}

/// <summary>
/// Clients TCP/WS
/// <summary>
var Clients = NewSessions()

/// <summary>
/// TCPClient TCP/WS
/// <summary>
type TCPClient struct {
	peer Session
}

func NewTCPClient() ITCPClient {
	s := &TCPClient{}
	return s
}

func (s *TCPClient) Session() Session {
	return s.peer
}

func (s *TCPClient) Write(msg interface{}) {
	if s.peer != nil {
		s.peer.Write(msg)
	}
}

func (s *TCPClient) connectTCP(name, address string) int {
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	//c, ok := conn.(net.Conn); ok
	s.peer = NewTCPConnection(
		NewConnID(), name, conn,
		SesClient, transmit.NewTCPChannel())
	s.peer.SetOnConnected(s.OnConnected)
	s.peer.SetOnMessage(s.OnMessage)
	s.peer.SetOnClosed(s.OnClosed)
	s.peer.SetOnWritten(s.OnWritten)
	s.peer.SetOnError(s.OnError)
	if !Clients.Add(s.peer) {
		s.peer.Close()
	}
	return 0
}

func (s *TCPClient) connectWS(name, address string) int {
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
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	//c, ok := conn.(*websocket.Conn); ok
	s.peer = NewTCPConnection(
		NewConnID(), name, conn,
		SesClient, transmit.NewWSChannel())
	s.peer.SetOnConnected(s.OnConnected)
	s.peer.SetOnMessage(s.OnMessage)
	s.peer.SetOnClosed(s.OnClosed)
	s.peer.SetOnWritten(s.OnWritten)
	s.peer.SetOnError(s.OnError)
	if !Clients.Add(s.peer) {
		s.peer.Close()
	}
	return 0
}

func (s *TCPClient) ConnectTCP(name, address string) {
	if s.connectWS(name, address) == -1 {
		s.connectTCP(name, address)
	}
}

func (s *TCPClient) OnConnected(peer Session) {
}

func (s *TCPClient) OnClosed(peer Session) {
}

func (s *TCPClient) OnMessage(msg interface{}, peer Session) {
}

func (s *TCPClient) OnWritten(msg interface{}, peer Session) {
}

func (s *TCPClient) OnError(peer Session, err error) {
}

func (s *TCPClient) Reconnect(d time.Duration) {
}

func (s *TCPClient) Disconnect() {
	if s.peer != nil {
		s.peer.Close()
	}
}
