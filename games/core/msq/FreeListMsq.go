package msq

import (
	"container/list"
	"sync"
	"sync/atomic"
)

/// <summary>
/// FreeListMsq 非阻塞链表类型
/// <summary>
type FreeListMsq struct {
	msq *list.List
	l   *sync.Mutex
	n   int64
}

func NewFreeListMsq() MsgQueue {
	s := &FreeListMsq{msq: list.New(),
		l: &sync.Mutex{}}
	return s
}

func (s *FreeListMsq) EnableNonBlocking(bv bool) {

}

func (s *FreeListMsq) Push(msg interface{}) {
	{
		s.l.Lock()
		s.msq.PushBack(msg)
		s.l.Unlock()
		atomic.AddInt64(&s.n, 1)
	}
}

func (s *FreeListMsq) Pop() (msg interface{}, exit bool) {
	{
		s.l.Lock()
		if elem := s.msq.Front(); elem != nil {
			msg = elem.Value
			if msg == nil {
				exit = true
				s.reset()
			} else {
				atomic.AddInt64(&s.n, -1)
				s.msq.Remove(elem)
			}
		}
		s.l.Unlock()
	}
	return
}

func (s *FreeListMsq) Pick() (msgs []interface{}, exit bool) {
	{
		s.l.Lock()
		var next *list.Element
		for elem := s.msq.Front(); elem != nil; elem = next {
			next = elem.Next()
			msg := elem.Value
			s.msq.Remove(elem)
			if msg == nil {
				exit = true
				break
			} else {
				msgs = append(msgs, msg)
			}
		}
		s.reset()
		s.l.Unlock()
	}
	return
}

func (s *FreeListMsq) Count() int64 {
	return atomic.LoadInt64(&s.n)
}

func (s *FreeListMsq) Signal() {
}

func (s *FreeListMsq) reset() {
	var next *list.Element
	for elem := s.msq.Front(); elem != nil; elem = next {
		next = elem.Next()
		s.msq.Remove(elem)
	}
	atomic.StoreInt64(&s.n, 0)
}

func (s *FreeListMsq) Close() {

}
