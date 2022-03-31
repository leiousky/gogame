package core

import (
	"games/comm/utils"
	"os"
	"os/signal"
	"time"
)

/// <summary>
/// IMailbox 邮槽管理器接口
/// <summary>
type IMailbox interface {
	/// 添加若干邮槽
	Add(d time.Duration, size int, creator IWorkerCreator, num int)
	/// 添加一个邮槽
	AddOne(d time.Duration, size int, creator IWorkerCreator) ISlot
	/// 遍历每个邮槽
	Range(cb func(ISlot, int))
	/// 启动所有邮槽协程处理
	Start()
	/// 获取下一个邮槽
	GetNextSlot() ISlot
	/// 等待退出
	Wait()
	/// 主动退出
	Stop()
}

/// <summary>
/// Mailbox cell处理池子
/// <summary>
type Mailbox struct {
	slots []ISlot
	next  int
	ch    chan os.Signal
	done  chan os.Signal
}

/// 创建邮槽管理器
func NewMailBox() IMailbox {
	return &Mailbox{next: 0}
}

/// 添加若干邮槽
func (s *Mailbox) Add(d time.Duration, size int, creator IWorkerCreator, num int) {
	for i := 0; i < num; i++ {
		slot := NewMsgSlot(d, size, creator)
		s.slots = append(s.slots, slot)
	}
}

/// 添加一个邮槽
func (s *Mailbox) AddOne(d time.Duration, size int, creator IWorkerCreator) ISlot {
	slot := NewMsgSlot(d, size, creator)
	s.slots = append(s.slots, slot)
	return slot
}

/// 遍历每个邮槽
func (s *Mailbox) Range(cb func(ISlot, int)) {
	utils.SafeCall(func() {
		for i, slot := range s.slots {
			cb(slot, i)
		}
	})
}

/// 启动所有邮槽协程处理
func (s *Mailbox) Start() {
	for _, slot := range s.slots {
		slot.Sched()
	}
	if len(s.slots) > 0 {
		s.ch = make(chan os.Signal)
		s.done = make(chan os.Signal)
		signal.Notify(s.ch, os.Interrupt, os.Kill)
		go s.watch()
	}
}

/// 监视器
func (s *Mailbox) watch() {
	sig := <-s.ch
	close(s.ch)
	s.done <- sig
}

/// 获取下一个邮槽
func (s *Mailbox) GetNextSlot() ISlot {
	var slot ISlot
	if len(s.slots) > 0 {
		slot = s.slots[s.next] //atomic.AddInt64(&s.x, 1)%int64(len(s.cells))
		s.next++
		if s.next >= len(s.slots) {
			s.next = 0
		}
	}
	return slot
}

/// 等待退出
func (s *Mailbox) Wait() {
	if s.done != nil {
		<-s.done
		//Stop()或CTRL+C前执行清理
		s.clear()
		close(s.done)
	}
}

/// 手动清理
func (s *Mailbox) clear() {
	for _, slot := range s.slots {
		slot.Stop()
	}
}

/// 主动退出
func (s *Mailbox) Stop() {
	if s.ch != nil {
		//通知监视器退出
		s.ch <- os.Interrupt
	}
}
