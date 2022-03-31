package tcp_client

import (
	"errors"
	"games/comm/utils"
	"games/core/conn"
	"games/core/conn/tcp/tcpclient"
	"games/core/conn/transmit"
	"games/server/stream/tcp_stream"
	"games/server/stream/ws_stream"
	"log"
)

/// <summary>
/// TCPClient TCP客户端
/// <summary>
type TCPClient struct {
	tcpclient.ITCPClient
	c tcpclient.TCPClient
}

func NewTCPClient() tcpclient.ITCPClient {
	s := &TCPClient{}
	s.c.SetProtocolCallback(s.onProtocol)
	s.c.SetConnectionCallback(s.OnConnection)
	s.c.SetMessageCallback(s.OnMessage)
	s.c.SetWriteCompleteCallback(s.OnWriteComplete)
	return s
}

func (s *TCPClient) onProtocol(proto string) transmit.IChannel {
	switch proto {
	case "tcp":
		return tcp_stream.NewChannel()
	case "ws":
		return ws_stream.NewChannel()
	}
	panic(errors.New("no proto setup"))
}

func (s *TCPClient) ConnectTCP(address string) {
	s.c.ConnectTCP(address)
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
