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
	"time"
)

/// <summary>
/// TCPClient TCP客户端
/// <summary>
type TCPClient struct {
	c tcpclient.TCPClient
}

func NewTCPClient(name string) tcpclient.TCPClient {
	s := &TCPClient{c: tcpclient.NewTCPClient(name)}
	s.c.(*tcpclient.Processor).SetProtocolCallback(s.onProtocol)
	s.c.(*tcpclient.Processor).SetConnectionCallback(s.OnConnection)
	s.c.(*tcpclient.Processor).SetMessageCallback(s.OnMessage)
	s.c.(*tcpclient.Processor).SetWriteCompleteCallback(s.OnWriteComplete)
	s.c.(*tcpclient.Processor).EnableRetry(true)
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

func (s *TCPClient) Session() conn.Session {
	return s.c.Session()
}

func (s *TCPClient) Write(msg interface{}) {
	s.c.Write(msg)
}

func (s *TCPClient) ConnectTCP(address string) {
	s.c.ConnectTCP(address)
}

func (s *TCPClient) OnConnection(peer conn.Session, reason conn.Reason) {
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

func (s *TCPClient) Reconnect(d time.Duration) {
	s.c.Reconnect(d)
}

func (s *TCPClient) Disconnect() {
	s.c.Disconnect()
}
