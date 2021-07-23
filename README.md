![header](https://capsule-render.vercel.app/api?type=waving&color=gradient&height=300&section=header&text=gooldi&fontSize=90&animation=fadeIn&fontAlignY=25&desc=go%20concurrency%20library%20for%20deterministic%20and%20non-%20deterministic%20stream%20processing&descAlignY=51&descAlign=50)


<a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> was born with the aim to use golang concurrency capabilities to provide a set of streaming patterns and approaches that allow to build very complex flows/ pipelines to fulfil the main paradigms for deterministic and no-deterministic stream processing.

<a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> brings an implementation of most of the concurrency patterns define in ["Concurrency in Go"](https://katherine.cox-buday.com/concurrency-in-go/) by Cox-Buday. And a set of generators and utilities to ease working with these concurrency patterns.

<a href="https://github.com/dcedyga/gooldi"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=ff9933&fontColor=ffffff&height=300&section=header&text=gooldi&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> provides a thread-safe implementation of Map, Slice, SortedMap and SortedSlice to access to the relevant maps and slice types shared across goroutines without race conditions.

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


- <a href="./concurrency/route.go#L01"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=300&section=header&text=Slice&fontSize=160&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> and <a href="./concurrency/route.go#L01"><img align="center" src="https://capsule-render.vercel.app/api?type=soft&color=6699ff&fontColor=ffffff&height=300&section=header&text=SortedSlice&fontSize=135&animation=fadeIn&fontAlignY=55" width="70" height="23"/></a> are the types within gooldi that are safe to share among threads. The only difference between them is that the SortedSlice is an ordered implementation of the Slice type. The performance among both is very similar as SortedSlice has a `dirty` property that is set when there is a change on the SortedSlice. The sort of the SortedSlice will only happen when attempting to read the SortedSlice and if the `dirty`property is set to `true`, keeping the SortedSlice very performant. 
Both of them can be safely shared between goroutines and the have their relevant constructors `NewSlice` and `NewSortedSlice`. Both have the following methods:
  * **Append** adds an item to the concurrent slice
  * **RemoveItemAtIndex** removes the item at the specified index
  * **IndexOf** returns the index of a specific item
  * **GetItemAtIndex** - Get item at index
  * **Iter** iterates over the items in the concurrent slice, Each item is sent over a channel, so that we can iterate over the slice using the builtin range keyword
  * **IterWithCancel** same as Iter but allows to pass a cancel chan to make the iteration cancelable
  * **Len** - length of the slice
  * **Cap** - capacity of the slice


```go
    s := concurrency.NewSortedSlice()
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

What it does

- Done Manager and Done Handler
- gooldi Stream Processing Entities: Message, MessagePair, Broadcast, Processor, Filter,MessageMultiplexer and MultiMessageMultiplexer
- Deterministic Stream Processing
- Non-Deterministic Stream Processing
- Highly customizable

How it does it

