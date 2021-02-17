package concurrency

import "sync"

// Slice type that can be safely shared between goroutines
type Slice struct {
	lock  *sync.RWMutex
	items []interface{}
}

// SliceItem contains the index/value pair of an item in a
// concurrent slice
type SliceItem struct {
	Index int
	Value interface{}
}

// NewSlice creates a new concurrent slice
func NewSlice() *Slice {
	cs := &Slice{
		items: make([]interface{}, 0),
		lock:  &sync.RWMutex{},
	}

	return cs
}

// Append adds an item to the concurrent slice
func (cs *Slice) Append(item interface{}) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	cs.items = append(cs.items, item)
}

//RemoveItemAtIndex removes the item at the specified index
func (cs *Slice) RemoveItemAtIndex(index int) {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.items = append(cs.items[:index], cs.items[index+1:]...)
}

//IndexOf returns the index of a specific item
func (cs *Slice) IndexOf(item interface{}) int {
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
func (cs *Slice) GetItemAtIndex(index int) interface{} {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	return cs.items[index]
}

// Iter iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builtin range keyword
func (cs *Slice) Iter() <-chan SliceItem {
	c := make(chan SliceItem)

	f := func() {
		cs.lock.Lock()
		defer cs.lock.Unlock()
		for index, value := range cs.items {
			c <- SliceItem{index, value}
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
func (cs *Slice) IterWithCancel(cancel chan interface{}) <-chan SliceItem {
	c := make(chan SliceItem)

	f := func() {
		cs.lock.Lock()
		defer cs.lock.Unlock()
		for index, value := range cs.items {
			select {
			case <-cancel:
				close(c)
				return
			case c <- SliceItem{index, value}:
			}
		}
		close(c)
	}
	go f()

	return c
}

//Len - length of the slice
func (cs *Slice) Len() int {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	return len(cs.items)
}

//Cap - capacity of the slice
func (cs *Slice) Cap() int {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	return cap(cs.items)
}
