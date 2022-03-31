package msq

import (
	"fmt"
	"games/comm/utils"
	"sync"
	"sync/atomic"
	"time"
)

/// <summary>
/// FreeChanMsq 非阻塞chan类型
/// <summary>
type FreeChanMsq struct {
	msq    chan interface{}
	l      *sync.Mutex
	closed bool
	n      int64
	tid    uint32
}

func NewFreeChanMsq() MsgQueue {
	return &FreeChanMsq{
		msq: make(chan interface{}, 9000),
		l:   &sync.Mutex{},
		tid: utils.GoroutineID(),
	}
}

func (s *FreeChanMsq) Push(data interface{}) {
	if len(s.msq) == cap(s.msq) {
		panic(fmt.Sprintf("pid[%v]FreeChanMsq is full", s.tid))
	}
	s.l.Lock()
	if data == nil {
		if !s.closed {
			s.msq <- data
			close(s.msq)
			s.closed = true
		} else {
			panic(fmt.Sprintf("pid[%v]FreeChanMsq repeat close", s.tid))
		}
	} else {
		if !s.closed {
			select {
			case s.msq <- data:
				atomic.AddInt64(&s.n, 1)
			}
		} else {
			panic(fmt.Sprintf("pid[%v]FreeChanMsq is closed", s.tid))
		}
	}
	s.l.Unlock()
}

func (s *FreeChanMsq) Pop() (data interface{}, exit bool) {
	select {
	case q, ok := <-s.msq:
		{
			if q == nil || !ok {
				//close(s.msq)
				exit = true
				break
			} else {
				data = q
			}
			atomic.AddInt64(&s.n, -1)
		}
	//case <-time.After(time.Nanosecond):
	//case <-time.After(time.Microsecond):
	case <-time.After(time.Millisecond):
		//log.Printf("--- *** ----------------------------- [%05d]Run time.After...\n", os.Getpid())
		//default:
	}
	return
}

func (s *FreeChanMsq) Pick() (v []interface{}, exit bool) {
	data, e := s.Pop()
	exit = e
	if data != nil && !exit {
		v = append(v, data)
	}
	return
}

func (s *FreeChanMsq) Count() int64 {
	return atomic.LoadInt64(&s.n)
}

func (s *FreeChanMsq) Signal() {
}

func (s *FreeChanMsq) EnableNonBlocking(bv bool) {
}
