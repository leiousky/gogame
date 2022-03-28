package net

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

/// <summary>
/// TCPClient TCP/WS客户端
/// <summary>
type ITCPClient interface {
	//会话ID
	ID() int64
	//会话
	Session() Session
	//关闭
	Close()
	//写
	Write(msg interface{})
	//连接
	ConnectTCP(address string)
}

type TCPClient struct {
	ses Session
}

func (s *TCPClient) ID() int64 {
	if s.ses != nil {
		return s.ses.ID()
	}
	return int64(0)
}

func (s *TCPClient) Session() Session {
	return s.ses
}

func (s *TCPClient) Close() {
	if s.ses != nil {
		s.ses.Close()
	}
}

func (s *TCPClient) Write(msg interface{}) {
	if s.ses != nil {
		s.ses.Write(msg)
	}
}

func (s *TCPClient) ConnectTCP(address string) {
	if s.connectWS(address) == -1 {
		s.connectTCP(address)
	}
}

func (s *TCPClient) connectTCP(address string) int {
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	peer := gSessMgr.Add(conn, SesClient)
	if peer != nil {
		// peer.SetOnConnected(s.onConnected)
		// peer.SetOnClosed(s.onClosed)
		// peer.SetOnMessage(s.onMessage)
		// peer.SetOnError(s.onError)
		// peer.SetCloseCallback(s.remove)
		// peer.OnEstablished()
	} else {
		conn.Close()
	}
	return 0
}

//
func (s *TCPClient) connectWS(address string) int {
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
	peer := gSessMgr.Add(conn, SesClient)
	if peer != nil {
		// peer.SetOnConnected(s.onConnected)
		// peer.SetOnClosed(s.onClosed)
		// peer.SetOnMessage(s.onMessage)
		// peer.SetOnError(s.onError)
		// peer.SetOnWritten(s.onWritten)
		// peer.SetCloseCallback(s.remove)
		// peer.OnEstablished()
	} else {
		conn.Close()
	}
	return 0
}
