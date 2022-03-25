package transmit

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
			log.Print("ReadFull : ", err)
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
			log.Print("WriteFull : ", err)
			return err
		}
		length += n
		if length == size {
			return nil
		}
	}
}
