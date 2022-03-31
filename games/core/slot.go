package core

import (
	"sync"
	"sync/atomic"
	"time"
)

/// <summary>
/// ISlot 消息处理单元接口
/// <summary>
type ISlot interface {
	/// 定时器间隔时间
	Duration() time.Duration
	/// 添加worker初始化参数
	Add(args ...interface{})
	/// 启动协程并返回消息处理器句柄
	Sched() IProc
	/// 获取消息处理器句柄
	GetProc() IProc
	/// 退出处理
	Stop()
}

const (
	Idle int32 = iota
	Running
)

/// <summary>
/// Slot 消息处理单元
/// <summary>
type Slot struct {
	proc    IProc //消息处理器
	args    []interface{}
	lock    *sync.Mutex
	cond    *sync.Cond
	creator IWorkerCreator
	d       time.Duration
	size    int
	sta     int32
}

/// 创建消息处理单元
func NewMsgSlot(d time.Duration, size int, creator IWorkerCreator) ISlot {
	s := &Slot{d: d, size: size, creator: creator, lock: &sync.Mutex{}}
	s.cond = sync.NewCond(s.lock)
	return s
}

/// 定时器间隔时间
func (s *Slot) Duration() time.Duration {
	return s.d
}

/// 添加worker初始化参数
func (s *Slot) Add(args ...interface{}) {
	s.args = append(s.args, args...)
}

/// 启动协程并返回消息处理器句柄
func (s *Slot) Sched() IProc {
	if atomic.CompareAndSwapInt32(&s.sta, Idle, Running) {
		go s.run()
	}
	{
		s.lock.Lock()
		for s.proc == nil {
			s.cond.Wait()
		}
		s.lock.Unlock()
	}
	return s.proc
}

/// 获取消息处理器句柄
func (s *Slot) GetProc() IProc {
	return s.proc
}

/// 执行协程处理任务
func (s *Slot) run() {
	proc := newMsgProc(s.d, s.size, s.creator, s.args...)
	s.lock.Lock()
	s.proc = proc
	s.cond.Signal()
	s.lock.Unlock()
	s.proc.Run()
	atomic.StoreInt32(&s.sta, Idle)
	s.proc = nil
}

/// 退出处理
func (s *Slot) Stop() {
	if s.proc != nil && atomic.CompareAndSwapInt32(&s.sta, Running, Idle) {
		s.proc.Quit()
	}
}
