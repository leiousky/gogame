package core

import (
	"errors"
	"fmt"
	"mygame/comm/utils"
	timer "mygame/core/timerv2"
	"runtime"
	"sync"
	"time"
)

/// <summary>
// 消息处理器接口
/// <summary>
type IProc interface {
	/// 获取线程内置局部定时器，肯定线程安全
	GetTimer() interface{}
	/// 获取业务句柄
	GetWorker() IWorker
	/// 添加任务
	AddTask(data interface{}) error
	/// 执行空闲回调
	Exec(cb func())
	/// 添加空闲回调
	Append(cb func())
	/// 进入消息循环开始轮询
	Run()
	/// 退出处理
	Quit()
}

/// <summary>
/// 消息处理器
/// <summary>
type Proc struct {
	tid    uint32                   //协程ID
	worker IWorker                  //业务处理
	args   []interface{}            //业务参数
	msq    chan interface{}         //消息队列
	cbs    []func()                 //空闲回调
	lock   *sync.RWMutex            //cbs锁
	exit   chan *struct{}           //退出处理，多了chan
	timer  *timer.SafeTimerScheduel // 协程安全定时器
	//exit *struct {}   	//退出处理，省掉chan
}

/// 创建消息处理器
func newMsgProc(creator IWorkerCreator, args ...interface{}) IProc {
	s := &Proc{
		tid:  utils.GoroutineID(), // newMsgProc()执行必须在Run()的go协程中调用，不然tid获取不对
		msq:  make(chan interface{}, 500),
		exit: make(chan *struct{}),
		//exit: &struct {}{},
		lock:  &sync.RWMutex{},
		timer: timer.NewSafeTimerScheduel()}
	s.worker = creator.Create(s)
	//worker初始化参数
	s.args = append(s.args, args...)
	return s
}

/// 获取线程内置局部定时器，肯定线程安全
func (s *Proc) GetTimer() interface{} {
	return s.timer
}

/// 获取业务句柄
func (s *Proc) GetWorker() IWorker {
	return s.worker
}

/// 添加任务
func (s *Proc) AddTask(data interface{}) error {
	if len(s.msq) == cap(s.msq) {
		return errors.New(fmt.Sprintf("pid=%v msq full", s.tid))
	}
	select {
	case <-s.exit:
		return errors.New(fmt.Sprintf("pid=%v Proc exit", s.tid))
	case s.msq <- data:
	}
	return nil
}

/// 线程安全检查
func (s *Proc) inThread() bool {
	return utils.GoroutineID() == s.tid
}

/// 线程安全检查
func (s *Proc) assertSafe() bool {
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
	//s.msq.Signal()
}

/// 执行空闲回调
func (s *Proc) execFunc() {
	s.assertSafe()
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

/// 进入消息循环开始轮询
func (s *Proc) Run() {
	s.assertSafe()
	worker := s.worker
	worker.OnInit(s.args...)
	i, t := 0, 200 //CPU分片
	for {
		if i > t {
			i = 0
			runtime.Gosched()
		}
		i++
		select {
		//退出处理
		case <-s.exit:
			goto end
			//定时消息
		case df := <-s.timer.Do():
			utils.SafeCall(df.Call)
		//任务消息
		case data := <-s.msq:
			t1 := time.Now()
			cmd := uint32(0)
			var peer interface{}
			worker.OnMessage(cmd, data, peer)
			duration := time.Since(t1)
			if duration > time.Second {
				fmt.Printf("elapsed:v% msg:%v", duration, data)
			}
		default:
			//处理空闲回调
			utils.SafeCall(
				func() {
					s.execFunc()
				})
		}
	}
end:
	//清理定时器
	//.......
	//关闭chan通道
	close(s.exit)
	close(s.msq)
}

/// 退出处理
func (s *Proc) Quit() {
	//{
	//	//方式1 省掉chan
	//	if s.msq != nil {
	//		s.msq <- s.exit
	//	}
	//}
	{
		//方式2
		if s.exit != nil {
			s.exit <- &struct{}{}
		}
	}
}
