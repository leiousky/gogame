package utils

import (
	"container/list"
)

/// <summary>
/// Pair key-value
/// <summary>
type Pair struct {
	Key interface{}
	Val interface{}
}

/// <summary>
/// Orderedmap 排序map
/// <summary>
type Orderedmap struct {
	list *list.List
}

func NewOrderedmap() *Orderedmap {
	return &Orderedmap{list: list.New()}
}

func (s *Orderedmap) Insert(key interface{}, value interface{}, compare func(a, b interface{}) bool) {
	pos := s.list.Front()
	for ; pos != nil; pos = pos.Next() {
		if !compare(key, pos.Value.(*Pair).Key) {
			data := &Pair{Key: key, Val: value}
			s.list.InsertBefore(data, pos)
			break
		}
	}
	if pos == nil {
		data := &Pair{Key: key, Val: value}
		s.list.PushBack(data)
	}
}

func (s *Orderedmap) Top() (interface{}, interface{}) {
	if elem := s.list.Front(); elem != nil {
		data := elem.Value.(*Pair)
		return data.Key, data.Val
	}
	return nil, nil
}

func (s *Orderedmap) Front() *list.Element {
	return s.list.Front()
}

func (s *Orderedmap) Pop() {
	if elem := s.list.Front(); elem != nil {
		s.list.Remove(elem)
	}
}

func (s *Orderedmap) Empty() bool {
	return s.list.Len() == 0
}
