package concurrency

// State establised the state of the relevant object
type State int

// State enum values
const (
	Init       State = iota //Init state
	Processing              // Processing state
	Stopped                 // Stopped state
	Closed                  // Closed state
)
