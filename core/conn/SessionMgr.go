package conn

import (
	"log"
	"sync"
)

/// <summary>
/// ISessions 连接会话容器
/// <summary>
type ISessions interface {
	Get(id int64) Session
	Count() int
	Add(peer Session) bool
	Remove(peer Session)
	CloseAll()
	Wait()
	Stop()
}

/// <summary>
/// Sessions 连接会话容器
/// <summary>
type Sessions struct {
	peers map[int64]Session
	n     int
	l     *sync.Mutex
	c     *sync.Cond
	stop  bool
	done  bool
}

func NewSessions() ISessions {
	s := &Sessions{l: &sync.Mutex{}, peers: map[int64]Session{}}
	s.c = sync.NewCond(s.l)
	return s
}

func (s *Sessions) Get(id int64) Session {
	s.l.Lock()
	if peer, ok := s.peers[id]; ok {
		s.l.Unlock()
		return peer
	}
	s.l.Unlock()
	return nil
}

func (s *Sessions) Count() int {
	c := 0
	s.l.Lock()
	//c = len(s.peers)
	c = s.n
	s.l.Unlock()
	return c
}

func (s *Sessions) Add(peer Session) bool {
	ok := false
	s.l.Lock()
	if !s.stop {
		log.Printf("Sessions.Add => %v", peer.Name())
		s.peers[peer.ID()] = peer
		s.n++
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
		s.n--
	}
	// if s.stop && len(s.peers) == 0 {
	if s.stop && s.n == 0 {
		s.done = true
		s.c.Signal()
	}
	s.l.Unlock()
}

/// s.closeAll -> peer.Close -> s.Remove
func (s *Sessions) closeAll(stop bool) {
	s.l.Lock()
	if stop {
		s.stop = true
	}
	for _, peer := range s.peers {
		peer.Close()
	}
	s.l.Unlock()
}

func (s *Sessions) CloseAll() {
	s.closeAll(false)
}

func (s *Sessions) Wait() {
	s.l.Lock()
	for !s.done {
		s.c.Wait()
	}
	s.l.Unlock()
}

func (s *Sessions) Stop() {
	s.closeAll(true)
}
