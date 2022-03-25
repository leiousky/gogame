package msq

import (
	"sync/atomic"
	"time"
)

/// <summary>
/// FreeChanMsq 非阻塞chan类型
/// <summary>
type FreeChanMsq struct {
	msq chan interface{}
	n   int64
}

func NewFreeChanMsq() MsgQueue {
	return &FreeChanMsq{msq: make(chan interface{}, 9000)}
}

func (s *FreeChanMsq) Push(msg interface{}) {
	s.msq <- msg
	atomic.AddInt64(&s.n, 1)
}

func (s *FreeChanMsq) Pop() (msg interface{}, exit bool) {
	select {
	case q := <-s.msq:
		{
			if q == nil {
				close(s.msq)
				exit = true
				break
			} else {
				msg = q
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

func (s *FreeChanMsq) Pick() (msgs []interface{}, exit bool) {
	msg, e := s.Pop()
	exit = e
	if msg != nil && !exit {
		msgs = append(msgs, msg)
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

func (s *FreeChanMsq) Close() {
	if s.msq != nil {
		close(s.msq)
	}
}
