package concurrency

//CloseChannel - Checks if the channel is not closed and closes it
func CloseChannel(ch chan interface{}) {
	select {
	case <-ch:
	default:
		close(ch)
	}
}
