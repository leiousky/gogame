package ws_channel

import (
	"errors"
	"games/comm/utils"
	"games/core/conn/transmit"

	"github.com/gorilla/websocket"
)

/// <summary>
/// Channel WS传输
/// <summary>
type Channel struct {
}

func NewChannel() transmit.IChannel {
	return &Channel{}
}

func (s *Channel) OnRecv(conn interface{}) (interface{}, error) {
	c, ok := conn.(*websocket.Conn)
	if !ok || c == nil {
		panic(errors.New("ws_channel.OnRecv conn == nil"))
	}
	c.SetReadLimit(512)
	msgType, buf, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}
	//TextMessage/BinaryMessage
	if websocket.BinaryMessage != msgType {
		return nil, nil
	}
	return buf, err
}

func (s *Channel) OnSend(conn interface{}, msg interface{}) error {
	c, ok := conn.(*websocket.Conn)
	if !ok || c == nil {
		panic(errors.New("ws_channel.OnSend conn == nil"))
	}
	buf, _ := utils.ToBytes(msg)
	return c.WriteMessage(websocket.BinaryMessage, buf)
}
