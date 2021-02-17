package concurrency

//AsChan - sends the contents of a slice through a channel
func AsChan(vs ...interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		for _, v := range vs {
			c <- v
		}
		CloseChannel(c)
	}()
	return c
}
