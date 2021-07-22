package concurrency

import (
	"fmt"
	"sort"
	"sync"
)

type interfaceArray []interface{}

func (a interfaceArray) Len() int      { return len(a) }
func (a interfaceArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a interfaceArray) Less(i, j int) bool {
	s := fmt.Sprintf("%v", a[i])
	s2 := fmt.Sprintf("%v", a[j])
	return s < s2

}

// Map is a map type that can be safely shared between
// goroutines that require read/write access to a map
type SortedMap struct {
	lock  *sync.RWMutex
	keys  []interface{}
	dirty bool
	items map[interface{}]interface{}
}

// MapItem contains a key/value pair item of a concurrent map
type SortedMapItem struct {
	Key   interface{}
	Value interface{}
}

// NewMap creates a new concurrent map
func NewSortedMap() *SortedMap {
	cm := &SortedMap{
		items: make(map[interface{}]interface{}),
		keys:  []interface{}{},
		dirty: true,
		lock:  &sync.RWMutex{},
	}

	return cm
}

// Set adds an item to a concurrent map
func (cm *SortedMap) Set(key, value interface{}) {

	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.items[key] = value
	cm.dirty = true
}

// Get retrieves the value for a concurrent map item
func (cm *SortedMap) Get(key interface{}) (interface{}, bool) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	value, ok := cm.items[key]

	return value, ok
}

// GetItemByIndex retrieves the value for a concurrent map item
func (cm *SortedMap) GetSortedMapItemByIndex(index int) (*MapItem, bool) {
	cm.sort()
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	key := cm.keys[index]
	if value, ok := cm.items[key]; ok {
		return &MapItem{key, value}, ok
	}

	return nil, false
}

// GetKeyByItem - retrieves the key for a concurrent map item by item
func (cm *SortedMap) GetKeyByItem(item interface{}) (interface{}, bool) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	for k, v := range cm.items {
		if item == v {
			return k, true
		}
	}

	return nil, false
}

// Delete removes the value/key pair of a concurrent map item
func (cm *SortedMap) Delete(key interface{}) bool {

	cm.lock.Lock()
	defer cm.lock.Unlock()
	_, ok := cm.items[key]
	if ok {
		delete(cm.items, key)
		cm.dirty = true
	}
	return ok

}

// Iter iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (cm *SortedMap) Iter() <-chan MapItem {
	c := make(chan MapItem)

	f := func() {
		cm.sort()
		cm.lock.Lock()
		defer cm.lock.Unlock()
		for _, k := range cm.keys {
			c <- MapItem{k, cm.items[k]}
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
func (cm *SortedMap) IterWithCancel(cancel chan interface{}) <-chan MapItem {
	c := make(chan MapItem)

	f := func() {
		cm.sort()
		cm.lock.Lock()
		defer cm.lock.Unlock()
		for _, k := range cm.keys {
			select {
			case <-cancel:
				close(c)
				return
			case c <- MapItem{k, cm.items[k]}:
			}
		}
		close(c)
	}
	go f()

	return c
}
func (cm *SortedMap) sort() {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	if cm.dirty {
		cm.keys = make([]interface{}, 0, len(cm.items))
		for i := range cm.items {
			cm.keys = append(cm.keys, i)
		}
		sort.Sort(interfaceArray(cm.keys))
		cm.dirty = false
	}
}

//Len - length of the map
func (cm *SortedMap) Len() int {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	return len(cm.items)
}

// GetKeys returns a slice of all the keys present
func (cm *SortedMap) GetKeys() []interface{} {
	cm.sort()
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	return cm.keys
}

// func (cm *Map) SortAndClone() *Map {
// 	cm.lock.Lock()
// 	defer cm.lock.Unlock()
// 	m := NewMap()
// 	sk := make([]interface{}, len(cm.keys))
// 	copy(sk, cm.keys)
// 	sort.Sort(interfaceArray(sk))
// 	for _, k := range sk {
// 		m.Set(k, cm.items[k])
// 	}
// 	return m
// }
