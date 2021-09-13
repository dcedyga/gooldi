package concurrency

// Repeat - Generator that repeats the values defined on the values slice indefinitely.
func Repeat(done <-chan interface{}, values ...interface{}) chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer CloseChannel(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

// RepeatFn - Generator that repeats a function indefinitely.
func RepeatFn(done <-chan interface{}, fn func() interface{}) chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer CloseChannel(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}

// RepeatParamFn - Generator that repeats a function with one parameter indefinitely.
func RepeatParamFn(done <-chan interface{}, fn func(param interface{}) interface{}, param interface{}) chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer CloseChannel(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn(param):
			}
		}
	}()
	return valueStream
}

// RepeatParamsFn - Generator that repeats a function with a list of parameters indefinitely.
func RepeatParamsFn(done <-chan interface{}, fn func(params ...interface{}) interface{}, params ...interface{}) chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer CloseChannel(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn(params...):
			}
		}
	}()
	return valueStream
}

// // RepeatChanParamFn - Generator that repeats a function with a channel as parameter indefinitely.
// func RepeatChanParamFn(done <-chan interface{}, fn func(param chan interface{}) interface{}, param chan interface{}) chan interface{} {
// 	valueStream := make(chan interface{})
// 	go func() {
// 		defer CloseChannel(valueStream)
// 		for {
// 			select {
// 			case <-done:
// 				return
// 			case valueStream <- fn(param):
// 			}
// 		}
// 	}()
// 	return valueStream
// }

// // RepeatChanParamsFn - Generator that repeats a function with a list of channels as parameters indefinitely.
// func RepeatChanParamsFn(done <-chan interface{}, fn func(params ...chan interface{}) interface{}, params ...chan interface{}) chan interface{} {
// 	valueStream := make(chan interface{})
// 	go func() {
// 		defer CloseChannel(valueStream)
// 		for {
// 			select {
// 			case <-done:
// 				return
// 			case valueStream <- fn(params...):
// 			}
// 		}
// 	}()
// 	return valueStream
// }
