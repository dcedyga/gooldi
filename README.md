
<div class="text-white bg-blue mb-2">
  .text-white on .bg-blue
</div>
<div class="bg-blue-light mb-2">
  .text-gray-dark on .bg-blue-light
</div>
<div class="text-white bg-red mb-2">
  .text-white on .bg-red
</div>
<div class="text-red bg-red-light mb-2">
  .text-red on .bg-red-light
</div>
<div class="bg-green-light mb-2">
  .text-gray-dark on .bg-green-light
</div>
<div class="bg-yellow mb-2">
  .text-gray-dark on .bg-yellow
</div>
<div class="bg-yellow-light mb-2">
  .text-gray-dark on .bg-yellow-light
</div>
<div class="text-white bg-purple mb-2">
  .text-white on .bg-purple
</div>
<div class="text-white bg-gray-dark mb-2">
  .text-white on .bg-gray-dark
</div>
<div class="bg-gray">
  .text-gray-dark on .bg-gray
</div>

# gooldi


<font color="orange">gooldi</font> was born with the aim to use golang concurrency capabilities to provide a set of streaming patterns and approaches that allow to build very complex flows/ pipelines to fulfil the main paradigms for deterministic and no-deterministic stream processing.

<span style="color:orange;">gooldi</span> brings an implementation of most of the concurrency patterns define in ["Concurrency in Go"](https://katherine.cox-buday.com/concurrency-in-go/) by Cox-Buday. And a set of generators and utilities to ease working with these concurrency patterns.

<span style="color:orange;">gooldi</span>  provides a thread-safe implementation of Map, Slice, SortedMap and SortedSlice to access to the relevant maps and slice types shared across goroutines without race conditions.

### Concurrency in Go
“Concurrency is about dealing with lots of things at once. It is a way to structure software, particularly as a way to write clean code that interacts well with the real world.” — Rob Pike

Concurrency refers to the design of the system, while parallelism relates to the execution. Concurrent programming is one of the most interesting aspects of the Go language. Go is designed with concurrency in mind and allows us to build complex concurrent pipelines. Go concurrency main building blocks are around parallel composition of sequential processes and communication between these processes. Go’s approach to concurrency can be best phrased by:

"Do not communicate by sharing memory; instead, share memory by communicating." 

##### Go Concurrency vs Multithreading

In any other mainstream programming language, when concurrent threads need to share data in order to communicate, a lock is applied to the piece of memory. Instead of applying a lock on the shared variable, Go allows you to communicate (or send) the value stored in that variable from one thread to another. The default behavior is that both the thread sending the data and the one receiving the data will wait till the value reaches its destination. The “waiting” of the threads forces proper synchronization between threads when data is being exchanged.

Some facts:
 * goroutines are managed by go runtime and has no hardware dependencies while OS threads are managed by kernal and has hardware dependencies
 * goroutines are smaller: typically 2KB of stack size, threads 1-2MB
 * Stack size of go is managed in run-time and can grow up to 1GB which is possible by allocating and freeing heap storage while for threads stack size needs to be determined at compile time
 * goroutine use channels to communicate with other goroutines with low latency. There is no easy communication medium between threads and huge latency between inter-thread communication
 * goroutine do not have any identity while threads do (TID)
 * goroutines are created and destoryed by the go's runtime. These operations are very cheap compared to threads as go runtime already maintain pool of threads for goroutines. In this case OS is not aware of goroutines
 * goroutines are coopertively scheduled,  when a goroutine switch occurs, only 3 registers need to be saved or restored. Threads are preemptively scheduled, switching cost between threads is high as scheduler needs to save/restore more than 50 registers and states. This can be quite significant when there is rapid switching between threads.

##### Interesting Reads and References

* [Concurrency is not Parallelism by Rob Pike](https://www.youtube.com/watch?v=oV9rvDllKEg)
* [Concurrency in Go by Thejas Babu](https://medium.com/@thejasbabu/concurrency-in-go-e4a61ec96491)
* [Achieving concurrency in Go by Uday Hiwarale](https://medium.com/rungo/achieving-concurrency-in-go-3f84cbf870ca)
* [Go's work-stealing scheduler](https://rakyll.org/scheduler/)

#### Concurrency patterns

Important concurrency patterns to highlight are:

- <span style="color:orange;">Bridge</span> - is a way to present a single-channel facade over a channel of channels. It is used to consume values from a sequence of channels (channel of channels) doing an ordered write from different sources. By bridging the channels it destructures the channel of channels into a simple channel, allowing to multiplex the input and simplify the consumption.With this pattern we can use the channel of channels from within a single range statement and focus on our loop’s logic.
- <span style="color:orange;">FanIn</span> - combines multiple results in the form of an slice of channels into one channel. This implementation uses a WaitGroup in order to multiplex all the results of the slice of channels. The output is not produced in sequence. This pattern is good for  Non-deterministic stream processing when order is not important and deterministic outcomes are not required.
- <span style="color:orange;">FanInRec</span> - combines multiple results in the form of an slice of channels into one channel. This implementation uses a a recursive approach in order to multiplex all the results of the slice of channels. The output is not produced in sequence. This pattern is good for  Non-deterministic stream processing when order is not important and deterministic outcomes are not required.
- <span style="color:orange;">OrDone</span> - Wraps a channel with a select statement that also selects from a done channel. Allows to cancel the channel avoiding go-routine leaks.
- <span style="color:orange;">Or</span> - returns the value of the fastest channel.
- <span style="color:orange;">Route</span> - Representation of the tee pattern. Takes a single input channel and an arbitrary number of output channels and duplicates each input into every output. When the input channel is closed, all outputs channels are closed. It allows to route or split an input into multiple outputs.

##### Utilities

- <span style="color:orange;">AsChan</span> - sends the contents of a slice through a channel
- <span style="color:orange;">CloseChannel</span> - Checks if the channel is not closed and closes it
- <span style="color:orange;">Take</span> - Takes a defined number of values by num from a channel
- <span style="color:orange;">ToString</span> - Converts any type of channel into a string channel

##### Generators

- <span style="color:orange;">Repeat</span> - Generator that repeats the values defined on the values slice indefinitely.
- <span style="color:orange;">RepeatFn</span> - Generator that repeats a function indefinitely.
- <span style="color:orange;">RepeatParamFn</span> - Generator that repeats a function with one parameter indefinitely.
- <span style="color:orange;">RepeatParamsFn</span> - Generator that repeats a function with a list of parameters indefinitely.
- <span style="color:orange;">RepeatChanParamFn</span> - Generator that repeats a function with a channel as parameter indefinitely.
- <span style="color:orange;">RepeatChanParamsFn</span> - Generator that repeats a function with a list of channels as parameters indefinitely.

#### Thread-safe Maps and Slices

Slice and Map - cannot be used safely from multiple goroutines without the risk of having a race condition.

<span style="color:orange;">gooldi´s</span> provides and implementation of Slice and Map types which can be safely shared between multiple goroutines by protecting the access to the shared data by a mutex. 

What a mutex does is basically to acquire a lock when it needs to access our concurrent type, When holding the lock a goroutine can read and/or write to the shared data protected by the mutex safely, this lock can be acquired by a single goroutine at a time, and if other goroutine needs access to the same shared data it waits until the lock has been released by the goroutine holding the lock. <span style="color:orange;">gooldi´s</span> implementation of these types manage the acquisition and releases of the locks avoiding deadlocks and unwanted behaviours.



What it does

- Done Manager and Done Handler
- gooldi Stream Processing Entities: Message, MessagePair, Broadcast, Processor, Filter,MessageMultiplexer and MultiMessageMultiplexer
- Deterministic Stream Processing
- Non-Deterministic Stream Processing
- Highly customizable

How it does it

