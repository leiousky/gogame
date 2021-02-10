package utils

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
)

/// 出错打印堆栈
func SafeCall(f func()) (err error) {
	CheckPanic()
	f()
	return
}

/// panic检查
func CheckPanic() {
	defer func() {
		if err := recover(); err != nil {
			printStack(err)
		}
	}()
}

/// 打印堆栈信息
func printStack(err interface{}) {
	log.Println(fmt.Sprintf("stack: %v\n", err))
	b := make([]byte, 4096)
	n := runtime.Stack(b, false)
	log.Println(string(b[:n]))
}

/// 获取协程ID
func GoroutineID() uint32 {
	CheckPanic()
	// b := make([]byte, 64)
	// b = b[:runtime.Stack(b, false)]
	// b = bytes.TrimPrefix(b, []byte("goroutine "))
	// b = b[:bytes.IndexByte(b, ' ')]
	// n, err := strconv.ParseUint(string(b), 10, 64)
	// if err != nil {
	// 	panic(fmt.Sprintf("goroutineID panic: %v\n", err))
	// }
	// return int(n)
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	str := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	ID, err := strconv.Atoi(str)
	if err != nil {
		panic(fmt.Sprintf("GoroutineID panic: %v\n", err))
	}
	return uint32(ID)
}
