package timerv2

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

/*
协程安全的timer
1.时钟原理说明：
1.1. 初始化一个三层时间轮:秒刻盘：0~59个节点, 分刻盘：0~59个节点, 时刻盘：0~12个节点;
1.2. 秒针由外界推动,每次循环为200毫秒，5次循环跳动一秒针，每跳一轮(60格),秒针复位至0,同时分针跳1格;
1.3. 同理分针每跳一轮(60格),分针复位至0,同时时针跳1格;
1.4. 最高层：时针跳一轮(12格）,时针复位至0，一个时间轮完整周期完成.
2.事件原理说明：
2.1. 设置时间为TimeOut的事件时,根据TimeOut算出发生此事件时刻的指针位置{TriggerHour,TriggerMin,TriggerSec};
2.2. 用{TriggerHour,TriggerMin,TriggerSec}与当前指针{NowHour,NowMin,NowSec}转换成系统毫米进行对比得出事件存放在哪一个指针(Tick);
2.3. 所有层的指针每跳到下一格(Tick01)都会触发格子的事件列表,处理每一个事件Event01：
2.3.1 根据事件Event01的剩余TimeOut算出Event01应该存在上一层(跳得更快)层的位置Pos;
2.3.2 把事件更新到新的Pos(更新TimeOut);
2.3.3 重复处理完Tick01里面所有的事件;
2.3.4 清空Tick01的事件;
2.3.5 最底层(跳最快)层所有的事件遇到指针Tick都会立即执行;
*/
const (
	//默认安全时间调度器的容量
	TIMERLEN = 2048
	//默认最大误差值100毫秒
	ERRORMAX = 100
	//默认最大触发队列缓冲大小
	TRIGGERMAX = 2048
	//默认hashwheel分级
	LEVEL = 12
)

func UnixTS() int64 {
	return time.Now().UnixNano() / 1e6
}

type ParamNull struct{}

type SafeTimer struct {
	//延迟调用的函数
	delayCall *DelayCall
	//调用的时间：单位毫秒
	unixts int64
}

// delay 以毫秒为单位， 一般输入可以使用 1000 标识1秒
func NewSafeTimer(tid uint32, delay int64, delayCall *DelayCall) *SafeTimer {
	unixts := UnixTS()
	if delay > 0 {
		unixts += delay
	}
	return &SafeTimer{
		delayCall: delayCall,
		unixts:    unixts,
	}
}

type SafeTimerScheduel struct {
	hashwheel   *HashWheel      // 时间轮的指针
	idGen       uint32          // 延时调用任务ID
	triggerChan chan *DelayCall // 调用任务输出通道
	cron        *cron.Cron      // 用于增加定时任务
	sync.RWMutex
}

func NewSafeTimerScheduel() *SafeTimerScheduel {
	scheduel := &SafeTimerScheduel{
		hashwheel:   NewHashWheel("wheel_hours", LEVEL, 3600*1e3, TIMERLEN),
		idGen:       0,
		triggerChan: make(chan *DelayCall, TRIGGERMAX),
		cron:        cron.New(cron.WithSeconds()),
	}

	//minute wheel
	minuteWheel := NewHashWheel("wheel_minutes", LEVEL, 60*1e3, TIMERLEN)
	//second wheel
	secondWheel := NewHashWheel("wheel_seconds", LEVEL, 1*1e3, TIMERLEN)
	minuteWheel.AddNext(secondWheel)
	scheduel.hashwheel.AddNext(minuteWheel)
	// 启动定时任务
	scheduel.cron.Start()
	// 启动延时任务
	go scheduel.StartScheduelLoop()
	return scheduel
}

func (t *SafeTimerScheduel) GetTriggerChannel() chan *DelayCall {
	return t.triggerChan
}

func (t *SafeTimerScheduel) Do() <-chan *DelayCall {
	return t.triggerChan
}

// 增加一个延时任务
func (t *SafeTimerScheduel) CreateTimer(delay int64, f func(v ...interface{}), args []interface{}) (uint32, error) {
	t.Lock()
	defer t.Unlock()

	t.idGen += 1
	// 当达到uint32 最大值时候 重置为1， 此时理论上不存在相同id的任务
	if t.idGen > 4294967290 {
		t.idGen = 1
	}
	d := &DelayCall{
		tid:  t.idGen,
		f:    f,
		args: args,
	}
	err := t.hashwheel.Add2WheelChain(t.idGen,
		NewSafeTimer(t.idGen, delay, d))
	if err != nil {
		return 0, err
	} else {
		return t.idGen, nil
	}
}

// 增加一个定时任务
func (t *SafeTimerScheduel) CreatCronFunc(spec string, cmd func()) (uint32, error) {
	t.Lock()
	defer t.Unlock()
	entryID, err := t.cron.AddFunc(spec, cmd)
	if err == nil {
		return uint32(entryID), nil
	} else {
		return 0, err
	}
}

// 删除一个定时任务
func (t *SafeTimerScheduel) DelCronFunc(cronId uint32) {
	t.cron.Remove(cron.EntryID(cronId))
}

func (t *SafeTimerScheduel) CancelTimer(timerId uint32) {
	t.hashwheel.RemoveFromWheelChain(timerId)
}

func (t *SafeTimerScheduel) StartScheduelLoop() {
	fmt.Printf("timer safe timer scheduelloop runing.")
	for {
		triggerList := t.hashwheel.GetTriggerWithIn(ERRORMAX)
		//trigger
		for _, v := range triggerList {
			//logger.Debug("want call: ", v.unixts, ".real call: ", UnixTS(), ".ErrorMS: ", UnixTS()-v.unixts)
			if math.Abs(float64(UnixTS()-v.unixts)) > float64(ERRORMAX) {
				fmt.Println("want call: ", v.unixts, ".real call: ", UnixTS(), ".ErrorMS: ", UnixTS()-v.unixts)
			}

			t.triggerChan <- v.delayCall
		}

		//wait for next loop
		time.Sleep(ERRORMAX / 2 * time.Millisecond)
	}
}
