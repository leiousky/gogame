package tcp_channel

import (
	"games/comm/utils"
	"games/core/conn/transmit"

	"net"
)

/// <summary>
/// Channel TCP传输
/// <summary>
type Channel struct {
}

func NewChannel() transmit.IChannel {
	return &Channel{}
}

func (s *Channel) OnRecv(conn interface{}) (msg interface{}, err error) {
	c, ok := conn.(net.Conn)
	if !ok || c == nil {
		return nil, nil
	}
	buf := make([]byte, 512)
	err = ReadFull(c, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (s *Channel) OnSend(conn interface{}, msg interface{}) error {
	c, ok := conn.(net.Conn)
	if !ok || c == nil {
		return nil
	}
	buf, _ := utils.ToBytes(msg)
	return WriteFull(c, buf)
}
