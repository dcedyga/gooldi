package concurrency

import (
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

// MultiMsgMultiplexerOption - option to initialize the MultiMsgMultiplexer
type MultiMsgMultiplexerOption func(*MultiMsgMultiplexer)

// MultiMsgMultiplexer - The default implementation MultiMsgMultiplexer allows to create complex patterns where multiple Broadcasters
// can emit a Message to multiple processors (consumers) that can potentially represent multiple processing systems,
// do the relevant calculation and multiplex the multiple outputs into a single channel for simplified consumption.
// Its main function is to Mulitplex a set of multiple messages that can be concurrently processed and converge the set of
// initial concurrency.Message into a SortedMap ordered by messageType that can be sent on one channel.
// values of the processors grouped by initial concurrency.Message and ordered by index value of each processor.
// Closure of MultiMsgMultiplexer is handle by a concurrency.DoneHandler that allows to control they way a set of go routines
// are closed in order to prevent deadlocks and unwanted behaviour
// MultiMsgMultiplexer outputs the multiplexed result in one channel using the channel bridge pattern.
// MultiMsgMultiplexer has several modes, the first one is to output the structure everytime a BCaster emits a message, giving
// an output of the last received message per BCaster. The second one is by using a timer to specify the sendPeriod, where the
// output represents the state of the last received messages at the specific point of time of the tick of the period.
// MultiMsgMultiplexer has also a waitForAll property that when true will just start emiting an output when the MultiMsgMultiplexer
// has at least received one message of each of the BCasters.
// MultiMsgMultiplexer has also a BufferSize property (default value is 1) where we can send the n number of last messages sent
// by each BCaster.
// MultiMsgMultiplexer default behaviour can be overridden by providing a MultiMsgMultiplexerItemKeyFn to the key of
// the items of a channel within the output SortedMap for a specific MessageType and MultiMsgMultiplexerTransformFn
// allows to transform the output into the desired structure.
type MultiMsgMultiplexer struct {
	id            string
	inputChannels *Map
	doneHandler   *DoneHandler
	outputMap     *SortedMap
	stream        chan (<-chan interface{})
	lock          *sync.RWMutex
	index         int64
	toStringIndex string
	isReady       bool
	waitForAll    bool
	BufferSize    int
	MsgType       string
	sendPeriod    *time.Duration
	getItemKeyFn  func(v interface{}) int64
	transformFn   func(mp *MultiMsgMultiplexer, sm *SortedMap) interface{}
}

// MultiMsgResultItem - The result item to be stored in the output SortedMap
// type MultiMsgResultItem struct {
// 	key   interface{}
// 	value interface{}
// }

//NewMultiMsgMultiplexer - Constructor
func NewMultiMsgMultiplexer(dh *DoneHandler, msgType string, opts ...MultiMsgMultiplexerOption) *MultiMsgMultiplexer {
	id := uuid.NewV4().String()
	mp := &MultiMsgMultiplexer{
		id:            id,
		inputChannels: NewMap(),
		doneHandler:   dh,
		outputMap:     NewSortedMap(),
		stream:        make(chan (<-chan interface{})),
		index:         0,
		toStringIndex: "00000000000000000000",
		isReady:       false,
		waitForAll:    false,
		BufferSize:    1,
		MsgType:       msgType,
		sendPeriod:    nil,
		lock:          &sync.RWMutex{},
	}

	mp.getItemKeyFn = defaultMultiMsgGetItemKey
	mp.transformFn = defaultMultiMsgTransformFn

	for _, opt := range opts {
		opt(mp)
	}
	if mp.sendPeriod != nil {
		go mp.sendWithTimer()
	}
	go mp.doneRn()
	return mp
}

// doneRn - Checks when the MultiMsgMultiplexer is done by listening to the closure of the DoneHandler.Done channel
func (mp *MultiMsgMultiplexer) doneRn() {
	select {
	case <-mp.doneHandler.Done():
		mp.close()
	}
}

// MultiMsgMultiplexerIndex - option to add a Index value to the MultiMsgMultiplexer
func MultiMsgMultiplexerIndex(idx int64) MultiMsgMultiplexerOption {
	return func(mp *MultiMsgMultiplexer) {
		mp.index = idx
		mp.toStringIndex = IndexToString(idx)
	}
}

// MultiMsgMultiplexerWaitForAll - option to add a waitforall value to the MultiMsgMultiplexer
func MultiMsgMultiplexerWaitForAll(waitforall bool) MultiMsgMultiplexerOption {
	return func(mp *MultiMsgMultiplexer) {
		mp.waitForAll = waitforall
	}
}

// MultiMsgMultiplexerBufferSize - option to add a buffersize value to the MultiMsgMultiplexer
func MultiMsgMultiplexerBufferSize(bufferSize int) MultiMsgMultiplexerOption {
	return func(mp *MultiMsgMultiplexer) {
		mp.BufferSize = bufferSize
	}
}

// MultiMsgMultiplexerSendPeriod - option to add a send period value to the MultiMsgMultiplexer
func MultiMsgMultiplexerSendPeriod(d *time.Duration) MultiMsgMultiplexerOption {
	return func(mp *MultiMsgMultiplexer) {
		mp.sendPeriod = d
	}
}

// MultiMsgMultiplexerTransformFn - option to add a function to transform the SortedMap output into
// the desired output structure to the MultiMsgMultiplexer
func MultiMsgMultiplexerTransformFn(fn func(mp *MultiMsgMultiplexer, sm *SortedMap) interface{}) MultiMsgMultiplexerOption {
	return func(mp *MultiMsgMultiplexer) {
		mp.transformFn = fn
	}
}

// MultiMsgMultiplexerItemKeyFn - option to add a function to resolve the set key value
// of an item of the channel to the map of the MultiMsgMultiplexer specific to a message
func MultiMsgMultiplexerItemKeyFn(fn func(v interface{}) int64) MultiMsgMultiplexerOption {
	return func(mp *MultiMsgMultiplexer) {
		mp.getItemKeyFn = fn
	}
}

// ID - retrieves the Id of the MultiMsgMultiplexer
func (mp *MultiMsgMultiplexer) ID() string {
	return mp.id
}

// Index - retrieves the index of the MultiMsgMultiplexer
func (mp *MultiMsgMultiplexer) Index() int64 {
	return mp.index
}

// DoneHandler - retrieves the DoneHandler of the MultiMsgMultiplexer
func (mp *MultiMsgMultiplexer) DoneHandler() *DoneHandler {
	return mp.doneHandler
}

// ToStringIndex - retrieves the ToStringIndex representation of the MsgMultiplexer
func (mp *MultiMsgMultiplexer) ToStringIndex() string {
	return mp.toStringIndex
}

// Set - Registers a channel in the MultiMsgMultiplexer and starts processing it
func (mp *MultiMsgMultiplexer) Set(key interface{}, value chan interface{}) {
	mp.lock.Lock()
	mp.inputChannels.Set(key, value)
	m := NewSortedMap()
	mp.outputMap.Set(key, m)
	mp.lock.Unlock()
	go mp.processChannel(key, value, m)

}

// Get - Retrieves a channel reqistered in the MultiMsgMultiplexer by key
func (mp *MultiMsgMultiplexer) Get(key interface{}) (chan interface{}, bool) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	v, ok := mp.inputChannels.Get(key)
	if !ok {
		return nil, false
	}
	return v.(chan interface{}), true
}

// Delete - Deletes a registered channel from the MsgMultiplexer map of inputChannels
func (mp *MultiMsgMultiplexer) clean(key interface{}) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	mp.inputChannels.Delete(key)
}

// processChannel - Retrieves all the values of an inputChannel using range and when done it deletes it
// from the MsgMultiplexer map of inputChannels.
func (mp *MultiMsgMultiplexer) processChannel(key interface{}, value chan interface{}, m *SortedMap) {
	for item := range value {
		select {
		//case <-mp.doneHandler.Done(): //need to check this
		default:
			cKey := mp.getItemKeyFn(item)
			mp.storeInOutputMap(cKey, item, m)
			if mp.sendPeriod == nil {
				mp.sendWhenReady()
			}
		}
	}
	mp.clean(key)
}

//storeInOutputMap - Adds an item to the SortedMap that is going to be send as part of the output, the SortedMap
// length is defined by the BufferSize property, allowing to retrieve the last n messages for a specific Messagetype.
func (mp *MultiMsgMultiplexer) storeInOutputMap(cKey int64, v interface{}, m *SortedMap) {
	mp.lock.Lock()
	//Check length
	if m.Len() == mp.BufferSize {
		//delete first Item
		itemToDelete, okd := m.GetSortedMapItemByIndex(0)
		if okd {
			m.Delete(itemToDelete.Key)
		}
	}
	//Add new item
	m.Set(cKey, v)
	mp.lock.Unlock()
}

func (mp *MultiMsgMultiplexer) sendWhenReady() {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	if mp.allReady() {
		mp.sendToMainBridge()
	}
}

// sendWithTimer - sends the output with every tick defined by the sendPeriod property.
func (mp *MultiMsgMultiplexer) sendWithTimer() {
	drain := false
loop:
	for !drain {
		select {
		case <-mp.doneHandler.Done():
			drain = true
			continue loop
		case <-time.After(*mp.sendPeriod):
			if mp.allReady() {
				mp.sendToMainBridge()
			}
		}
	}

}

// sendToMainBridge - sends the output SortedMap item to the output stream for consumption, using the bridge pattern.
func (mp *MultiMsgMultiplexer) sendToMainBridge() {
	e := mp.transformFn(mp, mp.outputMap)
	s := make(chan interface{}, 1)
	s <- e
	close(s)
	select {
	case <-mp.doneHandler.Done():
	case mp.stream <- s:
	}

}

// Checks if all the inputChannels have at least sent one item
func (mp *MultiMsgMultiplexer) allReady() bool {
	if !mp.waitForAll {
		return true
	}
	if !mp.isReady {
		allready := true
		for item := range mp.outputMap.Iter() {
			i := item.Value.(*SortedMap)
			if i.Len() == 0 {
				allready = false
			}
		}
		mp.isReady = allready
	}
	return mp.isReady
}

/////////////////////////////////////////////////////////////////////////////////////////////

// close - Closes the MultiMsgMultiplexer
func (mp *MultiMsgMultiplexer) close() {
	fmt.Printf("MultiMsgMultiplexer - is closed\n")
}

// Iter iterates over the items in the MultiMsgMultiplexer
// Each item is sent over a channel, so that
// we can iterate over the it using the builtin range keyword
func (mp *MultiMsgMultiplexer) Iter() chan interface{} {
	return Bridge(mp.doneHandler.Done(), mp.stream)
}

func (mp *MultiMsgMultiplexer) String() string {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	return "Index: " + mp.toStringIndex + "ID: " + mp.ID()
}

/******************************************************************************
Plug and Play and Transformation functions
*******************************************************************************/

// defaultMultiMsgGetItemKey - gets the last registered key when a new channel is added or removed. It
// is used in conjunction with the defaultMultiMsgGetItemKey to extract the length of the output SortedMap.
// Can be overridden for a more generic implementation
func defaultMultiMsgGetItemKey(v interface{}) int64 {
	return v.(*Message).CorrelationKey
}

// defaultMultiMsgTransformFn - Transforms the SortedMap output into a Message for future consumption as part of the
// output channel of the MultiMsgMultiplexer. Can be overridden for a more generic implementation
// Oriented
func defaultMultiMsgTransformFn(mp *MultiMsgMultiplexer, r *SortedMap) interface{} {
	m := NewSortedMap()
	for item := range r.Iter() {
		it := item.Value.(*SortedMap).Clone()
		m.Set(item.Key, it)
	}
	rmsg := NewMessage(m,
		mp.MsgType,
		MessageWithIndex(mp.index),
	)

	return rmsg
}
