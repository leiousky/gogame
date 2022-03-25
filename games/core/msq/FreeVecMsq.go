package msq

import (
	"sync"
	"sync/atomic"
)

/// <summary>
/// FreeVecMsq 非阻塞切片类型
/// <summary>
type FreeVecMsq struct {
	msq []interface{}
	l   *sync.Mutex
	n   int64
}

func NewFreeVecMsq() MsgQueue {
	s := &FreeVecMsq{l: &sync.Mutex{}}
	return s
}

func (s *FreeVecMsq) EnableNonBlocking(bv bool) {

}

func (s *FreeVecMsq) Push(msg interface{}) {
	{
		s.l.Lock()
		s.msq = append(s.msq, msg)
		s.l.Unlock()
		atomic.AddInt64(&s.n, 1)
	}
}

func (s *FreeVecMsq) Pop() (msg interface{}, exit bool) {
	{
		s.l.Lock()
		if len(s.msq) > 0 {
			msg = s.msq[0]
			if msg == nil {
				exit = true
				s.reset()
			} else {
				s.msq = s.msq[1:]
				atomic.AddInt64(&s.n, -1)
			}
		}
		s.l.Unlock()
	}
	return
}

func (s *FreeVecMsq) Pick() (msgs []interface{}, exit bool) {
	{
		s.l.Lock()
		for _, msg := range s.msq {
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

func (s *FreeVecMsq) Count() int64 {
	return atomic.LoadInt64(&s.n)
}

func (s *FreeVecMsq) Signal() {
}

func (s *FreeVecMsq) reset() {
	s.msq = s.msq[0:0]
	atomic.StoreInt64(&s.n, 0)
}

func (s *FreeVecMsq) Close() {

}
