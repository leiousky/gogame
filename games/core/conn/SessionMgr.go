package conn

import (
	"sync"
)

/// <summary>
/// ISessions 连接会话容器
/// <summary>
type ISessions interface {
	Add(peer Session) bool
	Remove(peer Session)
	Get(sesID int64) Session
	Count() int
	Stop()
	Wait()
}

/// <summary>
/// Sessions 连接会话容器
/// <summary>
type Sessions struct {
	peers map[int64]Session
	l     *sync.Mutex
	c     *sync.Cond
	stop  bool
}

func NewSessions() ISessions {
	s := &Sessions{l: &sync.Mutex{}, peers: map[int64]Session{}}
	s.c = sync.NewCond(s.l)
	return s
}

func (s *Sessions) Add(peer Session) bool {
	ok := false
	s.l.Lock()
	if !s.stop {
		s.peers[peer.ID()] = peer
		ok = true
	}
	s.l.Unlock()
	return ok
}

func (s *Sessions) Remove(peer Session) {
	s.l.Lock()
	if _, ok := s.peers[peer.ID()]; ok {
		delete(s.peers, peer.ID())
	}
	if s.stop && len(s.peers) == 0 {
		s.c.Signal()
	}
	s.l.Unlock()
}

func (s *Sessions) Get(sesID int64) Session {
	s.l.Lock()
	if peer, ok := s.peers[sesID]; ok {
		s.l.Unlock()
		return peer
	}
	s.l.Unlock()
	return nil
}

func (s *Sessions) Count() int {
	c := 0
	s.l.Lock()
	c = len(s.peers)
	s.l.Unlock()
	return c
}

func (s *Sessions) Stop() {
	s.l.Lock()
	s.stop = true
	for _, peer := range s.peers {
		peer.Close()
	}
	s.peers = map[int64]Session{}
	s.l.Unlock()
}

func (s *Sessions) Wait() {
	s.l.Lock()
	s.c.Wait()
	s.l.Unlock()
}
