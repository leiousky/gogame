package tcpconns

import (
	"errors"
	"games/core/conn"
	"games/core/conn/tcp/tcpclient"
	"games/core/conn/tcp/tcpserver"
)

func Get(connType conn.Type, id int64) conn.Session {
	switch connType {
	case conn.KServer:
		return tcpserver.Get(id)
	case conn.KClient:
		return tcpclient.Get(id)
	}
	panic(errors.New("connType error"))
}

func Count(connType conn.Type) int {
	switch connType {
	case conn.KServer:
		return tcpserver.Count()
	case conn.KClient:
		return tcpclient.Count()
	}
	panic(errors.New("connType error"))
}

func CloseAll() {
	tcpclient.CloseAll()
	tcpserver.CloseAll()
}

func Wait() {
	tcpclient.Wait()
	tcpserver.Wait()
}

func Stop() {
	tcpclient.Stop()
	tcpserver.Stop()
}
