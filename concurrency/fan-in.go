package concurrency

import (
	"sync"
)

// FanIn - combines multiple results in the form of an slice of channels into one channel.
// This implementation uses a waitgroup in order to multiplex all the results of the slice of channels.
// The output is not produced in sequence. This pattern is good when order is not important
func FanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	// Select from all the channels
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	// Wait for all the reads to complete
	go func() { // <5>
		wg.Wait()
		CloseChannel(multiplexedStream)
	}()

	return multiplexedStream
}

// FanInRec - combines multiple results in the form of an slice of channels into one channel.
// This implementation uses a a recursive approach in order to multiplex all the results of the slice of channels.
// The output is not produced in sequence. This pattern is good when order is not important
func FanInRec(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		c := make(chan interface{})
		CloseChannel(c)
		return c
	case 1:
		return channels[0]
	case 2:
		return fanInTwo(done, channels[0], channels[1])
	default:
		m := len(channels) / 2
		return fanInTwo(done,
			FanInRec(done, channels[:m]...),
			FanInRec(done, channels[m:]...))
	}
}

// fanInTwo - multiplexes two channels
func fanInTwo(done <-chan interface{}, a, b <-chan interface{}) <-chan interface{} {
	c := make(chan interface{})

	go func() {
		defer CloseChannel(c)
		for a != nil || b != nil {
			select {
			case <-done:
				return
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				c <- v
			}
		}
	}()
	return c
}
