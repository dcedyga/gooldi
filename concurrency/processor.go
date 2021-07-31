package concurrency

import (
	"sync"

	uuid "github.com/satori/go.uuid"
)

// ProcessorOption - option to initialize the processor
type ProcessorOption func(*Processor)

// Processor - Unit that listen to an input channel (inputChan) and process work.
// Closing the inputChan channel needs to be managed outside the Processor using a DoneHandler
// It has a DoneHandler to manage the lifecycle of the processor, an index to determine the
// order in which the processor output results might be stored in a multiplexed pattern, an id
// of the processor, the name of the processor, the state of the processor and an output channel
// that emits the processed results for consumption.
type Processor struct {
	Name          string
	doneHandler   *DoneHandler
	inputChan     chan interface{}
	index         int64
	id            string
	outputChannel chan interface{}
	lock          *sync.RWMutex
	state         State
	transformFn   func(p *Processor, input interface{}, result interface{}) interface{}
}

// NewProcessor - Constructor
func NewProcessor(name string, dh *DoneHandler, opts ...ProcessorOption) *Processor {
	id := uuid.NewV4().String()
	p := &Processor{
		id:            id,
		doneHandler:   dh,
		Name:          name,
		inputChan:     nil,
		outputChannel: make(chan interface{}),
		index:         0,
		state:         Init,
		lock:          &sync.RWMutex{},
	}
	p.transformFn = defaultProcessorTransformFn
	for _, opt := range opts {
		opt(p)
	}
	go p.doneRn()
	return p
}

// doneRn - Checks when the Processor is done by listening to the closure of the DoneHandler.Done channel
func (p *Processor) doneRn() {
	select {
	case <-p.doneHandler.Done():
		p.close()
	}
}

// ProcessorTransformFn - option to add a function to transform the output into
// the desired output structure to the Processor
func ProcessorTransformFn(fn func(pr *Processor, input interface{}, result interface{}) interface{}) ProcessorOption {
	return func(p *Processor) {
		p.transformFn = fn
	}
}

// ProcessorWithIndex - option to add a index value to the processor to define its order of execution
func ProcessorWithIndex(idx int64) ProcessorOption {
	return func(p *Processor) {
		p.index = idx
	}
}

// ProcessorWithInputChannel - option to add an inputchannel to the processor
func ProcessorWithInputChannel(in chan interface{}) ProcessorOption {
	return func(p *Processor) {
		p.inputChan = in
	}
}

// ID - retrieves the Id of the Processor
func (p *Processor) ID() string {
	return p.id
}

// Index - retrieves the Index of the Processor
func (p *Processor) Index() int64 {
	return p.index
}

// InputChannel - retrieves the InputChannel of the Processor
func (p *Processor) InputChannel() chan interface{} {
	return p.inputChan
}

// OutputChannel - retrieves the OutputChannel of the Processor
func (p *Processor) OutputChannel() chan interface{} {
	return p.outputChannel
}

// Process - When the processor is in Processing state processes a defined function. When the Processor is
// in stop state the processor will still consume messages from the input channel but it will produce a nil
// output as no process will be involved.
func (p *Processor) Process(f func(pr *Processor, input interface{}, params ...interface{}) interface{}, params ...interface{}) {
	p.setState(Processing)
	drain := false
	for input := range p.inputChan {
		var result interface{}
		if p.GetState() == Processing {
			result = f(p, input, params...)
		}
		r := p.transformFn(p, input, result)
		nextItem := false
	loop:
		for !nextItem && !drain {
			select {
			case <-p.doneHandler.Done():
				nextItem = true
				drain = true
				continue loop
			default:
				select {
				case p.outputChannel <- r:
					nextItem = true
					continue loop
				}
			}
		}
	}
	CloseChannel(p.outputChannel)
}

// Stop - stops the processor.
func (p *Processor) Stop() {
	p.lock.Lock()
	p.state = Stopped
	p.inputChan = nil
	p.lock.Unlock()

}

// Start - starts the processor.
func (p *Processor) Start() {
	p.lock.Lock()
	p.state = Processing
	p.lock.Unlock()

}

// GetState - retrieves the state of the Processor
func (p *Processor) GetState() State {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.state
}

// setState - retrieves the state of the Processor
func (p *Processor) setState(val State) {
	p.lock.Lock()
	p.state = val
	p.lock.Unlock()
}

// HasValidInputChan - checks if the input channel is valid and not nil.
func (p *Processor) HasValidInputChan() bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.inputChan != nil
}

//close - Closes the processor
func (p *Processor) close() {
	//fmt.Printf("Processor %v: Closing\n", p.Name)
	p.lock.Lock()
	p.inputChan = nil
	p.lock.Unlock()

}

/******************************************************************************
Plug and Play and Transformation functions
*******************************************************************************/

// ProcessorMessagePairTransformFn - Gets the processor, input Message and output message and returns the processed
// output in the form of an MessagePair
func ProcessorMessagePairTransformFn(p *Processor, input interface{}, output interface{}) interface{} {

	if IsMessagePtr(input) && IsMessagePtr(output) {
		return NewMessagePair(input.(*Message),
			output.(*Message),
			MessagePairWithIndex(p.index),
			MessagePairWithCorrelationKey(input.(*Message).CorrelationKey),
		)
	}
	out := NewMessage(nil, input.(*Message).MsgType, MessageWithIndex(p.index), MessageWithCorrelationKey(input.(*Message).CorrelationKey))
	return &MessagePair{
		In:             input.(*Message),
		Out:            out,
		CorrelationKey: input.(*Message).CorrelationKey,
	}
}

// defaultProcessorTransformFn - Gets the processor, input and result and outputs the result
func defaultProcessorTransformFn(p *Processor, input interface{}, output interface{}) interface{} {
	if output != nil {
		output.(*Message).Index = p.index
		return output
	}
	return NewMessage(nil,
		input.(*Message).MsgType,
		MessageWithIndex(p.index),
		MessageWithCorrelationKey(input.(*Message).CorrelationKey),
	)

}
