package tcp_server

import "games/core/conn/tcp"

/// <summary>
/// TCPServer TCP服务端
/// <summary>
type TCPServer struct {
	tcp.TCPServer
	c tcp.Acceptor
}
