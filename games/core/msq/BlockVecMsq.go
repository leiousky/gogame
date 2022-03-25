package msq

import (
	"sync"
	"sync/atomic"
)

/// <summary>
/// BlockVecMsq 切片类型
/// <summary>
type BlockVecMsq struct {
	msq         []interface{}
	l           *sync.Mutex
	c           *sync.Cond
	n           int64
	nonblocking bool
}

func NewBlockVecMsq() MsgQueue {
	s := &BlockVecMsq{l: &sync.Mutex{}}
	s.c = sync.NewCond(s.l)
	return s
}

func (s *BlockVecMsq) EnableNonBlocking(bv bool) {
	if s.nonblocking != bv {
		s.nonblocking = bv
		// str := " FALSE"
		// if s.nonblocking {
		// 	str = " TRUE"
		// }
		// log.Println("NonBlocking: ", str)
	}
}

func (s *BlockVecMsq) Push(msg interface{}) {
	{
		s.l.Lock()
		s.msq = append(s.msq, msg)
		s.l.Unlock()
		atomic.AddInt64(&s.n, 1)
	}
	s.c.Signal()
}

func (s *BlockVecMsq) Pop() (msg interface{}, exit bool) {
	{
		s.l.Lock()
		if !s.nonblocking && len(s.msq) == 0 {
			s.c.Wait()
		}
		s.l.Unlock()
	}
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

func (s *BlockVecMsq) Pick() (msgs []interface{}, exit bool) {
	{
		s.l.Lock()
		if !s.nonblocking && len(s.msq) == 0 {
			s.c.Wait()
		}
		s.l.Unlock()
	}
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

func (s *BlockVecMsq) Count() int64 {
	return atomic.LoadInt64(&s.n)
}

func (s *BlockVecMsq) Signal() {
	if !s.nonblocking {
		s.c.Signal()
	}
}

func (s *BlockVecMsq) reset() {
	s.msq = s.msq[0:0]
	atomic.StoreInt64(&s.n, 0)
}

func (s *BlockVecMsq) Close() {

}
