package main

import (
	"games/comm/utils"
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
	s.c.SetConnectionCallback(s.OnConnection)
	s.c.SetMessageCallback(s.OnMessage)
	s.c.SetWriteCompleteCallback(s.OnWriteComplete)
	return s
}

func (s *TCPClient) ConnectTCP(name string, address string) {
	s.c.ConnectTCP(name, address)
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
