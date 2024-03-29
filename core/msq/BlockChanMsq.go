package msq

import (
	"fmt"
	"games/comm/utils"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

//
const (
	KBlocking int32 = iota
	KNonblocking
)

/// <summary>
/// BlockChanMsq chan类型
/// <summary>
type BlockChanMsq struct {
	msq         chan interface{}
	l           *sync.Mutex
	closed      bool
	signal      chan bool
	n           int64
	tid         uint32
	nonblocking int32
}

func NewBlockChanMsq() MsgQueue {
	return &BlockChanMsq{
		msq:    make(chan interface{}, 100),
		l:      &sync.Mutex{},
		signal: make(chan bool, 1),
		tid:    utils.GoroutineID()}
}

func (s *BlockChanMsq) Push(data interface{}) {
	if len(s.msq) == cap(s.msq) {
		panic(fmt.Sprintf("pid[%v]BlockChanMsq is full", s.tid))
	}
	s.l.Lock()
	if data == nil {
		if !s.closed {
			s.msq <- data
			close(s.msq)
			s.closed = true
		} else {
			panic(fmt.Sprintf("pid[%v]BlockChanMsq repeat close", s.tid))
		}
	} else {
		if !s.closed {
			select {
			case s.msq <- data:
				atomic.AddInt64(&s.n, 1)
			}
		} else {
			panic(fmt.Sprintf("pid[%v]BlockChanMsq is closed", s.tid))
		}
	}
	s.l.Unlock()
}

func (s *BlockChanMsq) blockPop() (data interface{}, exit bool) {
	select {
	case q, ok := <-s.msq:
		{
			if q == nil || !ok {
				//close(s.msq)
				close(s.signal)
				exit = true
				break
			} else {
				data = q
			}
			atomic.AddInt64(&s.n, -1)
		}
	case <-s.signal:
	}
	return
}

func (s *BlockChanMsq) nonblockPop() (data interface{}, exit bool) {
	select {
	case q, ok := <-s.msq:
		{
			if q == nil || !ok {
				//close(s.msq)
				close(s.signal)
				exit = true
				break
			} else {
				data = q
			}
			atomic.AddInt64(&s.n, -1)
		}
	case <-time.After(time.Microsecond):
	default:
	}
	return
}

func (s *BlockChanMsq) Pop() (data interface{}, exit bool) {
	if KBlocking == atomic.LoadInt32(&s.nonblocking) {
		data, exit = s.blockPop()
	} else {
		data, exit = s.nonblockPop()
	}
	return
}

func (s *BlockChanMsq) Pick() (v []interface{}, exit bool) {
	data, e := s.Pop()
	exit = e
	if data != nil && !exit {
		v = append(v, data)
	}
	return
}

func (s *BlockChanMsq) Count() int64 {
	return atomic.LoadInt64(&s.n)
}

func (s *BlockChanMsq) Signal() {
	if KBlocking == atomic.LoadInt32(&s.nonblocking) {
		s.signalx()
	}
}

func (s *BlockChanMsq) EnableNonBlocking(bv bool) {
	val := KBlocking
	if bv == true {
		val = KNonblocking
	}
	old := atomic.LoadInt32(&s.nonblocking)
	if old != val {
		atomic.StoreInt32(&s.nonblocking, val)
	}
	//阻塞变为非阻塞
	if old == KBlocking && val == KNonblocking {
		s.signalx()
	}
}

func (s *BlockChanMsq) signalx() {
	s.signal <- true
}

var msq = NewBlockChanMsq()

func onInput1(str string) int {
	switch str {
	case "q":
		{
			msq.Push(nil)
			return -1
		}
	case "s":
		{
			msq.Signal()

		}
	case "b":
		{
			msq.EnableNonBlocking(false)
		}
	case "nb":
		{
			msq.EnableNonBlocking(true)
		}
	default:
		{

		}
	}
	return 0
}

func testChanMsq() {
	go func(msq MsgQueue) {
		for {
			msg, exit := msq.Pop()
			if exit == true {
				break
			}
			if msg != nil && !exit {
				log.Println("testChanMsq receive: ", msg)
			} else {
				log.Println("continue ...")
			}
		}
		log.Println("exit...")
	}(msq)
	for {
		utils.ReadConsole(onInput1)
	}
}
