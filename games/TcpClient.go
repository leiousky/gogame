package main

import (
	"games/core/conn"
	"games/core/conn/tcp"
	"log"
)

type TCPClient struct {
	tcp.TCPClient
	c tcp.Connector
}

func NewTCPClient() tcp.TCPClient {
	return &TCPClient{
		c: tcp.NewConnector()}
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
}

func (s *TCPClient) OnError(peer conn.Session, err error) {
	log.Print("--- *** TCPClient - TCPClient:: OnError \n")
}
