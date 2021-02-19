package concurrency

import (
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

// MsgMultiplexerOption - option to initialize the MsgMultiplexer
type MsgMultiplexerOption func(*MsgMultiplexer)

// MsgMultiplexer - The default implementation MsgMultiplexer allows to create complex patterns where a Broadcaster can emit an event to multiple processors (consumers)
// that can potentially represent multiple processing systems, do the relevant calculation and multiplex the multiple outputs
// into a single channel for simplified consumption.
// Its main function is to Mulitplex a set of parallel processors that process a common initial concurrency.Event/ Message
// converging them into one channel,where the output is an Event which Event.OutMessage is a sortedmap of the output
// values of the processors grouped by initial concurrency.Event/Message and ordered by sequence value of each processor.
// Closure of MsgMultiplexer is handle by a concurrency.DoneHandler that allows to control they way a set of go routines
// are closed in order to prevent deadlocks and unwanted behaviour
// MsgMultiplexer outputs the multiplexed result in one channel using the channel bridge pattern.
// MsgMultiplexer default behaviour can be overridden by providing a MsgMultiplexerGetChannelItemKeyFn to provide the comparison key of
// the items of a channel, MsgMultiplexerGetLastRegKeyFn which should give the key to compare to. With these two items MsgMultiplexer
// has an algorithm to group the processed messages related to the same source into a SortedMap.
// MsgMultiplexerGetChannelItemSequenceFn allows to get the sequence order of the relevant channel and MsgMultiplexerTransformFn
// allows to transform the output into the desired structure.
type MsgMultiplexer struct {
	id                      string
	inputChannels           *Map
	doneHandler             *DoneHandler
	outputMap               *SortedMap
	inputChannelsCountLog   *SortedMap
	inputChannelsLastRegKey int64
	inputChannelsLastRegLen int
	preStream               chan (<-chan interface{})
	stream                  chan (<-chan interface{})
	lock                    *sync.RWMutex
	sequence                interface{}
	getKeyFn                func(v interface{}) int64
	getLastRegisteredKey    func() int64
	getSequenceFn           func(v interface{}) interface{}
	transformFn             func(mp *MsgMultiplexer, sm *SortedMap) interface{}
}

// inputChannelsCountLogItem - keeps track of the length of the inputChannels Map when a channel is registered.
// inputChannelsCountLogItem is stored in inputChannelsCountLog SortedMap with a key that references the
// time that the channel is registered. The Multiplexer algorithm uses this structure to derive how many
// messages need to gather before sending the outputMap related to an initial input Event/ Message.
type inputChannelsCountLogItem struct {
	channel interface{}
	length  int
}

//NewMsgMultiplexer - Constructor
func NewMsgMultiplexer(dh *DoneHandler, opts ...MsgMultiplexerOption) *MsgMultiplexer {
	id := uuid.NewV4().String()
	mp := &MsgMultiplexer{
		id:                      id,
		inputChannels:           NewMap(),
		doneHandler:             dh,
		inputChannelsLastRegKey: -1,
		inputChannelsLastRegLen: 0,
		outputMap:               NewSortedMap(),
		inputChannelsCountLog:   NewSortedMap(),
		preStream:               make(chan (<-chan interface{})),
		stream:                  make(chan (<-chan interface{})),
		sequence:                0,
		lock:                    &sync.RWMutex{},
	}
	mp.getLastRegisteredKey = mp.defaultGetLastRegisteredKey
	mp.getKeyFn = mp.defaultGetKeyFn
	mp.getSequenceFn = mp.defaultGetSequenceFn
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

// MsgMultiplexerSequence - option to add a sequence value to the MsgMultiplexer
func MsgMultiplexerSequence(seq interface{}) MsgMultiplexerOption {
	return func(mp *MsgMultiplexer) {
		mp.sequence = seq
	}
}

// MsgMultiplexerGetChannelItemKeyFn - option to add a function to resolve the key value
// of an item of the channel to the MsgMultiplexer
func MsgMultiplexerGetChannelItemKeyFn(fn func(v interface{}) int64) MsgMultiplexerOption {
	return func(mp *MsgMultiplexer) {
		mp.getKeyFn = fn
	}
}

// MsgMultiplexerGetLastRegKeyFn - option to add a function to resolve the last registered key value
// for later comparison with the key of an item of the channel to the MsgMultiplexer
func MsgMultiplexerGetLastRegKeyFn(fn func() int64) MsgMultiplexerOption {
	return func(mp *MsgMultiplexer) {
		mp.getLastRegisteredKey = fn
	}
}

// MsgMultiplexerGetChannelItemSequenceFn - option to add a function to resolve the sequence value
// of an item of the channel to the MsgMultiplexer
func MsgMultiplexerGetChannelItemSequenceFn(fn func(v interface{}) interface{}) MsgMultiplexerOption {
	return func(mp *MsgMultiplexer) {
		mp.getSequenceFn = fn
	}
}

// MsgMultiplexerTransformFn - option to add a function to transform the SortedMap output into
// the desired output structure to the MsgMultiplexer
func MsgMultiplexerTransformFn(fn func(mp *MsgMultiplexer, sm *SortedMap) interface{}) MsgMultiplexerOption {
	return func(mp *MsgMultiplexer) {
		mp.transformFn = fn
	}
}

// ID - retrieves the Id of the MsgMultiplexer
func (mp *MsgMultiplexer) ID() string {
	return mp.id
}

// Sequence - retrieves the Sequence of the MsgMultiplexer
func (mp *MsgMultiplexer) Sequence() interface{} {
	return mp.sequence
}

// Set - Registers a channel in the MsgMultiplexer, starts processing it
// and logs the length of the registered channels map.
func (mp *MsgMultiplexer) Set(key interface{}, value chan interface{}) {
	mp.lock.Lock()
	mp.inputChannels.Set(key, value)
	mp.inputChannelsLastRegLen = mp.inputChannels.Len()
	item := inputChannelsCountLogItem{channel: key, length: mp.inputChannelsLastRegLen}
	mp.inputChannelsLastRegKey = mp.getLastRegisteredKey()
	mp.inputChannelsCountLog.Set(mp.inputChannelsLastRegKey, item)
	mp.lock.Unlock()
	go mp.processChannel(key, value)

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
func (mp *MsgMultiplexer) delete(key interface{}, keyToRegister int64) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	mp.inputChannelsLastRegKey = keyToRegister
	mp.inputChannels.Delete(key)
	mp.inputChannelsLastRegLen = mp.inputChannels.Len()
	item := inputChannelsCountLogItem{channel: key, length: mp.inputChannelsLastRegLen}
	mp.inputChannelsCountLog.Set(mp.inputChannelsLastRegKey, item)

}

// processChannel - Retrieves all the values of an inputChannel using range and when done it deletes it
// from the MsgMultiplexer map of inputChannels.
func (mp *MsgMultiplexer) processChannel(key interface{}, value chan interface{}) {
	var lastValue interface{}
	for item := range value {
		select {
		default:
			s := make(chan interface{}, 1)
			s <- item
			close(s)
			mp.preStream <- s
			lastValue = item
		}
	}
	keyToRegister := mp.getKeyFn(lastValue) + 1
	mp.delete(key, keyToRegister)
}

// getChannelsCountLogItemLenFromTime - Gets the length stored in inputChannelsCountLogItem for
// the closest item in the past for a given point of time.
func (mp *MsgMultiplexer) getChannelsCountLogItemLenFromTime(time int64) int {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	l := mp.inputChannelsCountLog.Len()
	for i := l - 1; i >= 0; i-- {
		k, ok := mp.inputChannelsCountLog.GetMapItemByIndex(i)
		if !ok {
			return 0
		} else if ok && time > k.Key.(int64) {
			return k.Value.(inputChannelsCountLogItem).length
		}
	}
	return 0
}

// Start - starts the main process of the MsgMultiplexer
func (mp *MsgMultiplexer) Start() {
	go mp.mainProcess()
}

// close - Closes the MsgMultiplexer
func (mp *MsgMultiplexer) close() {
	fmt.Printf("Multiplexer - is closed\n")
}

// mainProcess - The main process of the multiplexer it retrieves all the input messages that are on
// the preStream bridge and checks when and output result is ready by grouping all the messages into a
// SortedMap by a defined key, in the case of the default implementation with Events this key is the
// TimeInNano of the InitMessage of the Event, which is unique and represents the time the initial event
// was broadcasted. The SortedMap also sorts the output by the sequence order of the processor, producing a
// deterministic output. Once a SortedMap is sent it will delete the entry from the OutputMap
func (mp *MsgMultiplexer) mainProcess() {
	i := 0
	drain := false

	for v := range mp.bridgeIter() {
		nextItem := false
	loop:
		for !nextItem && !drain {
			select {
			case <-mp.doneHandler.Done():
				nextItem = true
				drain = true
				continue loop
			default:
				inputChannelsLastRegKey := mp.getInputChannelsLastRegKey()
				if inputChannelsLastRegKey != -1 {
					key := mp.getKeyFn(v) // getTime
					sequence := mp.getSequenceFn(v)
					mp.storeInOutputSortedMap(key, sequence, v)
					keysToDelete := mp.sendWhenReady()
					mp.cleanUpOutputMap(keysToDelete)
					nextItem = true
					continue loop
				}
			}
		}
		i++
	}
	mp.lock.Lock()
	fmt.Printf("Multiplexer Total Messages Multiplexed: %v,\n", mp.outputMap.Len())
	mp.lock.Unlock()
	fmt.Printf("Multiplexer Total Items Multiplexed: %v,\n", i)
}

// defaultGetLastRegisteredKey - gets the last registered key when a new channel is added or removed. It
// is used in conjunction with the defaultGetKeyFn to extract the length of the output SortedMap.
// Can be overridden for a more generic implementation that is not Event Oriented
func (mp *MsgMultiplexer) defaultGetLastRegisteredKey() int64 {
	return time.Now().UnixNano()
}

// defaultGetKeyFn - Gets the key value used by the comparator in order to group the Events in the
// SortedMap, in this case the TimeInNano property of the InitMessage of the event.
// Can be overridden for a more generic implementation that is not Event Oriented
func (mp *MsgMultiplexer) defaultGetKeyFn(v interface{}) int64 {
	return v.(*Event).InitMessage.TimeInNano
}

// defaultGetSequenceFn - Gets the sequence value of the event to be able to sort it when being added to the
// SortedMap. Can be overridden for a more generic implementation that is not Event Oriented
func (mp *MsgMultiplexer) defaultGetSequenceFn(v interface{}) interface{} {
	return v.(*Event).Sequence
}

// defaultTransformFn - Transforms the SortedMap output into an Event for future consumption as part of the
// output channel of the MsgMultiplexer. Can be overridden for a more generic implementation that is not Event
// Oriented
func defaultTransformFn(mp *MsgMultiplexer, sm *SortedMap) interface{} {
	var initmsg *Message

	eventInitMsgSeq := NewSlice()
	resultSlice := NewSlice()
	for item := range sm.Iter() {
		e := item.Value.(*Event)
		if initmsg == nil {
			initmsg = e.InitMessage
		}
		if e.OutMessage != nil {
			es := e.InMessageSequence
			es.Append(e.OutMessage)
			eventInitMsgSeq.Append(es)
			resultSlice.Append(e.OutMessage)
		}
	}
	rmsg := &Message{
		ID:         uuid.NewV4().String(),
		Message:    resultSlice,
		TimeInNano: initmsg.TimeInNano,
		MsgType:    initmsg.MsgType,
	}
	initMsgSeq := *NewSlice()
	initMsgSeq.Append(eventInitMsgSeq)
	re := &Event{
		InitMessage:       initmsg,
		InMessageSequence: initMsgSeq,
		OutMessage:        rmsg,
		Sequence:          mp.sequence,
	}
	return re
}

// getInputChannelsLastRegKey - Gets the inputChannelsLastRegKey property in a threadsafe manner
func (mp *MsgMultiplexer) getInputChannelsLastRegKey() int64 {
	mp.lock.Lock()
	inputChannelsLastRegKey := mp.inputChannelsLastRegKey
	mp.lock.Unlock()
	return inputChannelsLastRegKey
}

// getInputChannelsLastRegLen - Gets the inputChannelsLastRegLen property in a threadsafe manner
func (mp *MsgMultiplexer) getInputChannelsLastRegLen() int {
	mp.lock.Lock()
	inputChannelsLastRegLen := mp.inputChannelsLastRegLen
	mp.lock.Unlock()
	return inputChannelsLastRegLen
}

// storeInOutputSortedMap - stores the item of an input channel as part of the SortedMap result in the outputMap, with the
// common key of all the input messages, in the case of the Event implementation the TimeInNano of the InitMessage of the
// Event, which is unique. It detects if the SortedMap exists and creates a new one in case of not being present.
func (mp *MsgMultiplexer) storeInOutputSortedMap(key int64, sequence interface{}, v interface{}) {
	item, ok := mp.outputMap.Get(key)
	var rm *SortedMap
	if ok {
		rm = item.(*SortedMap)
		rm.Set(sequence, v) // Sequence
	} else {
		rm = NewSortedMap()
		rm.Set(sequence, v) // Sequence
		mp.outputMap.Set(key, rm)
	}
}

// sendWhenReady - Detects if the SortedMap has the required length based on the comparison of the key and lastRegKey.
// If the result Map is ready (has the required length) is sent to the output stream for consumption, using the bridge pattern
// and stores the key for later deletion.
func (mp *MsgMultiplexer) sendWhenReady() []int64 {
	keysToDelete := []int64{}
	inputChannelsLastRegLen := mp.getInputChannelsLastRegLen()
	inputChannelsLastRegKey := mp.getInputChannelsLastRegKey()
	cancel := make(chan interface{})
	for item := range mp.outputMap.IterWithCancel(cancel) {
		k := item.Key.(int64)
		val := item.Value.(*SortedMap)
		if k < inputChannelsLastRegKey {
			inputChannelsLastRegLen = mp.getChannelsCountLogItemLenFromTime(k)
		}
		if val.Len() == inputChannelsLastRegLen {
			e := mp.transformFn(mp, val)
			s := make(chan interface{}, 1)
			s <- e
			close(s)
			mp.stream <- s
			keysToDelete = append(keysToDelete, k)
		} else {
			close(cancel)
			break
		}
	}
	return keysToDelete
}

// cleanUpOutputMap - Deletes the identified keys from the outputMap
func (mp *MsgMultiplexer) cleanUpOutputMap(keysToDelete []int64) {
	for _, k := range keysToDelete {
		mp.lock.Lock()
		mp.outputMap.Delete(k)
		mp.lock.Unlock()
	}
}

// Iter iterates over the items in the Multiplexer
// Each item is sent over a channel, so that
// we can iterate over the it using the builtin range keyword
func (mp *MsgMultiplexer) Iter() chan interface{} {
	return Bridge(mp.doneHandler.Done(), mp.stream)
}

// bridgeIter iterates over the items of the preStream channel in the Multiplexer
// Each item is sent over a channel, so that
// we can iterate over it using the builtin range keyword
func (mp *MsgMultiplexer) bridgeIter() <-chan interface{} {
	return Bridge(mp.doneHandler.Done(), mp.preStream)
}
