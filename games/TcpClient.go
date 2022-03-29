package main

import (
	"games/core/conn/tcp"
	"log"
)

type TCPClient struct {
	tcp.TCPClient
	c *tcp.Connector
}

func NewTCPClient() tcp.TCPClient {
	return &TCPClient{
		c: tcp.NewConnector()}
}

func (s *TCPClient) OnConnected(peer interface{}) {
	log.Print("--- *** TCPClient - TCPClient:: OnConnected \n")
}

func (s *TCPClient) OnMessage(msg interface{}, peer interface{}) {
	log.Print("--- *** TCPClient - TCPClient:: OnMessage \n")
}

func (s *TCPClient) OnClosed(peer interface{}) {
	log.Print("--- *** TCPClient - TCPClient:: OnClosed \n")
}

func (s *TCPClient) OnWritten(msg interface{}, peer interface{}) {
}

func (s *TCPClient) OnError(peer interface{}, err error) {
	log.Print("--- *** TCPClient - TCPClient:: OnError \n")
}
