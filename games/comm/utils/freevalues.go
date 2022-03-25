package utils

import (
	"container/list"
	"runtime/debug"
	"sync"
)

/// <summary>
/// FreeValues
/// <summary>
type FreeValues struct {
	values *list.List
	new    func() interface{}
	lock   *sync.Mutex
}

func NewFreeValues() *FreeValues {
	return &FreeValues{values: list.New(), lock: &sync.Mutex{}}
}

func NewFreeValuesWith(new func() interface{}) *FreeValues {
	return &FreeValues{values: list.New(), new: new, lock: &sync.Mutex{}}
}

func (s *FreeValues) SetNew(new func() interface{}) {
	s.new = new
}

func (s *FreeValues) Len() uint32 {
	return uint32(s.values.Len())
}

func (s *FreeValues) Alloc() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.values.Len() > 0 {
		elem := s.values.Front()
		value := elem.Value
		s.values.Remove(elem)
		return value
	}
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	return s.new()
}

func (s *FreeValues) Free(value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values.PushBack(value)
}

func (s *FreeValues) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values.Init()
}

func (s *FreeValues) ResetValues(callback func(value interface{})) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var next *list.Element
	for elem := s.values.Front(); elem != nil; elem = next {
		next = elem.Next()
		callback(elem.Value)
		s.values.Remove(elem)
	}
}

func (s *FreeValues) Visit(callback func(i int32, value interface{})) {
	s.lock.Lock()
	defer s.lock.Unlock()
	i := int32(0)
	for elem := s.values.Front(); elem != nil; elem = elem.Next() {
		callback(i, elem.Value)
		i++
	}
}
