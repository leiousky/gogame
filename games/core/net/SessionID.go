package net

import "sync/atomic"

var gConnID int64

func NewConnID() int64 {
	return atomic.AddInt64(&gConnID, 1)
}
