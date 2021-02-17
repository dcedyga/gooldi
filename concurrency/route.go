package concurrency

// Route - Representation of the tee pattern. Takes a single input channel and an arbitrary number of
// output channels and duplicates each input into every output. When the input channel is closed, all
// outputs channels are closed. It allows to route or split an input into multiple outputs.
func Route(done <-chan interface{}, in <-chan interface{}, size int) (channels []chan interface{}) {
	out := make([]chan interface{}, size)
	for i := 0; i < size; i++ {
		out[i] = make(chan interface{})
	}
	go func() {
		for i := 0; i < size; i++ {
			defer CloseChannel(out[i])
		}
		for val := range OrDone(done, in) {
			for i := 0; i < size; i++ {
				select {
				case <-done:
				case out[i] <- val:
				}
			}
		}
	}()
	return out
}
