package net

import "sync/atomic"

var gSessionID int64

//
func createSessionID() int64 {
	return atomic.AddInt64(&gSessionID, 1)
}
