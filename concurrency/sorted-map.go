package concurrency

import (
	"fmt"
	"sort"
	"sync"
)

// SortedMap is a sorted map type that can be safely shared between
// goroutines that require read/write access to a map
type SortedMap struct {
	items map[interface{}]interface{}
	keys  interfaceArray
	lock  *sync.RWMutex
}

// SortedMapItem contains a key/value pair item of a concurrent map
type SortedMapItem struct {
	Key   interface{}
	Value interface{}
}

type interfaceArray []interface{}

func (a interfaceArray) Len() int      { return len(a) }
func (a interfaceArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a interfaceArray) Less(i, j int) bool {
	s := fmt.Sprintf("%v", a[i])
	s2 := fmt.Sprintf("%v", a[j])
	return s < s2

}

// NewSortedMap creates a new concurrent sorted map
func NewSortedMap() *SortedMap {
	return &SortedMap{
		items: make(map[interface{}]interface{}),
		keys:  []interface{}{},
		lock:  &sync.RWMutex{},
	}
}

// Set adds an item to a concurrent map
func (sm *SortedMap) Set(key, value interface{}) {
	if index := sm.indexOf(key); index == -1 {
		sm.lock.Lock()
		defer sm.lock.Unlock()
		sm.items[key] = value
		sm.keys = append(sm.keys, key)
		sort.Sort(interfaceArray(sm.keys))
	}

}

// Get retrieves the value for a concurrent map item
func (sm *SortedMap) Get(key interface{}) (interface{}, bool) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	value, ok := sm.items[key]

	return value, ok
}

// GetMapItemByIndex retrieves the SortedMapItem for a concurrent map item given the index
func (sm *SortedMap) GetMapItemByIndex(index int) (*SortedMapItem, bool) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	if index < len(sm.keys) {
		key := sm.keys[index]
		value, ok := sm.items[key]
		i := &SortedMapItem{Key: key, Value: value}
		return i, ok
	}
	return nil, false
}

// GetByIndex retrieves the value for a concurrent map item given the index
func (sm *SortedMap) GetByIndex(index int) (interface{}, bool) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	if index < len(sm.keys) {
		key := sm.keys[index]
		value, ok := sm.items[key]

		return value, ok
	}
	return nil, false
}

// GetKeyByIndex retrieves the key for a concurrent map item given the index
func (sm *SortedMap) GetKeyByIndex(index int) (interface{}, bool) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	if index < len(sm.keys) {
		key := sm.keys[index]
		return key, true
	}
	return nil, false
}

// GetKeyByItem retrieves the key for a concurrent map item
func (sm *SortedMap) GetKeyByItem(item interface{}) (interface{}, bool) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	for k, v := range sm.items {
		if item == v {
			return k, true
		}
	}

	return nil, false
}

func (sm *SortedMap) indexOf(item interface{}) int {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	for k, v := range sm.keys {
		if item == v {
			return k
		}
	}
	return -1 //not found.
}

// Delete removes the value/key pair of a concurrent sorted map item
func (sm *SortedMap) Delete(key interface{}) {

	if index := sm.indexOf(key); index > -1 {
		sm.lock.Lock()
		delete(sm.items, key)
		sm.keys = append(sm.keys[:index], sm.keys[index+1:]...)
		sm.lock.Unlock()
	}

}

// Iter iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (sm *SortedMap) Iter() <-chan SortedMapItem {
	c := make(chan SortedMapItem)

	f := func() {
		sm.lock.Lock()
		defer sm.lock.Unlock()

		for _, k := range sm.keys {
			v := sm.items[k]
			c <- SortedMapItem{k, v}
		}
		close(c)
	}
	go f()

	return c
}

// IterWithCancel iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
// allows to pass a cancel chan to make the iteration cancelable
func (sm *SortedMap) IterWithCancel(cancel chan interface{}) <-chan SortedMapItem {
	c := make(chan SortedMapItem)
	f := func() {
		sm.lock.Lock()
		defer sm.lock.Unlock()
		for _, k := range sm.keys {
			v := sm.items[k]
			select {
			case <-cancel:
				close(c)
				return
			case c <- SortedMapItem{k, v}:
			}
		}
		close(c)
	}
	go f()

	return c
}

//Len - length of the sortedmap
func (sm *SortedMap) Len() int {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	return len(sm.items)
}
