package core

import (
	"games/comm/utils"
	"games/core/cell"
	"games/core/slot"
	"os"
	"os/signal"
	"time"
)

/// <summary>
/// IMailbox 邮槽管理器接口
/// <summary>
type IMailbox interface {
	/// 添加若干邮槽
	Add(d time.Duration, size int, creator cell.IWorkerCreator, num int)
	/// 添加一个邮槽
	AddOne(d time.Duration, size int, creator cell.IWorkerCreator) slot.ISlot
	/// 遍历每个邮槽
	Range(cb func(slot.ISlot, int))
	/// 启动所有邮槽协程处理
	Start()
	/// 获取下一个邮槽
	GetNextSlot() slot.ISlot
	/// 等待退出
	Wait()
	/// 主动退出
	Stop()
}

/// <summary>
/// Mailbox 邮槽管理器实现
/// <summary>
type Mailbox struct {
	slots []slot.ISlot
	next  int
	ch    chan os.Signal
	done  chan os.Signal
}

func NewMailBox() IMailbox {
	return &Mailbox{next: 0}
}

/// 添加若干邮槽
func (s *Mailbox) Add(d time.Duration, size int, creator cell.IWorkerCreator, num int) {
	for i := 0; i < num; i++ {
		slot := slot.NewMsgSlot(d, size, creator)
		s.slots = append(s.slots, slot)
	}
}

/// 添加一个邮槽
func (s *Mailbox) AddOne(d time.Duration, size int, creator cell.IWorkerCreator) slot.ISlot {
	slot := slot.NewMsgSlot(d, size, creator)
	s.slots = append(s.slots, slot)
	return slot
}

/// 遍历每个邮槽
func (s *Mailbox) Range(cb func(slot.ISlot, int)) {
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
func (s *Mailbox) GetNextSlot() slot.ISlot {
	var slot slot.ISlot
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
