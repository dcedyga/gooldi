![header](https://capsule-render.vercel.app/api?type=waving&color=gradient&height=300&section=header&text=gooldi&fontSize=90&animation=fadeIn&fontAlignY=25&desc=go%20concurrency%20library%20for%20deterministic%20and%20non-%20deterministic%20stream%20processing&descAlignY=51&descAlign=50)


<a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> was born with the aim to use golang concurrency capabilities to provide a set of streaming patterns and approaches that allow to build very complex flows/ pipelines to fulfil the main paradigms for deterministic and no-deterministic stream processing.

<a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> brings an implementation of most of the concurrency patterns define in ["Concurrency in Go"](https://katherine.cox-buday.com/concurrency-in-go/) by Cox-Buday. And a set of generators and utilities to ease working with these concurrency patterns.

<a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> provides a thread-safe implementation of Map, Slice, SortedMap and SortedSlice to access to the relevant maps and slice types shared across goroutines without race conditions.

<a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> main purpose is to provide deterministic and non-deterministic stream processing capabilities and to do so, the following entities have been implemented: Message, MessagePair, Bcaster, Processor, Filter,MessageMultiplexer and MultiMessageMultiplexer. With these entities complex pipeline and flow architectures can be implemented with throughputs that are close to 500k processed messages per second in a macbook pro with 32G of RAM and 2,9 GHz Intel Core i9 processor.

## Concurrency in Go
```
Concurrency is about dealing with lots of things at once. It is a way to structure software, 
particularly as a way to write clean code that interacts well with the real world.

— Rob Pike
```
Concurrency refers to the design of the system, while parallelism relates to the execution. Concurrent programming is one of the most interesting aspects of the Go language. Go is designed with concurrency in mind and allows us to build complex concurrent pipelines. Go concurrency main building blocks are around parallel composition of sequential processes and communication between these processes. Go’s approach to concurrency can be best phrased by:
```
Do not communicate by sharing memory; instead, share memory by communicating. 
```
### Go Concurrency vs Multithreading

In any other mainstream programming language, when concurrent threads need to share data in order to communicate, a lock is applied to the piece of memory. Instead of applying a lock on the shared variable, Go allows you to communicate (or send) the value stored in that variable from one thread to another. The default behavior is that both the thread sending the data and the one receiving the data will wait till the value reaches its destination. The “waiting” of the threads forces proper synchronization between threads when data is being exchanged.

Some facts:
 - [x] goroutines are managed by go runtime and has no hardware dependencies while OS threads are managed by kernal and has hardware dependencies
- [x] goroutines are smaller: typically 2KB of stack size, threads 1-2MB
- [x] Stack size of go is managed in run-time and can grow up to 1GB which is possible by allocating and freeing heap storage while for threads, stack size needs to be determined at compile time
- [x] goroutine use channels to communicate with other goroutines with low latency. There is no easy communication medium between threads and huge latency between inter-thread communication
- [x] goroutine do not have any identity while threads do (TID)
- [x] goroutines are created and destoryed by the go's runtime. These operations are very cheap compared to threads as go runtime already maintain pool of threads for goroutines. In this case OS is not aware of goroutines
- [x] goroutines are coopertively scheduled,  when a goroutine switch occurs, only 3 registers need to be saved or restored. Threads are preemptively scheduled, switching cost between threads is high as scheduler needs to save/restore more than 50 registers and states. This can be quite significant when there is rapid switching between threads.

### Interesting Reads and References

* [Concurrency is not Parallelism by Rob Pike](https://www.youtube.com/watch?v=oV9rvDllKEg)
* [Concurrency in Go by Thejas Babu](https://medium.com/@thejasbabu/concurrency-in-go-e4a61ec96491)
* [Achieving concurrency in Go by Uday Hiwarale](https://medium.com/rungo/achieving-concurrency-in-go-3f84cbf870ca)
* [Go's work-stealing scheduler](https://rakyll.org/scheduler/)

## gooldi: Concurrency patterns
Important concurrency patterns to highlight are:

- <a href="./concurrency/bridge.go#L01"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=300&section=header&text=Bridge&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> is a way to present a single-channel facade over a channel of channels. It is used to consume values from a sequence of channels (channel of channels) doing an ordered write from different sources. By bridging the channels it destructures the channel of channels into a simple channel, allowing to multiplex the input and simplify the consumption.With this pattern we can use the channel of channels from within a single range statement and focus on our loop’s logic.
- <a href="./concurrency/fan-in.go#L01"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=300&section=header&text=FanIn&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> combines multiple results in the form of an slice of channels into one channel. This implementation uses a WaitGroup in order to multiplex all the results of the slice of channels. The output is not produced in sequence. This pattern is good for  Non-deterministic stream processing when order is not important and deterministic outcomes are not required.
- <a href="./concurrency/fan-in.go#L40"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=300&section=header&text=FanInRec&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> combines multiple results in the form of an slice of channels into one channel. This implementation uses a a recursive approach in order to multiplex all the results of the slice of channels. The output is not produced in sequence. This pattern is good for  Non-deterministic stream processing when order is not important and deterministic outcomes are not required.
- <a href="./concurrency/or-done.go#L01"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=300&section=header&text=OrDone&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a>  Wraps a channel with a select statement that also selects from a done channel. Allows to cancel the channel avoiding go-routine leaks.
- <a href="./concurrency/or.go#L01"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=300&section=header&text=Or&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> returns the value of the fastest channel.
- <a href="./concurrency/route.go#L01"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=300&section=header&text=Route&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> Representation of the tee pattern. Takes a single input channel and an arbitrary number of output channels and duplicates each input into every output. When the input channel is closed, all outputs channels are closed. It allows to route or split an input into multiple outputs.

### Utilities and Generators

- <a href="./concurrency/as-chan.go#L01">AsChan</a>  sends the contents of a slice through a channel
- <a href="./concurrency/close-chan.go#L01">CloseChannel</a> Checks if the channel is not closed and closes it
- <a href="./concurrency/take.go#L01">Take</a> Takes a defined number of values by num from a channel
- <a href="./concurrency/to-string.go#L01">ToString</a>  Converts any type of channel into a string channel
- <a href="./concurrency/repeat.go#L01">Repeat</a> Generator that repeats the values defined on the values slice indefinitely.

## gooldi: Thread-safe Maps and Slices
```
Slice and Map - cannot be used safely from multiple goroutines without the risk of having 
a race condition.
```
<a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> provides and implementation of Slice and Map types which can be safely shared between multiple goroutines by protecting the access to the shared data by a mutex. 

What a mutex does is basically to acquire a lock when it needs to access our concurrent type, When holding the lock a goroutine can read and/or write to the shared data protected by the mutex safely, this lock can be acquired by a single goroutine at a time, and if other goroutine needs access to the same shared data it waits until the lock has been released by the goroutine holding the lock. <a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi's&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> implementation of these types manages the acquisition and release of locks avoiding deadlocks and unwanted behaviours.

<a href="./concurrency/slice.go#L01"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=300&section=header&text=Slice&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> and <a href="./concurrency/sorted-slice.go#L01"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=200&section=header&text=SortedSlice&fontSize=100&animation=fadeIn&fontAlignY=55" width="100" height="23"/></a> are `slice` types and 
<a href="./concurrency/map.go#L01"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=300&section=header&text=Map&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> and <a href="./concurrency/sorted-map.go#L01"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=200&section=header&text=SortedMap&fontSize=100&animation=fadeIn&fontAlignY=55" width="100" height="23"/></a> are `map` types within <a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> that are safe to share among threads. 
`SortedSlice` is an ordered implementation of the `Slice` type and `SortedMap` is an ordered implementation of the `Map` type. The performance among sorted vs non-sorted implementations of these types are very similar as `SortedSlice` and `SortedMap` have a `dirty` property that is set when there is a change on the `SortedSlice` or the `SortedMap`. The sort mechanism will only happen when attempting to read the relevant type and if the `dirty`property is set to `true`, keeping the `SortedSlice` and the `SortedMap` very performant. 

```go
//Sample SortedSlice
s := NewSortedSlice()
cancel := make(chan interface{})
s.Append("Juan")
s.Append("Pepi")
s.Append("David")
s.Append("Ayleen")
s.Append("Lila")
s.Append("Freddy")
s.Append("Moncho")
s.Append("Zac")
s.Append("Caty")
s.Append("Tom")
for item := range s.IterWithCancel(cancel) {
    fmt.Printf("%v:%v\n", item.Index, item.Value)
    if item.Index == 5 {
        close(cancel)
        break
    }
}
```
  
```go
//Sample SortedMap
m := NewSortedMap()
cancel := make(chan interface{})
for i := 0; i < 10; i++ {
    m.Set(i, i)
}
for item := range m.IterWithCancel(cancel) {
    fmt.Printf("%v:%v\n", item.Key, item.Value)
}
```
## gooldi: Handling channel cancellation

<a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> provides an ordered layered mechanism in order to handle cancellation of channels. While we explored the use of the go `context`package to use cancelable contexts, the way it deals with the cancellation is not fit for purpose for the processing mechanisms that are part of the <a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> library. Cancellation within the `context` package happens in a non-deterministic way and it can lead to exceptions due to the closure of channels that are still active.

<a href="./concurrency/done-manager.go#L01"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=200&section=header&text=DoneManager&fontSize=100&animation=fadeIn&fontAlignY=55" width="100" height="23"/></a> and <a href="./concurrency/done-handler.go#L01"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=200&section=header&text=DoneHandler&fontSize=100&animation=fadeIn&fontAlignY=55" width="100" height="23"/></a> are the two types that manage channel cancellation in <a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a>. `DoneManager` registers `DoneHandlers` and keeps track of them. It has the capability to organize the `DoneHandlers`in a layered `Map` and when cancelling the `DoneManager`loops through this map by layers. It has also the capability to add a delay between layer cancellation, allowing for a graceful shutdown.
In `DoneManager` a deadline and a timeout can be set, so cancellation can be triggered when these are reached.
In `DoneHandler` we can setup a deadline that will cancel that specific handler and notify the `DoneManager` for proper de-registration of the `DoneHandler`.

```go
//DoneManager and DoneHandler Sample
dm := NewDoneManager(
    DoneManagerWithDelay(1*time.Millisecond),
    DoneManagerWithTimeout(100*time.Millisecond),
)
dh := dm.AddNewDoneHandler(0)
dh1 := dm.AddNewDoneHandler(1)
go func() {
    for {
        select {
        case <-dh.Done():
            fmt.Printf(
                "We are done with delay and timeout - on layer 0 -> %v: %v\n",
                dh.ID(),
                dh.Err(),
            )
            return
        }
    }
}()
go func() {
    for {
        select {
        case <-dh1.Done():
            fmt.Printf(
                "We are done with delay and timeout - on layer 1 -> %v: %v\n",
                dh1.ID(),
                dh1.Err(),
            )
            return
        }
    }
}()
go func() {
    for {
        select {
        case <-dm.Done():
            fmt.Printf(
                "We are done with delay and timeout - on DoneManager -> %v: %v\n",
                dm.ID(),
                dm.Err(),
            )
            return
        }
    }
}()

time.Sleep(1000 * time.Millisecond)

```

## gooldi: Stream Processing Entities

The foundation of <a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi's&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> stream processing is based on the following concepts. We have:
- <a href="./concurrency/message.go#L15"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=200&section=header&text=Message&fontSize=100&animation=fadeIn&fontAlignY=55" width="100" height="23"/></a> and <a href="./concurrency/message.go#L43"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=200&section=header&text=MessagePair&fontSize=100&animation=fadeIn&fontAlignY=55" width="100" height="23"/></a> entities that are shared across the pipeline and acts as the interchangeable entity within the flow of the process. These entities are the representation of the Stream as a Stream of Messages or MessagePairs. 
- <a href="./concurrency/bcaster.go#L13"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=200&section=header&text=BCaster&fontSize=100&animation=fadeIn&fontAlignY=55" width="100" height="23"/></a> is a broadcaster that allows to broadcast any type of message to its listeners. You can register a listener with the `AddListener`method that returns a channel of type interface{}. <a href="./concurrency/bcaster.go#L13"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=200&section=header&text=BCaster&fontSize=100&animation=fadeIn&fontAlignY=55" width="100" height="23"/></a> uses a transformation function, that is highly customizable (you can use your own), to map the input message to an output structure. the default transform function called `defaultBCasterTransformFn` Gets the bcaster, input and outputs the input with no changes.
- 

#### *Message and MessagePair*

```go
/*
Message 
*/
type Message struct {
	ID         string
	Message    interface{}
	TimeInNano int64
	MsgType    string
	CorrelationKey int64
	Index int64
}

/*
MessagePair
*/
type MessagePair struct {
	In             *Message
	Out            *Message
	Index          int64
	CorrelationKey int64
}

```
<a href="./concurrency/message.go#L15"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=200&section=header&text=Message&fontSize=100&animation=fadeIn&fontAlignY=55" width="100" height="23"/></a>  is a Struct that represents an message in the context of the concurrency package. With the following properties:
 - **ID**: the ID of the message
 - **Message**: the Message,
 - **TimeInNano**: the time that was produced
 - **MsgType**: the type of the message,
 - **CorrelationKey**: a correlation key to correlate to other messages
 - **Index**: the index related to the order that can be used to produce deterministic outputs

<a href="./concurrency/message.go#L43"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=200&section=header&text=MessagePair&fontSize=100&animation=fadeIn&fontAlignY=55" width="100" height="23"/></a> is a Struct that represents an input/output message Pair in the context of the concurrency package. With the following properties:
- **In**: represents an Input Message,
- **Out**: represents an Output Message,
- **Index**: a index number that is usefull to define the processing order of the MessagePair 
- **CorrelationKey**: a correlationKey to correlate to other messages

`NewMessage`and `NewMessagePair`are constructors that simplify the creation of these entities.
```go
//Message constructor
NewMessage(fmt.Sprintf("This is message 1"),
    "string",
    MessageWithCorrelationKey(1),
    MessageWithIndex(1),
)
```

#### *BCaster*

#### *Processor and Filter*

#### *MsgMultiplexer*

#### *MultiMsgMultiplexer*

## gooldi: Deterministic and Non-Deterministic Stream Processing
## gooldi:Highly customizable

<a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> is highly customizable. For example one can define its own Message entity to represent better the Stream that it wants to process. With its own `BCaster`that has a specific transformation/mapping function of the Broadcast output.

#### BCaster and processor customization
As an example we could have a `CustomEvent` with an InitMessage and OutMessage properties, the BCaster has a custom transform function called `BCasterCustomEventTransformFn` which transforms an standard input `Message` into a `CustomEvent`that has the InitMessage and OutMessage properties filled with the input `Message`. The `BCaster` broadcasts the `CustomEvent`to all the registered listeners, in this case one `Processor`that uses it process function to square the Message property(payload) of the OutMessage of the incoming `CustomEvent` producing a new standard `Message` which is transformed by the `ProcessorCustomEventTransformFn` to produce a new `CustomEvent` with the InitMessage of the incomming `CustomEvent` as InitMessage property and the result `Message`of the process function as OutMessage property. This example just gives us an idea of how flexible <a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> is.

**Process functions**
* The `Processor` process functions follow this primitive:
```go 
func func(pr *Processor, input interface{}, params ...interface{}) interface{} {}
```

**Mapping/Transformation functions**
* The `BCaster` transformation functions follow this primitive:
```go 
func func(b *BCaster, input interface{}) interface{}{
    //Code
}
```
* The `Processor` transformation functions follow this primitive:
```go
func func(p *Processor, input interface{}, result interface{}) interface{} {
    //Code
}
```
Please find below the entire code snippet:
```go
//BCaster and processor customization

//CustomEvent
type CustomEvent struct {
	InitMessage *Message
	OutMessage  *Message
}

//BCasterCustomEventTransformFn
func BCasterCustomEventTransformFn(b *BCaster, 
    input interface{}) interface{} {
	var msg *Message
	if input != nil {
		msg = input.(*Message)
	}
	e := &CustomEvent{
		InitMessage: msg,
		OutMessage:  msg,
	}
	return e
}

//customProcess3
func customProcess3(pr *Processor, 
    input interface{}, 
    params ...interface{}) interface{} {
	event := input.(*CustomEvent)
	//id := params[0]
	i := event.OutMessage.Message.(int) * event.OutMessage.Message.(int)
	return NewMessage(i,
		event.OutMessage.MsgType,
		MessageWithIndex(event.OutMessage.Index),
		MessageWithCorrelationKey(event.OutMessage.CorrelationKey),
	)
}
//ProcessorCustomEventTransformFn
func ProcessorCustomEventTransformFn(p *Processor, 
    input interface{}, 
    result interface{}) interface{} {
	var event *CustomEvent
	var resultMsg *Message
	if input != nil {
		event = input.(*CustomEvent)
	}
	if result != nil {
		resultMsg = result.(*Message)
	}
	r := &CustomEvent{
		InitMessage: event.InitMessage,
		OutMessage:  resultMsg,
	}
	return r
}
//main
func main(){
    dm := NewDoneManager(DoneManagerWithDelay(1 * time.Millisecond))
	defer dm.GetDoneFunc()()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	//Create caster
	b := NewBCaster(dh,
		"bcaster",
		BCasterTransformFn(BCasterCustomEventTransformFn),
	)
	//Create Processors
	w2 := NewProcessor(
		fmt.Sprintf("worker%v", 2),
		dh1,
		ProcessorWithIndex(2),
		ProcessorWithInputChannel(b.AddListener(dh1)),
		ProcessorTransformFn(ProcessorCustomEventTransformFn),
	)
	go w2.Process(customProcess3, 2)
	time.Sleep(300 * time.Microsecond)
	// Get all multiplexed messages
	go func() {
		i := 0
		for v := range w2.OutputChannel() {
			_ = v
			i++
		}
		fmt.Printf("Total Messages from Processor %v: %v,\n", w2.ID(), i)
	}()
	// Broadcast Messages
	go func() {
		for msgId := 0; msgId < 300000; msgId++ {
			e := NewMessage(
				msgId,
				b.MsgType,
			)
			b.Broadcast(e)
		}
	}()
	time.Sleep(1000 * time.Millisecond)
} 
```
#### Processor and MsgMultiplexer customization

```go
//Processor and MsgMultiplexer customization


```



