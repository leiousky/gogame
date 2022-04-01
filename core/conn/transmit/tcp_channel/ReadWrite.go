package tcp_channel

import (
	"log"
	"net"
)

/// 读指定长度
func ReadFull(conn net.Conn, buf []byte) error {
	length := 0
	size := len(buf)
	for {
		n, err := conn.Read(buf[length:size])
		if err != nil {
			//log.Print("tcp_channel.ReadFull : ", err)
			return err
		}
		length += n
		if length == size {
			return nil
		}
	}
}

/// 写指定长度
func WriteFull(conn net.Conn, buf []byte) error {
	length := 0
	size := len(buf)
	for {
		n, err := conn.Write(buf[length:size])
		if err != nil {
			log.Print("tcp_channel.WriteFull : ", err)
			return err
		}
		length += n
		if length == size {
			return nil
		}
	}
}
