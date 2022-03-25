package transmit

import (
	"io"
	"net"
)

/// <summary>
/// IChannel 消息传输接口
/// <summary>
type IChannel interface {
	OnRecvMessage(conn interface{}) (msg interface{}, err error)
	OnSendMessage(conn interface{}, msg interface{}) error
}

func IsEOFOrReadError(err error) bool {
	if err == io.EOF {
		return true
	}
	ne, ok := err.(*net.OpError)
	return ok && ne.Op == "read"
}

func IsEOFOrWriteError(err error) bool {
	if err == io.EOF {
		return true
	}
	ne, ok := err.(*net.OpError)
	return ok && ne.Op == "write"
}
