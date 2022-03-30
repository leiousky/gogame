package tcp

import cb "games/core/callback"

/// <summary>
/// Acceptor TCP接受器
/// <summary>
type Acceptor interface {
	SetProtocolCallback(cb cb.OnProtocol)
	SetConditionCallback(cb cb.OnCondition)
	SetNewConnectionCallback(cb cb.OnNewConnection)
}

/// <summary>
/// acceptor TCP接受器
/// <summary>
type acceptor struct {
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
