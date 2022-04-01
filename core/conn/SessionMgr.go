package conn

import (
	"log"
	"sync"
)

/// <summary>
/// ISessions 连接会话容器
/// <summary>
type ISessions interface {
	Get(sesID int64) Session
	Count() int
	Add(peer Session) bool
	Remove(peer Session)
	Wait()
	Stop()
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

func (s *Sessions) Add(peer Session) bool {
	ok := false
	s.l.Lock()
	if !s.stop {
		log.Printf("Sessions.Add => %v", peer.Name())
		s.peers[peer.ID()] = peer
		ok = true
	}
	s.l.Unlock()
	return ok
}

func (s *Sessions) Remove(peer Session) {
	s.l.Lock()
	if _, ok := s.peers[peer.ID()]; ok {
		log.Printf("Sessions.Remove => %v", peer.Name())
		delete(s.peers, peer.ID())
	}
	s.l.Unlock()
}

func (s *Sessions) removeAll(stop bool) {
	s.l.Lock()
	for _, peer := range s.peers {
		peer.Close()
	}
	s.peers = map[int64]Session{}
	if stop {
		s.stop = true
		s.c.Signal()
	}
	s.l.Unlock()
}

func (s *Sessions) Wait() {
	s.l.Lock()
	for !s.stop {
		s.c.Wait()
	}
	s.l.Unlock()
}

func (s *Sessions) Stop() {
	s.removeAll(true)
}
