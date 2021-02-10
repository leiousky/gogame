package timer

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func test(a ...interface{}) {
	fmt.Println(a[0], "============", a[1])
}

var tt = int64(0)

var s *SafeTimerScheduel

func TestMain(m *testing.M) {
	s = NewSafeTimerScheduel()
	m.Run()
}

func Test(t *testing.T) {
	go func() {
		for {
			df := <-s.GetTriggerChannel()
			df.Call()
			atomic.AddInt64(&tt, -1)
		}
	}()
	go func() {
		i := 0
		for i < 50000 {
			s.CreateTimer(int64(rand.Int31n(3600*1e3)), test, []interface{}{22, 33})
			atomic.AddInt64(&tt, 1)
			time.Sleep(1 * time.Second)
			i += 1
		}
	}()
	go func() {
		ii := 0
		for ii < 50000 {
			s.CreateTimer(int64(rand.Int31n(3600*1e3)), test, []interface{}{22, 33})
			atomic.AddInt64(&tt, 1)
			time.Sleep(1 * time.Second)
			ii += 1
		}
	}()

	for {
		time.Sleep(60 * time.Second)
		fmt.Printf("last timer: ", atomic.LoadInt64(&tt))
	}
}

func TestGet(t *testing.T) {
	var a map[int]string
	fmt.Println(a)
	var f Foo
	fmt.Println(f)
	b := f.GetBar()
	fmt.Println(b)
	n := b.GetName()
	fmt.Println(n)
}

type Foo struct {
	bar *Bar
}
type Bar struct {
	name string
}

func (t Foo) GetBar() *Bar {
	return t.bar
}
func (t Bar) GetName() string {
	return t.name
}
