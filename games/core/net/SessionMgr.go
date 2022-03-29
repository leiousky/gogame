package net

import (
	"games/core/net/transmit"
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

type SessionMgr interface {
	Add(name string, conn interface{}, connType SesType) Session
	Remove(peer Session)
	Get(sesID int64) Session
	Count() int
	Stop()
	Wait()
}

var gSessMgr = newSessionMgr()

type defaultSessionMgr struct {
	peers map[int64]Session
	l     *sync.Mutex
	c     *sync.Cond
	exit  bool
}

func newSessionMgr() SessionMgr {
	s := &defaultSessionMgr{l: &sync.Mutex{}, peers: map[int64]Session{}}
	s.c = sync.NewCond(s.l)
	s.exit = false
	return s

}

func (s *defaultSessionMgr) Add(name string, conn interface{}, connType SesType) Session {
	if !s.exit {
		if c, ok := conn.(net.Conn); ok {
			peer := newTCPConnection(name, c, connType, transmit.NewTCPChannel())
			s.l.Lock()
			s.peers[peer.ID()] = peer
			s.l.Unlock()
			return peer
		} else if c, ok := conn.(*websocket.Conn); ok {
			peer := newTCPConnection(name, c, connType, transmit.NewWSChannel())
			s.l.Lock()
			s.peers[peer.ID()] = peer
			s.l.Unlock()
			return peer
		}
	}
	return nil
}

func (s *defaultSessionMgr) Remove(peer Session) {
	s.l.Lock()
	if _, ok := s.peers[peer.ID()]; ok {
		delete(s.peers, peer.ID())
	}
	if s.exit && len(s.peers) == 0 {
		s.c.Signal()
	}
	s.l.Unlock()
}

func (s *defaultSessionMgr) Get(sesID int64) Session {
	s.l.Lock()
	if peer, ok := s.peers[sesID]; ok {
		s.l.Unlock()
		return peer
	}
	s.l.Unlock()
	return nil
}

func (s *defaultSessionMgr) Count() int {
	c := 0
	s.l.Lock()
	c = len(s.peers)
	s.l.Unlock()
	return c
}

func (s *defaultSessionMgr) Stop() {
	s.l.Lock()
	s.exit = true
	for _, peer := range s.peers {
		peer.Close()
	}
	s.l.Unlock()
}

func (s *defaultSessionMgr) Wait() {
	s.l.Lock()
	s.c.Wait()
	//log.Printf("SessionMgr::Wait exit...")
	s.l.Unlock()
}
