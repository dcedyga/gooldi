package concurrency

import (
	"fmt"
	"strconv"
	"strings"
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
func IndexToString(idx int64) string {
	s := strconv.FormatInt(idx, 10)
	return strings.Repeat("0", 20-len(s)) + s
}
