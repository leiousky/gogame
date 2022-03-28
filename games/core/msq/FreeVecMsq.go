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

func (s *FreeVecMsq) Push(data interface{}) {
	{
		s.l.Lock()
		s.msq = append(s.msq, data)
		s.l.Unlock()
		atomic.AddInt64(&s.n, 1)
	}
}

func (s *FreeVecMsq) Pop() (data interface{}, exit bool) {
	{
		s.l.Lock()
		if len(s.msq) > 0 {
			data = s.msq[0]
			if data == nil {
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

func (s *FreeVecMsq) Pick() (v []interface{}, exit bool) {
	{
		s.l.Lock()
		for _, data := range s.msq {
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

func (s *FreeVecMsq) Count() int64 {
	return atomic.LoadInt64(&s.n)
}

func (s *FreeVecMsq) Signal() {
}

func (s *FreeVecMsq) reset() {
	s.msq = s.msq[0:0]
	atomic.StoreInt64(&s.n, 0)
}
