package concurrency

import (
	"fmt"
	"sort"
	"sync"
)

// SortedSlice type that can be safely shared between goroutines
type SortedSlice struct {
	lock  *sync.RWMutex
	items sortedSliceArray
	dirty bool
}

type sortedSliceArray []interface{}

// SortedSliceItem contains the index/value pair of an item in a
// concurrent slice
type SortedSliceItem struct {
	Index int
	Value interface{}
}

func (a sortedSliceArray) Len() int { return len(a) }
func (a sortedSliceArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a sortedSliceArray) Less(i, j int) bool {
	switch a[i].(type) {
	case int:
		return a[i].(int) < a[j].(int)
	case int8:
		return a[i].(int8) < a[j].(int8)
	case int16:
		return a[i].(int16) < a[j].(int16)
	case int32:
		return a[i].(int32) < a[j].(int32)
	case int64:
		return a[i].(int64) < a[j].(int64)
	case uint:
		return a[i].(uint) < a[j].(uint)
	case uint8:
		return a[i].(uint8) < a[j].(uint8)
	case uint16:
		return a[i].(uint16) < a[j].(uint16)
	case uint32:
		return a[i].(uint32) < a[j].(uint32)
	case uint64:
		return a[i].(uint64) < a[j].(uint64)
	case float32:
		return a[i].(float32) < a[j].(float32)
	case float64:
		return a[i].(float64) < a[j].(float64)
	default:
		s := fmt.Sprintf("%v", a[i])
		s2 := fmt.Sprintf("%v", a[j])
		return s < s2
	}

}

// NewSortedSlice creates a new concurrent slice
func NewSortedSlice() *SortedSlice {
	cs := &SortedSlice{
		items: make([]interface{}, 0),
		lock:  &sync.RWMutex{},
		dirty: true,
	}

	return cs
}

// Append adds an item to the concurrent slice
func (cs *SortedSlice) Append(item interface{}) {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.items = append(cs.items, item)
	cs.dirty = true

}

func (cs *SortedSlice) sort() {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	if cs.dirty {
		sort.Sort(sortedSliceArray(cs.items))
		cs.dirty = false
	}
}

//RemoveItemAtIndex removes the item at the specified index
func (cs *SortedSlice) RemoveItemAtIndex(index int) {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.items = append(cs.items[:index], cs.items[index+1:]...)
	cs.dirty = true
}

//IndexOf returns the index of a specific item
func (cs *SortedSlice) IndexOf(item interface{}) int {
	cs.sort()
	cs.lock.Lock()
	defer cs.lock.Unlock()
	for index, v := range cs.items {
		if item == v {
			return index
		}
	}
	return -1 //not found.
}

//GetItemAtIndex - Get item at index
func (cs *SortedSlice) GetItemAtIndex(index int) interface{} {
	cs.sort()
	cs.lock.Lock()
	defer cs.lock.Unlock()
	return cs.items[index]
}

// Iter iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builtin range keyword
func (cs *SortedSlice) Iter() <-chan SortedSliceItem {
	c := make(chan SortedSliceItem)

	f := func() {
		cs.sort()
		cs.lock.Lock()
		defer cs.lock.Unlock()
		for index, value := range cs.items {
			c <- SortedSliceItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}

// IterWithCancel iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builtin range keyword
// allows to pass a cancel chan to make the iteration cancelable
func (cs *SortedSlice) IterWithCancel(cancel chan interface{}) <-chan SortedSliceItem {
	c := make(chan SortedSliceItem)

	f := func() {
		cs.sort()
		cs.lock.Lock()
		defer cs.lock.Unlock()
		for index, value := range cs.items {
			select {
			case <-cancel:
				close(c)
				return
			case c <- SortedSliceItem{index, value}:
			}
		}
		close(c)
	}
	go f()

	return c
}

//Len - length of the slice
func (cs *SortedSlice) Len() int {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	return len(cs.items)
}

//Cap - capacity of the slice
func (cs *SortedSlice) Cap() int {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	return cap(cs.items)
}
