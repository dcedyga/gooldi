package concurrency

// Bridge is a way to present a single-channel facade over a channel of channels. It is used
// to consume values from a sequence of channels (channel of channels) doing an ordered write
// from different sources. By bridging the channels it destructures the channel of channels
// into a simple channel, allowing to multiplex the input and simplify the consumption.
// With this pattern we can use the channel of channels from within a single range statement
// and focus on our loop’s logic.
func Bridge(done <-chan interface{}, chanStream <-chan <-chan interface{}) chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer CloseChannel(valStream)
		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range OrDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}
