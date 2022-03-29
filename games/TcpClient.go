package main

import (
	"games/core/net"
	"log"
)

type TCPClient struct {
	net.TCPClient
	c *net.Connector
}

func NewTCPClient() net.TCPClient {
	return &TCPClient{
		c: net.NewConnector()}
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
