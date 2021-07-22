package concurrency

import (
	"fmt"
	"sync"
	uuid "github.com/satori/go.uuid"
)

// FilterOption - option to initialize the filter
type FilterOption func(*Filter)

// Filter - Unit that listen to an input channel (inputChan) and filter work.
// Closing the inputChan channel needs to be managed outside the Filter using a DoneHandler
// It has a DoneHandler to manage the lifecycle of the filter, an index to determine the
// order in which the filter might be run, an id of the filter,
// the name of the filter, the state of the filter and an output channel
// that emits the processed results for consumption.
type Filter struct {
	Name          string
	doneHandler   *DoneHandler
	inputChan     chan interface{}
	index      	  int64
	id            string
	outputChannel chan interface{}
	lock          *sync.RWMutex
	transformFn   func(f *Filter, input interface{}) interface{}
	state         State
}

// NewFilter - Constructor
func NewFilter(name string, dh *DoneHandler, opts ...FilterOption) *Filter {
	id := uuid.NewV4().String()
	f := &Filter{
		id:            id,
		doneHandler:   dh,
		Name:          name,
		inputChan:     nil,
		outputChannel: make(chan interface{}),
		index:      0,
		state:         Init,
		lock:          &sync.RWMutex{},
	}
	f.transformFn = defaultFilterTransformFn
	for _, opt := range opts {
		opt(f)
	}
	go f.doneRn()
	return f
}

// doneRn - Checks when the Filter is done by listening to the closure of the DoneHandler.Done channel
func (f *Filter) doneRn() {
	select {
	case <-f.doneHandler.Done():
		f.close()
	}
}

// FilterWithSequence - option to add a index value to the filter to define its order of execution
func FilterWithIndex(idx int64) FilterOption {
	return func(f *Filter) {
		f.index = idx
	}
}

// FilterWithInputChannel - option to add an inputchannel to the filter
func FilterWithInputChannel(in chan interface{}) FilterOption {
	return func(f *Filter) {
		f.inputChan = in
	}
}

// FilterTransformFn - option to add a function to transform the output into
// the desired output structure to the Filter
func FilterTransformFn(fn func(fr *Filter, input interface{}) interface{}) FilterOption {
	return func(f *Filter) {
		f.transformFn = fn
	}
}

// ID - retrieves the Id of the Filter
func (f *Filter) ID() string {
	return f.id
}

// Sequence - retrieves the Index of the Filter
func (f *Filter) Index() interface{} {
	return f.index
}

// InputChannel - retrieves the InputChannel of the Filter
func (f *Filter) InputChannel() chan interface{} {
	return f.inputChan
}

// OutputChannel - retrieves the OutputChannel of the Filter
func (f *Filter) OutputChannel() chan interface{} {
	return f.outputChannel
}

// Filter - When the filter is in Processing state filters a defined function. When the Filter is
// in stop state the filter will still consume messages from the input channel and it will
// output the input Message as no filter will be involved.
func (f *Filter) Filter(fn func(f *Filter, input interface{}, params ...interface{}) bool, params ...interface{}) {
	f.setState(Processing)
	drain := false
	for input := range f.inputChan {
		result := true
		if f.GetState() == Processing {
			result = fn(f, input, params...)
		}
		if result {
			r := f.transformFn(f, input)
			nextItem := false
		loop:
			for !nextItem && !drain {
				select {
				case <-f.doneHandler.Done():
					nextItem = true
					drain = true
					continue loop
				default:
					select {
					case f.outputChannel <- r:
						nextItem = true
						continue loop
					}
				}
			}
		}
	}
	CloseChannel(f.outputChannel)
}

// Stop - stops the filter.
func (f *Filter) Stop() {
	f.lock.Lock()
	f.state = Stopped
	f.inputChan = nil
	f.lock.Unlock()

}

// Start - starts the filter.
func (f *Filter) Start() {
	f.lock.Lock()
	f.state = Processing
	f.lock.Unlock()

}

// GetState - retrieves the state of the Filter
func (f *Filter) GetState() State {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.state
}

// setState - retrieves the state of the Filter
func (f *Filter) setState(val State) {
	f.lock.Lock()
	f.state = val
	f.lock.Unlock()
}

// HasValidInputChan - checks if the input channel is valid and not nil.
func (f *Filter) HasValidInputChan() bool {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.inputChan != nil
}

//close - Closes the filter
func (f *Filter) close() {
	fmt.Printf("Filter %v: Closing\n", f.Name)
	f.lock.Lock()
	f.inputChan = nil
	f.lock.Unlock()

}

/******************************************************************************
Plug and Play and Transformation functions 
*******************************************************************************/

// defaultFilterTransformFn - Gets the filter and input outputs the input with the 
// filter index, if the input has the Index property
func defaultFilterTransformFn(f *Filter, input interface{}) interface{} {
	if (HasField(input,"Index")){
		SetFieldInt64Val("Index",input,f.index)
	}
	return input
}
