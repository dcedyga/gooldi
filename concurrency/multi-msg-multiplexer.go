package concurrency

import (
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

//MultiMsgMultiplexer struct - converge multiple inputChannels into one channel. Closing the Listener and done inputChannels need to be managed outside the Multiplexer
type MultiMsgMultiplexer struct {
	id            string
	inputChannels *Map
	doneHandler   *DoneHandler
	outputMap     *SortedMap
	preStream     chan (<-chan interface{})
	stream        chan (<-chan interface{})
	lock          *sync.RWMutex
	sequence      interface{}
	waitForAll    bool
	BufferSize    int
	MsgType       string
	sendPeriod    *time.Duration
}

type MultiMsgResultItem struct {
	key   interface{}
	value interface{}
}

//NewMultiMsgMultiplexer - Constructor
func NewMultiMsgMultiplexer(dh *DoneHandler, msgType string, opts ...func(*MultiMsgMultiplexer)) *MultiMsgMultiplexer {
	id := uuid.NewV4().String()
	mp := &MultiMsgMultiplexer{
		id:            id,
		inputChannels: NewMap(),
		doneHandler:   dh,
		outputMap:     NewSortedMap(),
		preStream:     make(chan (<-chan interface{})),
		stream:        make(chan (<-chan interface{})),
		sequence:      0,
		waitForAll:    false,
		BufferSize:    1,
		MsgType:       msgType,
		sendPeriod:    nil,
		lock:          &sync.RWMutex{},
	}
	for _, opt := range opts {
		opt(mp)
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

// MultiMsgMultiplexerSequence - option to add a sequence value to the MultiMsgMultiplexer
func MultiMsgMultiplexerSequence(seq interface{}) func(*MultiMsgMultiplexer) {
	return func(mp *MultiMsgMultiplexer) {
		mp.sequence = seq
	}
}

// MultiMsgMultiplexerWaitForAll - option to add a waitforall value to the MultiMsgMultiplexer
func MultiMsgMultiplexerWaitForAll(waitforall bool) func(*MultiMsgMultiplexer) {
	return func(mp *MultiMsgMultiplexer) {
		mp.waitForAll = waitforall
	}
}

// MultiMsgMultiplexerBufferSize - option to add a buffersize value to the MultiMsgMultiplexer
func MultiMsgMultiplexerBufferSize(bufferSize int) func(*MultiMsgMultiplexer) {
	return func(mp *MultiMsgMultiplexer) {
		mp.BufferSize = bufferSize
	}
}

// MultiMsgMultiplexerSendPeriod - option to add a send period value to the MultiMsgMultiplexer
func MultiMsgMultiplexerSendPeriod(d *time.Duration) func(*MultiMsgMultiplexer) {
	return func(mp *MultiMsgMultiplexer) {
		mp.sendPeriod = d
	}
}

// ID - retrieves the Id of the MultiMsgMultiplexer
func (mp *MultiMsgMultiplexer) ID() string {
	return mp.id
}

// Sequence - retrieves the Sequence of the MultiMsgMultiplexer
func (mp *MultiMsgMultiplexer) Sequence() interface{} {
	return mp.sequence
}

func (mp *MultiMsgMultiplexer) Set(key interface{}, value chan interface{}) {
	mp.lock.Lock()

	mp.inputChannels.Set(key, value)
	m := NewSortedMap()
	mp.outputMap.Set(key, m)

	mp.lock.Unlock()
	go func(key interface{}) {
		for item := range value {
			select {
			default:
				s := make(chan interface{}, 1)
				i := &MultiMsgResultItem{
					key:   key,
					value: item,
				}
				s <- i
				close(s)
				mp.preStream <- s

			}
		}

	}(key)

}

func (mp *MultiMsgMultiplexer) Get(key interface{}) (chan interface{}, bool) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	v, ok := mp.inputChannels.Get(key)
	if !ok {
		return nil, false
	}
	return v.(chan interface{}), true
}

func (mp *MultiMsgMultiplexer) delete(key interface{}) {
	mp.lock.Lock()
	defer mp.lock.Unlock()

	mp.inputChannels.Delete(key)

}
func (mp *MultiMsgMultiplexer) Start() {
	go mp.mainProcess()
	if mp.sendPeriod != nil {
		go mp.sendWithTimer()
	}

}
func (mp *MultiMsgMultiplexer) sendWithTimer() {
	drain := false
loop:
	for !drain {
		select {
		case <-mp.doneHandler.Done():
			drain = true
			continue loop
		case <-time.After(*mp.sendPeriod):
			mp.sendToMainBridge()
		}
	}

}
func (mp *MultiMsgMultiplexer) AddItemToMap(event *Event, m *SortedMap) {

	//Check length
	if m.Len() == mp.BufferSize {
		//delete first Item
		keyToDelete, okd := m.GetKeyByIndex(0)
		if okd {
			m.Delete(keyToDelete)
		}
	}
	//Add new item
	mp.lock.Lock()
	m.Set(event.InitMessage.TimeInNano, event)
	mp.lock.Unlock()
}
func (mp *MultiMsgMultiplexer) sendToMainBridge() {
	e := mp.transformSortedMapToEvent(mp.outputMap)
	s := make(chan interface{}, 1)
	s <- e
	close(s)
	mp.stream <- s
}
func (mp *MultiMsgMultiplexer) mainProcess() {
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
				key := v.(*MultiMsgResultItem).key
				event := v.(*MultiMsgResultItem).value.(*Event)
				rm, ok := mp.outputMap.Get(key)
				if ok {
					m := rm.(*SortedMap)
					mp.AddItemToMap(event, m)
					//If waitForAll - we just wait till all the maps have at least one entry
					if mp.waitForAll {
						if !mp.allReady() {
							nextItem = true
							continue loop
						}
					}
					if mp.sendPeriod == nil {
						mp.sendToMainBridge()
					}
					nextItem = true
					continue loop
				}

			}
		}

	}
}

func (mp *MultiMsgMultiplexer) allReady() bool {
	allReady := true
	for item := range mp.outputMap.Iter() {
		i := item.Value.(*SortedMap)
		if i.Len() == 0 {
			allReady = false

		}
	}
	return allReady
}

func (mp *MultiMsgMultiplexer) close() {

	fmt.Printf("MultiMsgMultiplexer - is closed\n")

}
func (mp *MultiMsgMultiplexer) transformSortedMapToEvent(r *SortedMap) *Event {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	m := NewSortedMap()
	for item := range r.Iter() {
		it := NewSortedMap()
		for i := range item.Value.(*SortedMap).Iter() {
			it.Set(i.Key, i.Value)
		}
		m.Set(item.Key, it)
	}
	//Init message
	initmsg := &Message{
		ID:         uuid.NewV4().String(),
		Message:    "MultiMsgMultiplexer message",
		TimeInNano: time.Now().UnixNano(),
		MsgType:    mp.MsgType,
	}
	//TODO - InMessagesequenceMessage

	rmsg := &Message{
		ID:         uuid.NewV4().String(),
		Message:    m,
		TimeInNano: initmsg.TimeInNano,
		MsgType:    initmsg.MsgType,
	}

	re := &Event{
		InitMessage:       initmsg,
		InMessageSequence: Slice{},
		OutMessage:        rmsg,
		Sequence:          mp.sequence,
	}
	return re
}

// Iter iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
//<-chan interface{}
func (mp *MultiMsgMultiplexer) Iter() chan interface{} {
	return Bridge(mp.doneHandler.Done(), mp.stream)
}

func (mp *MultiMsgMultiplexer) bridgeIter() <-chan interface{} {
	return Bridge(mp.doneHandler.Done(), mp.preStream)
}
