package msq

import (
	"errors"
	"fmt"
	"games/comm/utils"
	"sync/atomic"
	"time"
)

/// <summary>
/// FreeChanMsq 非阻塞chan类型
/// <summary>
type FreeChanMsq struct {
	msq chan interface{}
	n   int64
	tid uint32
}

func NewFreeChanMsq() MsgQueue {
	return &FreeChanMsq{
		msq: make(chan interface{}, 9000),
		tid: utils.GoroutineID(),
	}
}

func (s *FreeChanMsq) Push(data interface{}) error {
	if len(s.msq) == cap(s.msq) {
		return errors.New(fmt.Sprintf("pid[%v]FreeChanMsq is full", s.tid))
	}
	select {
	case s.msq <- data:
		atomic.AddInt64(&s.n, 1)
	}
	return nil
}

func (s *FreeChanMsq) Pop() (data interface{}, exit bool) {
	select {
	case q := <-s.msq:
		{
			if q == nil {
				close(s.msq)
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

func (s *FreeChanMsq) Close() {
	if s.msq != nil {
		close(s.msq)
	}
}
