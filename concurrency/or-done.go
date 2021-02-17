package concurrency

import (
	"fmt"
)

// OrDone - Wraps a channel with a select statement that also selects from a done channel. Allows to cancel the channel
// avoiding go-routine leaks
func OrDone(done, c <-chan interface{}) chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer CloseChannel(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

// OrDoneFn - Wraps a channel with a select statement that also selects from a done channel with a function that can
// be run before returning. Allows to cancel the channel avoiding go-routine leaks
func OrDoneFn(done, c <-chan interface{}, fn func()) chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer CloseChannel(valStream)
		for {
			select {
			case <-done:
				fn()
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

// OrDoneParamFn - Wraps a channel with a select statement that also selects from a done channel with a function
// with parameter that can be run before returning. Allows to cancel the channel avoiding go-routine leaks.
func OrDoneParamFn(done, c <-chan interface{}, fn func(param interface{}), param interface{}) chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer CloseChannel(valStream)
		for {
			select {
			case <-done:
				fn(param)
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

// OrDoneParamsFn - Wraps a channel with a select statement that also selects from a done channel with a function
// with a list of parameters that can be run before returning. Allows to cancel the channel avoiding go-routine
// leaks.
func OrDoneParamsFn(done, c <-chan interface{}, fn func(params ...interface{}), params ...interface{}) chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer CloseChannel(valStream)
		for {
			select {
			case <-done:
				fn(params...)
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

// OrDoneChanParamFn - Wraps a channel with a select statement that also selects from a done channel with a function
// with a channel as parameter that can be run before returning. Allows to cancel the channel avoiding go-routine
// leaks.
func OrDoneChanParamFn(done, c <-chan interface{}, fn func(param chan interface{}), param chan interface{}) chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer CloseChannel(valStream)
		for {
			select {
			case <-done:
				fmt.Printf("We are done with: %v\n", param)
				fn(param)
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

// OrDoneChanParamsFn - Wraps a channel with a select statement that also selects from a done channel with a function
// with a list of channels as parameters that can be run before returning. Allows to cancel the channel avoiding go-routine
// leaks.
func OrDoneChanParamsFn(done, c <-chan interface{}, fn func(params ...chan interface{}), params ...chan interface{}) chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer CloseChannel(valStream)
		for {
			select {
			case <-done:
				fn(params...)
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}
