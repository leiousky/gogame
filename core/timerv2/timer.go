package timerv2

import (
	"fmt"
	"reflect"
	"time"
)

type DelayCall struct {
	tid  uint32
	f    func(v ...interface{})
	args []interface{}
}

func (this *DelayCall) Call() {
	defer func() {
		//if err := recover(); err != nil {
		//	logger.Error(this.String(), "Call Error: ", err)
		//}
	}()

	this.f(this.args...)
}

func (this *DelayCall) String() string {
	funcType := reflect.TypeOf(this.f)
	return fmt.Sprintf("DelayCall function: %s. args: %v.", funcType.Name(), this.args)
}

func (this *DelayCall) GetTid() uint32 {
	return this.tid
}

type Timer struct {
	durations time.Duration
	delayCall *DelayCall
}

func NewTimer(durations time.Duration, f func(v ...interface{}), args []interface{}) *Timer {
	return &Timer{
		durations: durations,
		delayCall: &DelayCall{
			f:    f,
			args: args,
		},
	}
}

func (this *Timer) Run() {
	go func() {
		time.Sleep(this.durations)
		this.delayCall.Call()
	}()
}

func (this *Timer) GetDurations() time.Duration {
	return this.durations
}

func (this *Timer) GetFunc() *DelayCall {
	return this.delayCall
}
