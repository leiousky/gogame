package tcp_server

import (
	"errors"
	"games/comm/utils"
	"games/core/conn"
	"games/core/conn/tcp/tcpserver"
	"games/core/conn/transmit"
	"games/server/stream/tcp_stream"
	"games/server/stream/ws_stream"
	"log"
)

/// <summary>
/// TCPServer TCP服务端
/// <summary>
type TCPServer struct {
	tcpserver.ITCPServer
	c tcpserver.ITCPServer
}

func NewTCPServer(name string) tcpserver.ITCPServer {
	s := &TCPServer{c: tcpserver.NewTCPServer(name)}
	s.c.(*tcpserver.TCPServer).SetProtocolCallback(s.onProtocol)
	s.c.(*tcpserver.TCPServer).SetConditionCallback(s.OnCondition)
	s.c.(*tcpserver.TCPServer).SetConnectionCallback(s.OnConnection)
	s.c.(*tcpserver.TCPServer).SetMessageCallback(s.OnMessage)
	s.c.(*tcpserver.TCPServer).SetWriteCompleteCallback(s.OnWriteComplete)
	return s
}

func (s *TCPServer) onProtocol(proto string) transmit.IChannel {
	switch proto {
	case "tcp":
		return tcp_stream.NewChannel()
	case "ws":
		return ws_stream.NewChannel()
	}
	panic(errors.New("no proto setup"))
}

func (s *TCPServer) ListenTCP(address string) {
	s.c.ListenTCP(address)
}

func (s *TCPServer) OnCondition(c interface{}) bool {
	return true
}

func (s *TCPServer) OnConnection(peer conn.Session) {
	if peer.Connected() {
		log.Print("--- *** TCPServer - TCPServer:: OnConnected \n")
	} else {
		log.Print("--- *** TCPServer - TCPServer:: OnClosed \n")
	}
}

func (s *TCPServer) OnMessage(peer conn.Session, msg interface{}, recvTime utils.Timestamp) {
	log.Print("--- *** TCPServer - TCPServer:: OnMessage \n")
}

func (s *TCPServer) OnWriteComplete(peer conn.Session) {
	log.Print("--- *** TCPServer - TCPServer:: OnWriteComplete \n")
}
