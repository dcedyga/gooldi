package concurrency

import (
	"sync"
)

// Map is a map type that can be safely shared between
// goroutines that require read/write access to a map
type Map struct {
	lock  *sync.RWMutex
	items map[interface{}]interface{}
}

// MapItem contains a key/value pair item of a concurrent map
type MapItem struct {
	Key   interface{}
	Value interface{}
}

// NewMap creates a new concurrent map
func NewMap() *Map {
	cm := &Map{
		items: make(map[interface{}]interface{}),
		lock:  &sync.RWMutex{},
	}

	return cm
}

// Set adds an item to a concurrent map
func (cm *Map) Set(key, value interface{}) {

	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.items[key] = value

}

// Get retrieves the value for a concurrent map item
func (cm *Map) Get(key interface{}) (interface{}, bool) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	value, ok := cm.items[key]

	return value, ok
}

// GetKeyByItem - retrieves the key for a concurrent map item by item
func (cm *Map) GetKeyByItem(item interface{}) (interface{}, bool) {
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
func (cm *Map) Delete(key interface{}) bool {

	cm.lock.Lock()
	defer cm.lock.Unlock()
	_, ok := cm.items[key]
	if ok {
		delete(cm.items, key)
	}
	return ok

}

// Iter iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (cm *Map) Iter() <-chan MapItem {
	c := make(chan MapItem)

	f := func() {
		cm.lock.Lock()
		defer cm.lock.Unlock()

		for k, v := range cm.items {
			c <- MapItem{k, v}
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
func (cm *Map) IterWithCancel(cancel chan interface{}) <-chan MapItem {
	c := make(chan MapItem)

	f := func() {
		cm.lock.Lock()
		defer cm.lock.Unlock()

		for k, v := range cm.items {
			select {
			case <-cancel:
				close(c)
				return
			case c <- MapItem{k, v}:
			}
		}
		close(c)
	}
	go f()

	return c
}

//Len - length of the map
func (cm *Map) Len() int {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	return len(cm.items)
}

// GetKeys returns a slice of all the keys present
func (cm *Map) GetKeys() []interface{} {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	keys := []interface{}{}
	for i := range cm.items {
		keys = append(keys, i)
	}
	return keys
}
