package timer

import (
	//_ "container/heap"
	"games/comm/utils"
	"log"

	"sync/atomic"
)

// https://www.ibm.com/developerworks/cn/linux/l-cn-timers/
// https://blog.csdn.net/yueguanghaidao/article/details/46290539
// https://github.com/cloudwu/skynet/blob/master/skynet-src/skynet_timer.c

type TimerCallback func(timerID uint32, dt int32, args interface{}) bool

/// <summary>
/// ScopedTimer 基于最小堆(最小生成树)实现的线程局部定时器
/// ScopedTimer 提供给线程内部使用，所以是安全的
/// <summary>
type ScopedTimer interface {
	/// 定时器协程ID
	ThreadID() uint32
	/// 不指定回调和timerID
	CreateTimer(delay, interval int32, args interface{}) uint32
	/// 不指定回调指定timerID
	CreateTimerWithID(timerID uint32, delay, interval int32, args interface{}) uint32
	/// 指定回调和timerID
	CreateTimerWithIDCB(timerID uint32, delay, interval int32, handler TimerCallback, args interface{}) uint32
	/// 指定回调不指定timerID
	CreateTimerWithCB(delay, interval int32, handler TimerCallback, args interface{}) uint32
	/// 撤销定时器
	RemoveTimer(timerID uint32)
	/// 撤销所有
	RemoveTimers()
	/// 轮询定时回调 默认 handler(timerID, dt, args)
	/// 如果创建定时器时指定了回调函数handler
	/// 则执行handler(timerID, dt, args)回调，否则执行update(timerID, dt, args)回调
	Poll(pid uint32, update TimerCallback) bool
}

/// <summary>
/// timerEvent 定时器事件
/// <summary>
type timerEvent struct {
	// 定时器ID
	timerID uint32
	// 延迟执行等待(s)
	//delay int32
	// 执行间隔时间(s)
	interval int32
	// 上次开始执行时间
	last utils.Timestamp
	// 下次开始执行时间
	//expr Timestamp
	// 回调函数
	handler TimerCallback
	// 回调参数
	args interface{}
}

/// <summary>
/// scopedTimer 基于最小堆(最小生成树)实现的定时器
/// <summary>
type scopedTimer struct {
	x        uint32           // 用于自动生成timerID
	tid      uint32           // 定时器所属goroutine
	timers   utils.Orderedmap // 排序 map[timestamp] = timer
	timerIDs map[uint32]bool  // 保存要删除的timerID集合
}

func NewScopedTimer(tid uint32) ScopedTimer {
	if utils.GoroutineID() != tid {
		log.Fatalln("NewScopedTimer")
	}
	return &scopedTimer{
		tid:      tid,
		timers:   *utils.NewOrderedmap(),
		timerIDs: map[uint32]bool{}}
}

/// 定时器协程ID
func (s *scopedTimer) ThreadID() uint32 {
	return s.tid
}

/// 撤销定时器
func (s *scopedTimer) RemoveTimer(timerID uint32) {
	s.addRemoves(timerID)
}

/// 添加到撤销表
func (s *scopedTimer) addRemoves(timerID uint32) {
	// 线程安全
	s.AssertInThread(utils.GoroutineID())
	if timerID != 0 {
		s.timerIDs[timerID] = true
	}
}

/// 判断是否撤销
func (s *scopedTimer) isRemoveID(timerID uint32) bool {
	if _, ok := s.timerIDs[timerID]; ok {
		// timerID 在删除表中则移除 timerIDs = append(timerIDs[:i], timerIDs[i+1:]...)
		delete(s.timerIDs, timerID)
		return true
	}
	return false
}

/// 撤销所有
func (s *scopedTimer) RemoveTimers() {

}

/// 不带回调带ID
func (s *scopedTimer) CreateTimerWithID(timerID uint32, delay, interval int32, args interface{}) uint32 {
	return s.createTimer(timerID, delay, interval, nil, args)
}

/// 不带回调不带ID
func (s *scopedTimer) CreateTimer(delay, interval int32, args interface{}) uint32 {
	return s.CreateTimerWithID(atomic.AddUint32(&s.x, 1), delay, interval, args)
}

/// 带回调带ID
func (s *scopedTimer) CreateTimerWithIDCB(timerID uint32, delay, interval int32, handler TimerCallback, args interface{}) uint32 {
	return s.createTimer(timerID, delay, interval, handler, args)
}

/// 带回调不带ID
func (s *scopedTimer) CreateTimerWithCB(delay, interval int32, handler TimerCallback, args interface{}) uint32 {
	return s.CreateTimerWithIDCB(atomic.AddUint32(&s.x, 1), delay, interval, handler, args)
}

/// 比较大小
func compare(a, b interface{}) bool {
	return a.(utils.Timestamp).Greater(b.(utils.Timestamp))
}

/// 安全断言
func (s *scopedTimer) AssertInThread(pollingThreadID uint32) {
	if s.tid != pollingThreadID {
		log.Fatalf(" scopedTimer::assertInThread\n")
	}
}

/// 带回调带ID
func (s *scopedTimer) createTimer(timerID uint32, delay, interval int32, handler TimerCallback, args interface{}) uint32 {
	// 线程安全
	s.AssertInThread(utils.GoroutineID())
	// 创建 timer
	timer := &timerEvent{timerID: timerID, interval: interval, last: utils.TimeNowMilliSec(), handler: handler, args: args}
	// 放在 map[timestamp] = timer 中，并对 timestamp 进行关键字排序
	s.timers.Insert(utils.TimeAdd(timer.last, delay), timer, compare)
	// 打印调试
	// s.Keys()
	// 返回定时器ID
	return timerID
}

/// 从栈顶节点开始打印
func (s *scopedTimer) Keys() {
	i := 0
	for elem := s.timers.Front(); elem != nil; elem = elem.Next() {
		key := elem.Value.(*utils.Pair).Key.(utils.Timestamp)
		val := elem.Value.(*utils.Pair).Val.(*timerEvent)
		log.Printf("--- *** ScopedTimer[%d:%v] = %d", i, key.SinceUnixEpoch(), val.timerID)
		i++
	}
}

/// 定时器轮询 true定时器已空 false定时器不空
func (s *scopedTimer) Poll(pid uint32, update TimerCallback) bool {
	// 线程安全
	s.AssertInThread(pid)
	if s.timers.Empty() {
		return true
	}
	// 进入循环
	for {
		//log.Printf("--- *** ScopedTimer:: Poll %s...", CreateToken())
		now := utils.TimeNowMilliSec()
		// 取出栈顶Timestamp
		k, v := s.timers.Top()
		ts := k.(utils.Timestamp)
		t := v.(*timerEvent)
		if ts.Greater(now) {
			return false
		}
		// 先移除
		s.timers.Pop()
		// 判断是否撤销
		if s.isRemoveID(t.timerID) {
			// 删除
		} else if t.handler != nil { // 先执行handler回调如果有的话
			// 执行handler回调 handler(timerID, elapsed, args)
			if t.handler(t.timerID, utils.TimeDiff(now, t.last), t.args) {
				// 下次开始执行时间，从当前handler执行之后开始算
				if t.interval > 0 {
					t.last = now
					// 再次添加到有序表
					s.timers.Insert(utils.TimeNowMilliSec().Add(t.interval), t, compare)
				} else {
					// 不再需要则销毁
				}
			} else {
				// 不再需要则销毁
			} // 否则执行update回调如果有的话
		} else if update != nil {
			// 执行update回调 update(timerID, elapsed, args)
			if update(t.timerID, utils.TimeDiff(now, t.last), t.args) {
				// 下次开始执行时间，从当前update执行之后开始算
				if t.interval > 0 {
					t.last = now
					// 再次添加到有序表
					s.timers.Insert(utils.TimeNowMilliSec().Add(t.interval), t, compare)
				} else {
					// 不再需要则销毁
				}
			} else {
				// 不再需要则销毁
			}
		} else {
			// 不再需要则销毁
		}
		// 容器已空则返回
		if s.timers.Empty() {
			return true
		}
	}
}
