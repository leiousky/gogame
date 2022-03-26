package msq

import (
	"container/list"
	"sync"
	"sync/atomic"
)

/// <summary>
/// BlockListMsq 链表类型
/// <summary>
type BlockListMsq struct {
	msq         *list.List
	l           *sync.Mutex
	c           *sync.Cond
	n           int64
	nonblocking bool
}

func NewBlockListMsq() MsgQueue {
	s := &BlockListMsq{msq: list.New(),
		l: &sync.Mutex{}}
	s.c = sync.NewCond(s.l)
	return s
}

func (s *BlockListMsq) EnableNonBlocking(bv bool) {
	if s.nonblocking != bv {
		s.nonblocking = bv
		// str := " FALSE"
		// if s.nonblocking {
		// 	str = " TRUE"
		// }
		// log.Println("NonBlocking: ", str)
	}
}

func (s *BlockListMsq) Push(data interface{}) error {
	{
		s.l.Lock()
		s.msq.PushBack(data)
		s.l.Unlock()
	}
	atomic.AddInt64(&s.n, 1)
	s.c.Signal()
	return nil
}

func (s *BlockListMsq) Pop() (data interface{}, exit bool) {
	{
		s.l.Lock()
		if !s.nonblocking && s.msq.Len() == 0 {
			s.c.Wait()
		}
		s.l.Unlock()
	}
	{
		s.l.Lock()
		if elem := s.msq.Front(); elem != nil {
			data = elem.Value
			if data == nil {
				exit = true
				s.reset()
			} else {
				s.msq.Remove(elem)
			}
			atomic.AddInt64(&s.n, -1)
		}
		s.l.Unlock()
	}
	return
}

func (s *BlockListMsq) Pick() (v []interface{}, exit bool) {
	{
		s.l.Lock()
		if !s.nonblocking && s.msq.Len() == 0 {
			s.c.Wait()
		}
		s.l.Unlock()
	}
	{
		s.l.Lock()
		var next *list.Element
		for elem := s.msq.Front(); elem != nil; elem = next {
			next = elem.Next()
			data := elem.Value
			s.msq.Remove(elem)
			if data == nil {
				exit = true
				break
			} else {
				v = append(v, data)
			}
		}

		s.reset()
		s.l.Unlock()
	}
	return
}

func (s *BlockListMsq) Count() int64 {
	return atomic.LoadInt64(&s.n)
}

func (s *BlockListMsq) Signal() {
	if !s.nonblocking {
		s.c.Signal()
	}
}

func (s *BlockListMsq) reset() {
	var next *list.Element
	for elem := s.msq.Front(); elem != nil; elem = next {
		next = elem.Next()
		s.msq.Remove(elem)
	}
	atomic.StoreInt64(&s.n, 0)
}

func (s *BlockListMsq) Close() {

}
