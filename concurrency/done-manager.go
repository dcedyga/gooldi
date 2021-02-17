package concurrency

import (
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

// DoneManagerOption - option to initialize the DoneManager
type DoneManagerOption func(*DoneManager)

// QueryDoneHandlerOption - option to initialize the QueryDoneHandler
type QueryDoneHandlerOption func(*QueryDoneHandler)

// DoneManager - It manages a set of DoneHandlers with a layered approach, allowing to structure
// the way components are closed. It has an ID, a Done channel, a SortedMap of layers which key is
// the number of the layer and the item is represented with a SortedMap of DoneHandlers which item key is
// the DoneHandler.ID. It also have a donefn, a deadline property to setup a deadline to the entire set that
// the DoneManager is handling, a delay to space the closure of layers within the DoneManager, an err and
// a lock to ensure that the operations are threadsafe.
type DoneManager struct {
	id           string
	done         chan interface{}
	doneHandlers *SortedMap
	donefn       func()
	deadline     *time.Time
	delay        *time.Duration
	err          error
	lock         *sync.RWMutex
}

// QueryDoneHandler - Struct with a key and a layer that is used to query the DoneManager by specifying the
// key of a DoneHandler
type QueryDoneHandler struct {
	key   interface{}
	layer int
}

// NewDoneManager - Constructor
func NewDoneManager(opts ...DoneManagerOption) *DoneManager {
	id := uuid.NewV4().String()
	d := 0 * time.Second
	dm := &DoneManager{
		id:           id,
		doneHandlers: NewSortedMap(),
		done:         make(chan interface{}),
		delay:        &d,
		lock:         &sync.RWMutex{},
	}
	dm.setDonefn()
	for _, opt := range opts {
		opt(dm)
	}
	go dm.doneRn()
	return dm
}

// doneRn - Checks when the DoneManager is done either by deadline or by closure of the Done channel
func (dm *DoneManager) doneRn() {
	if dm.deadline != nil {
		select {
		case <-time.After(time.Until(*dm.deadline)):
			dm.closeManager(fmt.Errorf("Deadline reached on DoneManager: %v", dm.ID()))
		case <-dm.done:
		}
	} else {
		select {
		case <-dm.done:
		}
	}
}

// setDonefn - sets the done function
func (dm *DoneManager) setDonefn() {
	dm.donefn = func() {
		dm.closeManager(fmt.Errorf("Task Finished on DoneManager: %v", dm.ID()))
	}
}

// ID - retrieves the Id of the DoneManager
func (dm *DoneManager) ID() string {
	return dm.id
}

// DoneManagerWithDeadline - option to add a dealine value to the DoneManager
func DoneManagerWithDeadline(deadline time.Time) DoneManagerOption {
	return func(dm *DoneManager) {
		dm.deadline = &deadline
	}
}

// DoneManagerWithDelay - option to add a delay value to the DoneManager in order to space in time
// the closure the different layers of the DoneHandler
func DoneManagerWithDelay(delay time.Duration) DoneManagerOption {
	return func(dm *DoneManager) {

		dm.delay = &delay
	}
}

// DoneManagerWithTimeout - option to add a timeout value to the DoneManager. It sets a deadline
// with the value time.Now + timeout
func DoneManagerWithTimeout(timeout time.Duration) DoneManagerOption {
	return func(dm *DoneManager) {
		deadline := time.Now().Add(timeout)
		dm.deadline = &deadline
	}
}

// AddNewDoneHandler - Creates a new DoneHandler with the relevant DoneHandlerOptions and adds it to the
// sorted map of the specified layer
func (dm *DoneManager) AddNewDoneHandler(layer int, opts ...DoneHandlerOption) *DoneHandler {
	dh := NewDoneHandler(opts...)
	dm.AddDoneHandler(dh, layer)
	return dh
}

// AddDoneHandler - adds a DoneHandler to the sorted map of the specified layer
func (dm *DoneManager) AddDoneHandler(dh *DoneHandler, layer int) {
	//Will be added to the sortedmap sequence on the map key
	m := dm.getOrCreateLayerSortedMap(layer)
	dm.lock.Lock()
	m.Set(dh.ID(), dh)
	dm.lock.Unlock()
	// get the donefn thing
	go dm.removeDoneHandlerRn(dh, layer)
}

// doneRn - detects when the DoneHandler is done and removes it from the DoneManager
func (dm *DoneManager) removeDoneHandlerRn(dh *DoneHandler, l int) {
	select {
	case <-dh.Done():
		if !dm.RemoveDoneHandler(QueryDoneHandlerWithKey(dh.ID()), QueryDoneHandlerWithLayer(l)) {
			fmt.Printf("%v - unable to delete doneHandler %v on layer %v\n", dm.ID(), dh.ID(), l)
		}
	}
}

// getOrCreateLayerSortedMap - Gets the SortedMap related to a layer. If not found it creates a new one
// and returns it
func (dm *DoneManager) getOrCreateLayerSortedMap(layer int) *SortedMap {
	var m *SortedMap
	item, ok := dm.getLayerSortedMap(layer)
	if !ok {
		m = NewSortedMap()
		dm.lock.Lock()
		dm.doneHandlers.Set(layer, m)
		dm.lock.Unlock()
	} else {
		m = item
	}
	return m
}

// getLayerSortedMap - Gets the SortedMap related to a layer
func (dm *DoneManager) getLayerSortedMap(layer int) (*SortedMap, bool) {
	dm.lock.Lock()
	item, ok := dm.doneHandlers.Get(layer)
	dm.lock.Unlock()
	if !ok {
		return nil, ok
	}
	return item.(*SortedMap), ok
}

// getQueryDoneHandler - Gets a QueryDoneHandler with the relevant opts for the query.
func (dm *DoneManager) getQueryDoneHandler(opts ...QueryDoneHandlerOption) *QueryDoneHandler {
	q := &QueryDoneHandler{key: nil, layer: -1}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// retrieveDoneHandlerByKeyNoLayerDefined - Retrieves a DoneHandler by key when the layer is not defined
func (dm *DoneManager) retrieveDoneHandlerByKeyNoLayerDefined(q *QueryDoneHandler) (*DoneHandler, int, bool) {
	cancelLayer := make(chan interface{})
	for item := range dm.doneHandlers.IterWithCancel(cancelLayer) {
		m := item.Value.(*SortedMap)
		cancelItem := make(chan interface{})
		for it := range m.IterWithCancel(cancelItem) {
			itdh := it.Value.(*DoneHandler)
			if itdh.ID() == q.key {
				close(cancelItem)
				close(cancelLayer)
				return itdh, item.Key.(int), itdh != nil
			}
		}
	}
	return nil, -1, false
}

// GetDoneHandler - Retrieves a DoneHandler that meets the query defined within the QueryDoneHandlerOptions
func (dm *DoneManager) GetDoneHandler(opts ...QueryDoneHandlerOption) (*DoneHandler, int, bool) {
	q := dm.getQueryDoneHandler(opts...)
	if q.key == nil {
		return nil, -1, false
	}
	if q.layer == -1 {
		return dm.retrieveDoneHandlerByKeyNoLayerDefined(q)
	}
	m, ok := dm.getLayerSortedMap(q.layer)
	if ok {
		dm.lock.Lock()
		it, ok := m.Get(q.key)
		dm.lock.Unlock()
		if !ok {
			return nil, q.layer, false
		}
		return it.(*DoneHandler), q.layer, true
	}

	return nil, -1, false

}

// RemoveDoneHandler - Removes a DoneHandler that meets the query defined within the QueryDoneHandlerOptions.
// If item is Not found it returns false
func (dm *DoneManager) RemoveDoneHandler(opts ...QueryDoneHandlerOption) bool {

	q := dm.getQueryDoneHandler(opts...)
	if q.key == nil {
		return false
	}
	if q.layer == -1 {
		var ok bool
		if _, q.layer, ok = dm.GetDoneHandler(opts...); !ok || q.layer == -1 {
			return false
		}
	}
	m, ok := dm.getLayerSortedMap(q.layer)
	if ok {
		dm.lock.Lock()
		mlen := m.Len()
		m.Delete(q.key)
		if mlen == 1 {
			dm.doneHandlers.Delete(q.layer)
		}
		dm.lock.Unlock()
		return true
	}
	return ok
}

// QueryDoneHandlerWithKey - option to add a key value to the QueryDoneHandler.
func QueryDoneHandlerWithKey(key interface{}) QueryDoneHandlerOption {
	return func(q *QueryDoneHandler) {
		q.key = key
	}
}

// QueryDoneHandlerWithLayer - option to add a layer value to the QueryDoneHandler.
func QueryDoneHandlerWithLayer(layer int) QueryDoneHandlerOption {
	return func(q *QueryDoneHandler) {
		q.layer = layer
	}
}

// GetDoneFunc - retrieves the GetDone Function of the DoneManager
func (dm *DoneManager) GetDoneFunc() func() {
	return dm.donefn
}

// Deadline - retrieves the Deadline of the DoneManager
func (dm *DoneManager) Deadline() *time.Time {
	return dm.deadline
}

// Done - retrieves the Done channel of the DoneManager
func (dm *DoneManager) Done() chan interface{} {
	return dm.done
}

// Err - retrieves the Error of the DoneManager
func (dm *DoneManager) Err() error {
	return dm.err
}

// closeManager - It closes all the DoneHandlers that are managed by the DoneManager
func (dm *DoneManager) closeManager(err error) {
	dhs := []interface{}{}
	for layerMap := range dm.doneHandlers.Iter() {
		m := layerMap.Value.(*SortedMap)
		mslice := []*DoneHandler{}
		for it := range m.Iter() {
			mslice = append(mslice, it.Value.(*DoneHandler))
		}
		dhs = append(dhs, mslice)
	}
	for _, item := range dhs {
		for _, it := range item.([]*DoneHandler) {
			go func(dh *DoneHandler) {
				dh.GetDoneFunc()()
			}(it)
		}
		time.Sleep(*(dm.delay))
	}
	dhs = nil
	dm.err = err
	CloseChannel(dm.done)

}
