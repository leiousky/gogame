package ws_channel

import (
	"log"

	"github.com/gorilla/websocket"
)

/// 读指定长度
func ReadFull(conn *websocket.Conn) (buf []byte, err error) {
	length := 0
	size := len(buf)
	for {
		conn.SetReadLimit(int64(size - length))
		_, b, e := conn.ReadMessage()
		if err != nil {
			err = e
			log.Print("ws_channel.ReadFull : ", err)
			return
		}
		n := len(b)
		//copy(buf[length:], n)
		buf = append(buf, b[:]...)
		length += n
		if length == size {
			return
		}
	}
}
