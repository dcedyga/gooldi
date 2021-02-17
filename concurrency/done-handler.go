package concurrency

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

// DoneHandlerOption - option to initialize the DoneHandler
type DoneHandlerOption func(*DoneHandler)

// DoneHandler - Handles when an object is done and it is ready to be closed.
type DoneHandler struct {
	id       string
	done     chan interface{}
	donefn   func()
	deadline *time.Time
	err      error
}

// NewDoneHandler - Constructor
func NewDoneHandler(opts ...DoneHandlerOption) *DoneHandler {
	id := uuid.NewV4().String()
	dh := &DoneHandler{
		id:       id,
		done:     make(chan interface{}),
		deadline: nil,
	}
	dh.setDonefn()
	for _, opt := range opts {
		opt(dh)
	}
	go dh.doneRn()
	return dh
}

// doneRn - Checks when the DoneHandler is done either by deadline or by closure of the Done channel
func (dh *DoneHandler) doneRn() {
	if dh.deadline != nil {
		select {
		case <-time.After(time.Until(*dh.deadline)):
			dh.close(fmt.Errorf("Deadline reached on DoneHandler: %v", dh.ID()))
		case <-dh.Done():
		}
	} else {
		select {
		case <-dh.Done():
		}
	}
}

// setDonefn - sets the done function
func (dh *DoneHandler) setDonefn() {
	dh.donefn = func() {
		dh.close(fmt.Errorf("Task Finished on DoneHandler: %v", dh.ID()))
	}
}

// ID - retrieves the Id of the DoneHandler
func (dh *DoneHandler) ID() string {
	return dh.id
}

// DoneHandlerWithDeadline - option to add a dealine value to the DoneHandler
func DoneHandlerWithDeadline(deadline time.Time) DoneHandlerOption {
	return func(dh *DoneHandler) {
		dh.deadline = &deadline
	}
}

// DoneHandlerWithTimeout - option to add a timeout value to the DoneHandler. It sets a deadline
// with the value time.Now + timeout
func DoneHandlerWithTimeout(timeout time.Duration) DoneHandlerOption {
	return func(dh *DoneHandler) {
		deadline := time.Now().Add(timeout)
		dh.deadline = &deadline
	}
}

// Done - retrieves the Done channel of the DoneHandler
func (dh *DoneHandler) Done() chan interface{} {
	return dh.done
}

// GetDoneFunc - retrieves the GetDone Function of the DoneHandler
func (dh *DoneHandler) GetDoneFunc() func() {
	return dh.donefn
}

// Err - retrieves the Error of the DoneHandler
func (dh *DoneHandler) Err() error {
	return dh.err
}

// Deadline - retrieves the Deadline of the DoneHandler
func (dh *DoneHandler) Deadline() *time.Time {
	return dh.deadline
}

// close - Closes the Done channel of the handler, allowing the object that is handling to proceed
// with closure actions. It accepts an error to trace if the closure of DoneHandler comes from a
// deadline or a done action.
func (dh *DoneHandler) close(err error) {
	dh.err = err
	CloseChannel(dh.done)
}
