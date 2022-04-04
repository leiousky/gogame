package cell

import (
	"fmt"
	"games/comm/utils"
	cb "games/core/callback"
	"games/core/conn"
	"games/core/timer"
	timerv2 "games/core/timerv2"
	"log"
	"time"

	"runtime"
	"sync"
)

/// <summary>
/// IProc 单元业务处理器接口
/// <summary>
type IProc interface {
	/// 协程ID
	GetPID() uint32
	/// 修改时钟间隔
	Reset(d time.Duration)
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
	/// s.Exec(func(v ...interface{}) {
	/// }, a, b, c)
	// Exec(f func(args ...interface{}), args ...interface{})
	/// 追加空闲回调
	/// s.Append(func(v ...interface{}) {
	/// }, a, b, c)
	// Append(f func(args ...interface{}), args ...interface{})
	/// 执行空闲回调
	/// s.Exec(func(v interface{}) {
	/// }, []interface{}{a, b, c})
	Exec(f func(args interface{}), args interface{})
	/// 追加空闲回调
	/// s.Append(func(v interface{}) {
	/// }, []interface{}{a, b, c})
	Append(f func(args interface{}), args interface{})
	/// 任务轮询(定时任务/网络任务/自定义任务/空闲任务)
	Run()
	/// 退出处理
	Quit()
}

/// <summary>
/// Proc 单元业务处理器实现
/// <summary>
type Proc struct {
	tid        uint32                     //协程ID
	msq        chan interface{}           //任务队列
	l          *sync.Mutex                //入队上锁
	idle       chan bool                  //空闲任务
	lock       *sync.RWMutex              //入队上锁
	closed     [2]bool                    //关闭标识
	funcs      []cb.Functor               //空闲回调
	worker     IWorker                    //任务处理
	ticker     *time.Ticker               //滴答时钟
	trigger    <-chan time.Time           //定时触发器
	timer      timer.ScopedTimer          //内置局部定时器，线程安全
	timerWheel timer.TimerWheel           //时间轮盘
	dispatcher IProc                      //分派任务到其他IProc
	timerv2    *timerv2.SafeTimerScheduel //协程安全定时器
}

/// NewProc()执行必须在Run()的go协程中调用，不然tid获取不对
func NewProc(d time.Duration, size int, creator IWorkerCreator, args ...interface{}) IProc {
	ticker := func(d time.Duration) *time.Ticker {
		if d <= 0 {
			return nil
		}
		return time.NewTicker(d)
	}(d)
	trigger := func(ticker *time.Ticker, d time.Duration) <-chan time.Time {
		if d <= 0 {
			return nil
		}
		return ticker.C
	}(ticker, d)
	s := &Proc{
		tid:     utils.GoroutineID(),
		msq:     make(chan interface{}, 1000),
		l:       &sync.Mutex{},
		idle:    make(chan bool, 1),
		lock:    &sync.RWMutex{},
		ticker:  ticker,
		trigger: trigger,
		timerv2: timerv2.NewSafeTimerScheduel()}
	s.worker = creator.Create(s, args...)                  //线程局部worker
	s.timerWheel = timer.NewTimerWheel(s.tid, int32(size)) //指定时间轮大小
	s.timer = timer.NewScopedTimer(s.tid)                  //线程局部定时器
	return s
}

/// 协程ID
func (s *Proc) GetPID() uint32 {
	return s.tid
}

/// 修改时钟间隔
func (s *Proc) Reset(d time.Duration) {
	s.ticker.Reset(d)
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
	s.push(data)
}

func (s *Proc) push(data interface{}) {
	if len(s.msq) == cap(s.msq) {
		panic(fmt.Sprintf("[%v]Proc.push msq is full", s.tid))
	}
	s.l.Lock()
	if data == nil {
		if !s.closed[0] {
			select {
			case s.msq <- data:
				break
			default:
				break
			}
			close(s.msq)
			s.closed[0] = true
		} else {
			panic(fmt.Sprintf("[%v]Proc.push msq repeat close", s.tid))
		}
	} else {
		if !s.closed[0] {
			select {
			case s.msq <- data:
				break
			default:
				//默认case防止msq满时阻塞
				break
			}
		} else {
			panic(fmt.Sprintf("[%v]Proc.push msq is closed", s.tid))
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
/// s.Exec(func(v ...interface{}) {
/// }, a, b, c)
// func (s *Proc) Exec(f func(args ...interface{}), args ...interface{}) {
// 	if s.inThread() {
// 		f(args...)
// 	} else {
// 		s.Append(f, args...)
// 	}
// }

/// 追加空闲回调
/// s.Append(func(v ...interface{}) {
/// }, a, b, c)
// func (s *Proc) Append(f func(args ...interface{}), args ...interface{}) {
// 	s.lock.Lock()
// 	s.funcs = append(s.funcs, cb.NewFunctor(f, args...))
// 	s.signal(false)
// 	s.lock.Unlock()
// }

/// 执行空闲回调
/// s.Exec(func(v interface{}) {
/// }, []interface{}{a, b, c})
func (s *Proc) Exec(f func(args interface{}), args interface{}) {
	if s.inThread() {
		f(args)
	} else {
		s.Append(f, args)
	}
}

/// 追加空闲回调
/// s.Append(func(v interface{}) {
/// }, []interface{}{a, b, c})
func (s *Proc) Append(f func(args interface{}), args interface{}) {
	s.lock.Lock()
	s.funcs = append(s.funcs, cb.NewFunctor(f, args))
	s.signal(false)
	s.lock.Unlock()
}

func (s *Proc) signal(stop bool) {
	// if len(s.idle) == cap(s.idle) {
	// 	panic(fmt.Sprintf("[%v]Proc.signal idle is full", s.tid))
	// }
	if stop {
		if !s.closed[1] {
			// select {
			// case s.idle <- false:
			// 	break
			// default:
			// 	break
			// }
			close(s.idle)
			s.closed[1] = true
		} else {
			panic(fmt.Sprintf("[%v]Proc.signal idle repeat close", s.tid))
		}
	} else {
		if !s.closed[1] {
			select {
			case s.idle <- true:
				break
			default:
				//默认case防止连续调用多次Append()时阻塞
				break
			}
		} else {
			panic(fmt.Sprintf("[%v]Proc.signal idle is closed", s.tid))
		}
	}
}

/// 执行空闲回调
func (s *Proc) call() {
	s.AssertInThread()
	var funcs []cb.Functor
	{
		s.lock.Lock()
		if len(s.funcs) > 0 {
			funcs = s.funcs[:]
			s.funcs = s.funcs[0:0]
		}
		s.lock.Unlock()
	}
	for _, f := range funcs {
		f.Call()
	}
}

/// 任务轮询(定时任务/网络任务/自定义任务/空闲任务)
func (s *Proc) Run() {
	s.run_msq()
}

func (s *Proc) run_msq() {
	s.AssertInThread()
	utils.CheckPanic()
	worker := s.worker
	timer := s.timer
	worker.OnInit()
	i, t := 0, 200 //CPU分片
EXIT:
	for {
		if i > t {
			i = 0
			runtime.Gosched()
		}
		i++
		//log.Printf("[%v]Proc.run_msq ...", s.tid)
		select {
		//定时任务
		case _, ok := <-s.trigger:
			{
				if ok {
					//log.Printf("[%v]Proc.run_msq timer.Poll ...", s.tid)
					timer.Poll(s.tid, worker.OnTimer)
					s.test001()
				}
				break
			}
		//定时任务
		case c, ok := <-s.timerv2.Do():
			{
				if ok {
					//log.Printf("[%v]Proc.run_msq c.Call ...", s.tid)
					utils.SafeCall(c.Call)
				}
				break
			}
		//网络任务/自定义任务
		case msg, ok := <-s.msq:
			{
				if ok {
					if msg == nil {
						//panic(errors.New("msg nil"))
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
						//panic(errors.New("channel closed, msg nil"))
						break EXIT
					} else if _, ok := msg.(*Event); ok {
						//panic(errors.New("channel closed, msg exist"))
						start := time.Now()
						s.proc(msg.(*Event), worker)
						elapsed := time.Since(start)
						if elapsed > time.Second {
						}
					}
				}
				break
			}
		//空闲任务
		case _, ok := <-s.idle:
			{
				if ok {
					//log.Printf("[%v]Proc.run_msq s.call...", s.tid)
					utils.SafeCall(s.call)
				}
				break
			}
			//轮询时默认case会导致CPU负载非常高，应该禁用
			//default:
			//	break
		}
	}
	s.cleanup()
	log.Printf("[%v]Proc.run_msq exit...", s.tid)
}

func (s *Proc) cleanup() {
	s.timer.RemoveTimers()
	s.close_idle()
	s.ticker.Stop()
}

func (s *Proc) close_idle() {
	s.lock.Lock()
	s.signal(true)
	s.lock.Unlock()
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
		break
	case EVTCustom:
		ev := data.obj.(*customEvent)
		if ev.handler != nil {
			ev.handler(ev.cmd, ev.msg, ev.peer)
		} else {
			worker.OnCustom(ev.cmd, ev.msg, ev.peer)
		}
		break
	}
	if s.dispatcher != nil {
		s.dispatcher.AddTask(data)
	}
}

/// 退出处理
func (s *Proc) Quit() {
	s.push(nil)
}

func (s *Proc) test001() {
	//s.test002()
}

func (s *Proc) test002() {
	log.Println("Proc.test002 ...")
	for i := 1; i <= 10; i++ {
		s.Append(func(v interface{}) {
			args := v.([]interface{})
			log.Printf("%v %v %v ...\n", args[0].(int), args[1].(string), args[2].(float64))
		}, []interface{}{i, "hello", 3.1415926})
	}
}
