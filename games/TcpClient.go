package main

import (
	"games/core/conn"
	"games/core/conn/tcp"
	"log"
)

/// <summary>
/// TCPClient TCP客户端
/// <summary>
type TCPClient struct {
	tcp.TCPClient
	c tcp.Connector
}

func NewTCPClient() tcp.TCPClient {
	s := &TCPClient{
		c: tcp.NewConnector()}
	s.c.SetOnConnected(s.OnConnected)
	s.c.SetOnMessage(s.OnMessage)
	s.c.SetOnClosed(s.OnClosed)
	s.c.SetOnWritten(s.OnWritten)
	s.c.SetOnError(s.OnError)
	return s
}

func (s *TCPClient) ConnectTCP(name string, address string) {
	s.c.ConnectTCP(name, address)
}

func (s *TCPClient) OnConnected(peer conn.Session) {
	log.Print("--- *** TCPClient - TCPClient:: OnConnected \n")
}

func (s *TCPClient) OnMessage(msg interface{}, peer conn.Session) {
	log.Print("--- *** TCPClient - TCPClient:: OnMessage \n")
}

func (s *TCPClient) OnClosed(peer conn.Session) {
	log.Print("--- *** TCPClient - TCPClient:: OnClosed \n")
}

func (s *TCPClient) OnWritten(msg interface{}, peer conn.Session) {
	log.Print("--- *** TCPClient - TCPClient:: OnWritten \n")
}

func (s *TCPClient) OnError(peer conn.Session, err error) {
	log.Print("--- *** TCPClient - TCPClient:: OnError \n")
}
