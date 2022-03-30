package tcp

import (
	"fmt"
	cb "games/core/callback"
	"log"
	"net"
)

/// <summary>
/// Acceptor TCP接受器
/// <summary>
type Acceptor interface {
	ListenTCP(address string)
	Close()
	SetProtocolCallback(cb cb.OnProtocol)
	SetConditionCallback(cb cb.OnCondition)
	SetNewConnectionCallback(cb cb.OnNewConnection)
}

/// <summary>
/// acceptor TCP接受器
/// <summary>
type acceptor struct {
	listener        interface{}
	onProtocol      cb.OnProtocol
	onCondition     cb.OnCondition
	onNewConnection cb.OnNewConnection
}

func NewAcceptor() Acceptor {
	s := &acceptor{}
	return s
}

func (s *acceptor) SetProtocolCallback(cb cb.OnProtocol) {
	s.onProtocol = cb
}

func (s *acceptor) SetConditionCallback(cb cb.OnCondition) {
	s.onCondition = cb
}

func (s *acceptor) SetNewConnectionCallback(cb cb.OnNewConnection) {
	s.onNewConnection = cb
}

func (s *acceptor) ListenTCP(address string) {

}

func (s *acceptor) listenTCP(address string) int {
	if s.onProtocol == nil {
		panic(fmt.Sprintf("listenTCP s.onProtocol == nil"))
	}
	if s.onCondition == nil {
		panic(fmt.Sprintf("listenTCP s.onCondition == nil"))
	}
	if s.onNewConnection == nil {
		panic(fmt.Sprintf("listenTCP s.onNewConnection == nil"))
	}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println(err)
		return 1
	}
	//channel := s.onProtocol("tcp")
	s.SetConditionCallback(s.onCondition)
	s.listener = listener
	return 0
}

func (s *acceptor) listenWS(address string) int {
	return 0
}

func (s *acceptor) Close() {
	if l, ok := s.listener.(net.Listener); ok {
		l.Close()
	}

}
