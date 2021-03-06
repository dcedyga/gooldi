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
}

type sortedSliceArray []interface{}

// SortedSliceItem contains the index/value pair of an item in a
// concurrent slice
type SortedSliceItem struct {
	Index int
	Value interface{}
}

func (a sortedSliceArray) Len() int      { return len(a) }
func (a sortedSliceArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortedSliceArray) Less(i, j int) bool {
	s := fmt.Sprintf("%v", a[i])
	s2 := fmt.Sprintf("%v", a[j])
	return s < s2

}

// NewSortedSlice creates a new concurrent slice
func NewSortedSlice() *SortedSlice {
	cs := &SortedSlice{
		items: make([]interface{}, 0),
		lock:  &sync.RWMutex{},
	}

	return cs
}

// Append adds an item to the concurrent slice
func (cs *SortedSlice) Append(item interface{}) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	cs.items = append(cs.items, item)
	sort.Sort(sortedSliceArray(cs.items))

}

//RemoveItemAtIndex removes the item at the specified index
func (cs *SortedSlice) RemoveItemAtIndex(index int) {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.items = append(cs.items[:index], cs.items[index+1:]...)
}

//IndexOf returns the index of a specific item
func (cs *SortedSlice) IndexOf(item interface{}) int {
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
