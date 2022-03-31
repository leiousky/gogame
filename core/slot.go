package core

import (
	"sync"
	"sync/atomic"
	"time"
)

/// <summary>
/// ISlot 业务邮槽接口(Proc启动器)
/// <summary>
type ISlot interface {
	/// 滴答时钟间隔
	Duration() time.Duration
	/// 修改时钟间隔
	Reset(d time.Duration)
	/// 添加worker初始化参数
	Add(args ...interface{})
	/// 启动协程并返回Proc句柄
	Sched() IProc
	/// 获取Proc句柄
	GetProc() IProc
	/// 退出处理
	Stop()
}

const (
	Idle int32 = iota
	Running
)

/// <summary>
/// Slot 业务邮槽实现(Proc启动器)
/// <summary>
type Slot struct {
	proc    IProc
	args    []interface{}
	lock    *sync.Mutex
	cond    *sync.Cond
	creator IWorkerCreator
	d       time.Duration
	size    int
	sta     int32
}

/// 创建业务邮槽
func NewMsgSlot(d time.Duration, size int, creator IWorkerCreator) ISlot {
	s := &Slot{d: d, size: size, creator: creator, lock: &sync.Mutex{}}
	s.cond = sync.NewCond(s.lock)
	return s
}

/// 滴答时钟间隔
func (s *Slot) Duration() time.Duration {
	return s.d
}

/// 修改时钟间隔
func (s *Slot) Reset(d time.Duration) {
	s.proc.Reset(d)
}

/// 添加worker初始化参数
func (s *Slot) Add(args ...interface{}) {
	s.args = append(s.args, args...)
}

/// 启动协程并返回Proc句柄
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

/// 获取Proc句柄
func (s *Slot) GetProc() IProc {
	return s.proc
}

/// 执行协程处理任务
func (s *Slot) run() {
	proc := newProc(s.d, s.size, s.creator, s.args...)
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
