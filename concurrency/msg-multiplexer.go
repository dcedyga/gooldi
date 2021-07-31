package concurrency

import (
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"
)

// MsgMultiplexerOption - option to initialize the MsgMultiplexer
type MsgMultiplexerOption func(*MsgMultiplexer)

// MsgMultiplexer - The default implementation MsgMultiplexer allows to create complex patterns where a Broadcaster can emit an message to multiple processors (consumers)
// that can potentially represent multiple processing systems, do the relevant calculation and multiplex the multiple outputs
// into a single channel for simplified consumption.
// Its main function is to Mulitplex a set of parallel processors that process a common initial concurrency.Message
// converging them into one channel,where the output is a concurrency.Message which Message property is a sortedmap of the output
// values of the processors grouped by initial concurrency.Message CorrelationKey and ordered by index value of each processor.
// Closure of MsgMultiplexer is handle by a concurrency.DoneHandler that allows to control they way a set of go routines
// are closed in order to prevent deadlocks and unwanted behaviour
// MsgMultiplexer outputs the multiplexed result in one channel using the channel bridge pattern.
// MsgMultiplexer default behaviour can be overridden by providing a MsgMultiplexerGetItemKeyFn to provide the comparison key of
// the items of a channel, with this function MsgMultiplexer
// has an algorithm to group the processed messages related to the same source into a SortedMap.

type MsgMultiplexer struct {
	id                      string
	inputChannels           *Map
	doneHandler             *DoneHandler
	preStreamMap            *Map
	inputChannelsLastRegLen int
	stream                  chan (<-chan interface{})
	lock                    *sync.RWMutex
	index                   int64
	minIndexKey             interface{}
	MsgType                 string
	getItemKeyFn            func(v interface{}) int64
	transformFn             func(mp *MsgMultiplexer, sm *SortedMap, correlationKey int64) interface{}
}

//NewMsgMultiplexer - Constructor
func NewMsgMultiplexer(dh *DoneHandler, msgType string, opts ...MsgMultiplexerOption) *MsgMultiplexer {
	id := uuid.NewV4().String()
	mp := &MsgMultiplexer{
		id:                      id,
		inputChannels:           NewMap(),
		doneHandler:             dh,
		inputChannelsLastRegLen: 0,
		preStreamMap:            NewMap(),
		stream:                  make(chan (<-chan interface{})),
		index:                   0,
		minIndexKey:             "",
		MsgType:                 msgType,
		lock:                    &sync.RWMutex{},
	}
	mp.getItemKeyFn = defaultGetItemKeyFn
	mp.transformFn = defaultTransformFn
	for _, opt := range opts {
		opt(mp)
	}
	go mp.doneRn()
	return mp
}

// doneRn - Checks when the MsgMultiplexer is done by listening to the closure of the DoneHandler.Done channel
func (mp *MsgMultiplexer) doneRn() {
	select {
	case <-mp.doneHandler.Done():
		mp.close()
	}
}

// MsgMultiplexerIndex - option to add the index value to the MsgMultiplexer
func MsgMultiplexerIndex(idx int64) MsgMultiplexerOption {
	return func(mp *MsgMultiplexer) {
		mp.index = idx
	}
}

// MsgMultiplexerGetItemKeyFn - option to add a function to resolve the key value
// of an item of the channel to the MsgMultiplexer
func MsgMultiplexerGetItemKeyFn(fn func(v interface{}) int64) MsgMultiplexerOption {
	return func(mp *MsgMultiplexer) {
		mp.getItemKeyFn = fn
	}
}

// MsgMultiplexerTransformFn - option to add a function to transform the SortedMap output into
// the desired output structure to the MsgMultiplexer
func MsgMultiplexerTransformFn(fn func(mp *MsgMultiplexer, sm *SortedMap, correlationKey int64) interface{}) MsgMultiplexerOption {
	return func(mp *MsgMultiplexer) {
		mp.transformFn = fn
	}
}

// ID - retrieves the Id of the MsgMultiplexer
func (mp *MsgMultiplexer) ID() string {
	return mp.id
}

// Index - retrieves the Index of the MsgMultiplexer
func (mp *MsgMultiplexer) Index() int64 {
	return mp.index
}

// Set - Registers a channel in the MsgMultiplexer, starts processing it
// and logs the length of the registered channels map.
func (mp *MsgMultiplexer) Set(key interface{}, value chan interface{}) {
	mp.lock.Lock()
	mp.inputChannels.Set(key, value)
	mp.inputChannelsLastRegLen = mp.inputChannels.Len()
	if mp.minIndexKey == "" {
		mp.minIndexKey = key
	} else if mp.isMinIndexKey(key) {
		mp.minIndexKey = key
	}
	mp.lock.Unlock()
	go mp.processChannel(key, value)

}

func (mp *MsgMultiplexer) isMinIndexKey(key interface{}) bool {
	s := fmt.Sprintf("%v", key)
	s2 := fmt.Sprintf("%v", mp.minIndexKey)
	return s < s2
}

// Get - Retrieves a channel reqistered in the MsgMultiplexer by key
func (mp *MsgMultiplexer) Get(key interface{}) (chan interface{}, bool) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	v, ok := mp.inputChannels.Get(key)
	if !ok {
		return nil, false
	}
	return v.(chan interface{}), true
}

// delete - Deletes a registered channel from the MsgMultiplexer map of inputChannels
// and logs the length of the remaining registered channels map.
func (mp *MsgMultiplexer) delete(key interface{}) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	mp.inputChannels.Delete(key)
	mp.inputChannelsLastRegLen = mp.inputChannels.Len()

}

// processChannel - Retrieves all the values of an inputChannel using range and when done it deletes it
// from the MsgMultiplexer map of inputChannels.
func (mp *MsgMultiplexer) processChannel(key interface{}, value chan interface{}) {
	for item := range value {
		select {
		default:
			cKey := mp.getItemKeyFn(item)
			mp.storeInPreStreamMap(cKey, key, item)
			mp.sendWhenReady(cKey)
		}
	}
	mp.delete(key)
}

//storeInPreStreamMap - store stream item in PreStreamMap per correlationKey
func (mp *MsgMultiplexer) storeInPreStreamMap(key int64, sequence interface{}, v interface{}) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	item, ok := mp.preStreamMap.Get(key)
	var rm *SortedMap
	if ok {
		rm = item.(*SortedMap)
		rm.Set(sequence, v) // Sequence
	} else {
		rm = NewSortedMap()
		rm.Set(sequence, v) // Sequence
		mp.preStreamMap.Set(key, rm)
	}
}

// close - Closes the MsgMultiplexer
func (mp *MsgMultiplexer) close() {
	fmt.Printf("Multiplexer - is closed\n")
}

func (mp *MsgMultiplexer) sendWhenReady(cKey int64) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	m, ok := (mp.preStreamMap.Get(cKey))
	if ok {
		sm := m.(*SortedMap)
		l := sm.Len()
		inputChannelsLastRegLen := mp.inputChannelsLastRegLen
		//fmt.Printf("len:%v,RegLen:%v\n", l, mp.inputChannelsLastRegLen)
		if l == inputChannelsLastRegLen {
			//send to stream
			e := mp.transformFn(mp, sm, cKey)
			s := make(chan interface{}, 1)
			s <- e
			close(s)
			mp.stream <- s
			//delete
			mp.preStreamMap.Delete(cKey)
		}
	}
}

// Iter iterates over the items in the Multiplexer
// Each item is sent over a channel, so that
// we can iterate over the it using the builtin range keyword
func (mp *MsgMultiplexer) Iter() chan interface{} {
	return Bridge(mp.doneHandler.Done(), mp.stream)
}

//PrintPreStreamMap - prints the preStreamMap
func (mp *MsgMultiplexer) PrintPreStreamMap() {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	fmt.Printf("preStreamMap:%v\n", mp.preStreamMap)
	l := mp.preStreamMap.Len()
	if l > 0 {
		for it := range mp.preStreamMap.Iter() {
			mop := it.Value.(*SortedMap)
			fmt.Printf("mop.len:%v\n", mop.Len())
		}
	}
}

/******************************************************************************
Plug and Play and Transformation functions
*******************************************************************************/

// defaultGetItemKeyFn - Gets the key value used by the comparator in order to group the Messages in the
// SortedMap. Can be overridden for a more generic implementation that is not Message Oriented
func defaultGetItemKeyFn(v interface{}) int64 {
	return v.(*Message).CorrelationKey
}

// MsgMultiplexerMessagePairGetItemKeyFn - Gets the key value used by the comparator in order to group the Messages in the
// SortedMap. Can be overridden for a more generic implementation that is not Message Oriented
func MsgMultiplexerMessagePairGetItemKeyFn(v interface{}) int64 {
	return v.(*MessagePair).CorrelationKey
}

// defaultTransformFn - Transforms the SortedMap output into an Message for future consumption as part of the
// output channel of the MsgMultiplexer. Can be overridden for a more generic implementation
func defaultTransformFn(mp *MsgMultiplexer, sm *SortedMap, correlationKey int64) interface{} {

	re := NewMessage(sm,
		mp.MsgType,
		MessageWithIndex(mp.index),
		MessageWithCorrelationKey(correlationKey),
	)

	return re
}

// MsgMultiplexerMessagePairTransformFn - Transforms the SortedMap output into an Message for future consumption as part of the
// output channel of the MsgMultiplexer. Can be overridden for a more generic implementation
func MsgMultiplexerMessagePairTransformFn(mp *MsgMultiplexer, sm *SortedMap, correlationKey int64) interface{} {

	inMsg, _ := sm.Get(mp.minIndexKey)

	in := inMsg.(*MessagePair).In

	out := NewMessage(sm,
		mp.MsgType,
		MessageWithIndex(mp.index),
		MessageWithCorrelationKey(correlationKey),
	)

	re := NewMessagePair(in, out)
	return re
}
