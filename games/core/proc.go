package core

import (
	"fmt"
	"games/comm/utils"
	cb "games/core/callback"
	"games/core/conn"
	"games/core/msq"
	"games/core/timer"
	timerv2 "games/core/timerv2"
	"log"
	"time"

	"runtime"
	"sync"
)

/// <summary>
// 消息处理器接口
/// <summary>
type IProc interface {
	/// 协程ID
	GetPID() uint32
	/// 内置局部定时器，线程安全
	GetTimer() interface{}
	GetTimerv2() interface{}
	/// 业务句柄
	GetWorker() IWorker
	/// 添加任务
	AddTask(data *Event)
	AddReadTask(cmd uint32, msg interface{}, peer conn.Session)
	AddReadTaskWith(handler cb.ReadCallback, cmd uint32, msg interface{}, peer conn.Session)
	AddCustomTask(cmd uint32, msg interface{}, peer conn.Session)
	AddCustomTaskWith(handler cb.CustomCallback, cmd uint32, msg interface{}, peer conn.Session)
	/// 定时器任务
	RunAfter(delay int32, args interface{}) uint32
	RunAfterWith(delay int32, handler timer.TimerCallback, args interface{}) uint32
	RunEvery(delay, interval int32, args interface{}) uint32
	RunEveryWith(delay, interval int32, handler timer.TimerCallback, args interface{}) uint32
	RemoveTimer(timerID uint32)
	RemoveTimers()
	/// 线程安全检查
	AssertInThread() bool
	/// 任务分发
	SetDispatcher(c IProc)
	GetDispatcher() IProc
	/// 执行空闲回调
	Exec(cb func())
	/// 添加空闲回调
	Append(cb func())
	/// 任务轮询
	Run()
	/// 退出处理
	Quit()
}

/// <summary>
/// 消息处理器
/// <summary>
type Proc struct {
	msQ        chan interface{}
	l          *sync.Mutex
	closed     bool
	tid        uint32            //协程ID
	msq        msq.MsgQueue      //任务队列
	worker     IWorker           //任务处理
	timer      timer.ScopedTimer //内置局部定时器，线程安全
	timerWheel timer.TimerWheel  //时间轮盘
	dispatcher IProc             //分派任务到其他IProc
	args       []interface{}     //任务参数
	cbs        []func()          //空闲回调
	lock       *sync.RWMutex
	timerv2    *timerv2.SafeTimerScheduel //协程安全定时器
	selectQ    int
}

const (
	TmsQ int = int(0)
	Tmsq int = int(1)
)

/// 创建消息处理器
func newMsgProc(creator IWorkerCreator, size int, args ...interface{}) IProc {
	s := &Proc{
		msQ:     make(chan interface{}, 1000),
		l:       &sync.Mutex{},
		tid:     utils.GoroutineID(), //newMsgProc()执行必须在Run()的go协程中调用，不然tid获取不对
		msq:     msq.NewFreeVecMsq(),
		lock:    &sync.RWMutex{},
		timerv2: timerv2.NewSafeTimerScheduel(),
		selectQ: TmsQ}
	s.worker = creator.Create(s)                           //线程局部worker
	s.timerWheel = timer.NewTimerWheel(s.tid, int32(size)) //指定时间轮大小
	s.timer = timer.NewScopedTimer(s.tid)                  //线程局部定时器
	//worker初始化参数
	s.args = append(s.args, args...)
	return s
}

/// 协程ID
func (s *Proc) GetPID() uint32 {
	return s.tid
}

/// 获取线程内置局部定时器，肯定线程安全
func (s *Proc) GetTimer() interface{} {
	return s.timer
}

func (s *Proc) GetTimerv2() interface{} {
	return s.timerv2
}

/// 获取业务句柄
func (s *Proc) GetWorker() IWorker {
	return s.worker
}

/// 任务分发
func (s *Proc) SetDispatcher(c IProc) {
	s.AssertInThread()
	s.dispatcher = c
}

func (s *Proc) GetDispatcher() IProc {
	s.AssertInThread()
	return s.dispatcher
}

/// 添加任务
func (s *Proc) AddTask(data *Event) {
	switch s.selectQ {
	case TmsQ:
		s.push(data)
		break
	case Tmsq:
		s.msq.Push(data)
		break
	}
}

func (s *Proc) push(data interface{}) {
	if len(s.msQ) == cap(s.msQ) {
		panic(fmt.Sprintf("pid[%v]msQ is full", s.tid))
	}
	s.l.Lock()
	if data == nil {
		if !s.closed {
			s.msQ <- data
			close(s.msQ)
			s.closed = true
		} else {
			panic(fmt.Sprintf("pid[%v]msQ repeat close", s.tid))
		}
	} else {
		if !s.closed {
			select {
			case s.msQ <- data:
			}
		} else {
			panic(fmt.Sprintf("pid[%v]msQ is closed", s.tid))
		}
	}
	s.l.Unlock()
}

func (s *Proc) AddReadTask(cmd uint32, msg interface{}, peer conn.Session) {
	s.AddTask(createEvent(EVTRead, createReadEvent(cmd, msg, peer), nil))
}

func (s *Proc) AddReadTaskWith(handler cb.ReadCallback, cmd uint32, msg interface{}, peer conn.Session) {
	s.AddTask(createEvent(EVTRead, createReadEventWith(handler, cmd, msg, peer), nil))
}

func (s *Proc) AddCustomTask(cmd uint32, msg interface{}, peer conn.Session) {
	s.AddTask(createEvent(EVTCustom, createCustomEvent(cmd, msg, peer), nil))
}

func (s *Proc) AddCustomTaskWith(handler cb.CustomCallback, cmd uint32, msg interface{}, peer conn.Session) {
	s.AddTask(createEvent(EVTCustom, createCustomEventWith(handler, cmd, msg, peer), nil))
}

func (s *Proc) RunAfter(delay int32, args interface{}) uint32 {
	return s.timer.CreateTimer(delay, 0, args)
}

func (s *Proc) RunAfterWith(delay int32, handler timer.TimerCallback, args interface{}) uint32 {
	return s.timer.CreateTimerWithCB(delay, 0, handler, args)
}

func (s *Proc) RunEvery(delay, interval int32, args interface{}) uint32 {
	return s.timer.CreateTimer(delay, interval, args)
}

func (s *Proc) RunEveryWith(delay, interval int32, handler timer.TimerCallback, args interface{}) uint32 {
	return s.timer.CreateTimerWithCB(delay, interval, handler, args)
}

func (s *Proc) RemoveTimer(timerID uint32) {
	s.timer.RemoveTimer(timerID)
}

func (s *Proc) RemoveTimers() {
	s.timer.RemoveTimers()
}

/// 线程安全检查
func (s *Proc) inThread() bool {
	return utils.GoroutineID() == s.tid
}

/// 线程安全检查
func (s *Proc) AssertInThread() bool {
	if !s.inThread() {
		panic(fmt.Sprintf("非线程安全 %v", s.tid))
	}
	return true
}

/// 执行空闲回调
func (s *Proc) Exec(cb func()) {
	if s.inThread() {
		cb()
	} else {
		s.Append(cb)
	}
}

/// 添加空闲回调
func (s *Proc) Append(cb func()) {
	s.lock.Lock()
	s.cbs = append(s.cbs, cb)
	s.lock.Unlock()
	s.msq.Signal()
}

/// 执行空闲回调
func (s *Proc) execFunc() {
	s.AssertInThread()
	var cbs []func()
	{
		s.lock.Lock()
		if len(s.cbs) > 0 {
			cbs = s.cbs[:]
			s.cbs = s.cbs[0:0]
		}
		s.lock.Unlock()
	}
	for _, cb := range cbs {
		cb()
	}
}

/// 任务轮询
func (s *Proc) Run() {
	switch s.selectQ {
	case TmsQ:
		s.run_msQ()
		break
	case Tmsq:
		s.run_msq()
		break
	}
}

func (s *Proc) run_msQ() {
	s.AssertInThread()
	utils.CheckPanic()
	worker := s.worker
	timer := s.timer
	worker.OnInit(s.args...)
	s.msq.EnableNonBlocking(true)
	i, t := 0, 200 //CPU分片
EXIT:
	for {
		if i > t {
			i = 0
			runtime.Gosched()
		}
		i++
		select {
		//定时消息
		case fn := <-s.timerv2.Do():
			utils.SafeCall(fn.Call)
			timer.Poll(s.tid, worker.OnTimer)
			break
		//任务消息
		case msg, ok := <-s.msQ:
			if ok {
				if msg == nil {
					//panic(fmt.Sprintf("msg nil"))
					break EXIT
				} else if _, ok := msg.(*Event); ok {
					start := time.Now()
					s.proc(msg.(*Event), worker)
					elapsed := time.Since(start)
					if elapsed > time.Second {
					}
				}
			} else {
				//channel closed
				if msg == nil {
					//panic(fmt.Sprintf("channel closed, msg nil"))
					break EXIT
				} else if _, ok := msg.(*Event); ok {
					//panic(fmt.Sprintf("channel closed, msg exist"))
					start := time.Now()
					s.proc(msg.(*Event), worker)
					elapsed := time.Since(start)
					if elapsed > time.Second {
					}
				}
			}
			break
		default:
			//log.Println("execFunc...")
			//处理空闲回调
			utils.SafeCall(
				func() {
					s.execFunc()
				})
		}
	}
	timer.RemoveTimers()
	log.Printf("proc run_msQ tid=%v exit...", s.tid)
}

func (s *Proc) run_msq() {
	s.AssertInThread()
	utils.CheckPanic()
	worker := s.worker
	timer := s.timer
	worker.OnInit(s.args...)
	s.msq.EnableNonBlocking(true)
	flag := 0
	exit := false
	i, t := 0, 200 //CPU分片
EXIT:
	for {
		if i > t {
			i = 0
			runtime.Gosched()
		}
		i++
		//定时器轮询
		//log.Printf("--- *** ----------------------------- [%05d]Run Poll begin...\n", s.tid)
		timer.Poll(s.tid, worker.OnTimer)
		//log.Printf("--- *** ----------------------------- [%05d]Run Poll end...\n", s.tid)
		switch flag {
		case 0:
			{
				//单条消息处理
				msg, b := s.msq.Pop()
				exit = b
				if msg != nil && !exit {
					if _, ok := msg.(*Event); ok {
						//log.Printf("--- *** ----------------------------- [%05d]Run proc begin...\n", s.pid)
						s.proc(msg.(*Event), worker)
						//log.Printf("--- *** ----------------------------- [%05d]Run proc end...\n", s.pid)
					}
				}
				if nil == msg && !exit {
					//log.Printf("--- *** ----------------------------- [%05d]Run time.Sleep...\n", s.pid)
					//time.Sleep(50 * time.Millisecond)
					time.Sleep(0)
				}
				break
			}
		case 1:
			{
				//批量消息处理
				msgs, b := s.msq.Pick()
				exit = b
				for _, msg := range msgs {
					if _, ok := msg.(*Event); ok {
						//log.Printf("--- *** ----------------------------- [%05d]Run proc begin...\n", s.pid)
						s.proc(msg.(*Event), worker)
						//log.Printf("--- *** ----------------------------- [%05d]Run proc end...\n", s.pid)
					}
				}
				if 0 == len(msgs) && !exit {
					//log.Printf("--- *** ----------------------------- [%05d]Run time.Sleep...\n", s.pid)
					//time.Sleep(50 * time.Millisecond)
					time.Sleep(0)
				}
				break
			}
		}
		log.Println("execFunc...")
		//处理空闲回调
		utils.SafeCall(
			func() {
				s.execFunc()
			})
		if exit {
			break EXIT
		}
	}
	timer.RemoveTimers()
	log.Printf("proc run_msq tid=%v exit...", s.tid)
}

/// 处理任务队列
func (s *Proc) proc(data *Event, worker IWorker) {
	s.dispatcher = nil
	switch data.ev {
	case EVTRead:
		ev := data.obj.(*readEvent)
		if ev.handler != nil {
			ev.handler(ev.cmd, ev.msg, ev.peer)
		} else {
			worker.OnRead(ev.cmd, ev.msg, ev.peer)
		}
	case EVTCustom:
		ev := data.obj.(*customEvent)
		if ev.handler != nil {
			ev.handler(ev.cmd, ev.msg, ev.peer)
		} else {
			worker.OnCustom(ev.cmd, ev.msg, ev.peer)
		}
	}
	if s.dispatcher != nil {
		s.dispatcher.AddTask(data)
	}
}

/// 退出处理
func (s *Proc) Quit() {
	switch s.selectQ {
	case TmsQ:
		s.push(nil)
		break
	case Tmsq:
		s.msq.Push(nil)
		break
	}
}
