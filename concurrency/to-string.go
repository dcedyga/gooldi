package concurrency

import (
	"fmt"
)

// ToString - Converts any type of channel into a string channel
func ToString(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {
	stringStream := make(chan string)
	go func() {

		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- fmt.Sprintf("%v", v):
			}
		}
	}()
	return stringStream
}
