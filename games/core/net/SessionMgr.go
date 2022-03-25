package net

import (
	"games/core/transmit"
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

//
type SessionMgr interface {
	Add(conn interface{}) Session
	Remove(peer Session)
	Get(sesID int64) Session
	Count() int64
	Stop()
	Wait()
}

//
var gSessMgr = newSessionMgr()

//
type defaultSessionMgr struct {
	peers map[int64]Session
	l     *sync.Mutex
	c     *sync.Cond
	exit  bool
}

//
func newSessionMgr() SessionMgr {
	s := &defaultSessionMgr{l: &sync.Mutex{}, peers: map[int64]Session{}}
	s.c = sync.NewCond(s.l)
	s.exit = false
	return s

}

//
func (s *defaultSessionMgr) Add(conn interface{}) Session {
	if !s.exit {
		if c, ok := conn.(net.Conn); ok {
			peer := newTCPConnection(c, transmit.NewTCPChannel())
			s.l.Lock()
			s.peers[peer.ID()] = peer
			s.l.Unlock()
			return peer
		} else if c, ok := conn.(*websocket.Conn); ok {
			peer := newTCPConnection(c, transmit.NewWSChannel())
			s.l.Lock()
			s.peers[peer.ID()] = peer
			s.l.Unlock()
			return peer
		}
	}
	return nil
}

//
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

//
func (s *defaultSessionMgr) Get(sesID int64) Session {
	s.l.Lock()
	if peer, ok := s.peers[sesID]; ok {
		s.l.Unlock()
		return peer
	}
	s.l.Unlock()
	return nil
}

//Count 有效连接数
func (s *defaultSessionMgr) Count() int64 {
	return 0
}

//
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
