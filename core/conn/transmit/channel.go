package transmit

import (
	"games/core/conn/def"
	"io"
	"net"
)

/// <summary>
/// IChannel 消息传输接口
/// <summary>
type IChannel interface {
	/// 接收数据
	OnRecv(conn interface{}) (interface{}, error, def.Reason)
	/// 发送数据
	OnSend(conn interface{}, msg interface{}) (error, def.Reason)
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
