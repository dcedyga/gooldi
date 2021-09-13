package concurrency

import (
	"fmt"
	//"strconv"
	"sync"
	//"time"
	uuid "github.com/satori/go.uuid"
)

// MsgMultiplexerOption - option to initialize the MsgMultiplexer
type MsgMultiplexerOption func(*MsgMultiplexer)
type MsgMultiplexerItem struct {
	value  chan interface{}
	cancel chan interface{}
	index  interface{}
}
type ActionRegistryItem struct {
	ckey int64
	indexesToDelete *Slice
	indexesToAdd *Slice
	len int
	elen int
}
// MsgMultiplexer - The default implementation MsgMultiplexer allows to create complex patterns where a Broadcaster can emit an message to multiple processors (consumers)
// that can potentially represent multiple processing systems, do the relevant calculation and multiplex the multiple outputs
// into a single channel for simplified consumption.
// Its main function is to Mulitplex a set of concurrent processors that process a common initial concurrency.Message
// converging them into one channel,where the output is a concurrency.Message which Message property is a sortedmap of the output
// values of the processors grouped by initial concurrency.Message CorrelationKey and ordered by index value of each processor.
// Closure of MsgMultiplexer is handle by a concurrency.DoneHandler that allows to control they way a set of go routines
// are closed in order to prevent deadlocks and unwanted behaviour.
// MsgMultiplexer outputs the multiplexed result in one channel using the channel bridge pattern.
// MsgMultiplexer default behaviour can be overridden by providing a MsgMultiplexerGetItemKeyFn to provide the comparison key of
// the items of a channel, with this function MsgMultiplexer
// has an algorithm to group the processed messages related to the same source into a SortedMap.
// A transformFn can be provided as well to map the output of the MsgMultiplexer to a specific structure. The defaultTransformFn
// transform the SortedMap of processor outputs into a Message which Message property is populated with this SortedMap
type MsgMultiplexer struct {
	id            string
	inputChannels *Map
	doneHandler   *DoneHandler
	//doneBridgeChan          chan interface{}
	debug                       			string
	preStreamMap                			*SortedMap
	//inputChannelsLastRegLen     			int
	//inputChannelsEffectiveLen 			int
	inputChannelTransitionsInitLen 			int
	taggedToDeleteMap						*SortedMap
	taggedToAddMap 							*SortedMap
	transitionRegistryMap					*SortedMap
	startFrom								int64
	stopAt									int64
	lastInsertedPreStreamMapKey 			int64
	initProcesses							*SortedMap
	lastSentPreStreamMapKey					int64
	stream                      			chan (<-chan interface{})
	lock                        			*sync.RWMutex
	index                       			int64
	//dirty                       			bool
	toStringIndex               			string
	minIndexKey                 			interface{}
	MsgType                     			string
	state                       			State
	getItemKeyFn                			func(v interface{}) int64
	transformFn                 			func(mp *MsgMultiplexer, sm *SortedMap, correlationKey int64) interface{}
}
//NewMsgMultiplexer - Constructor
func NewMsgMultiplexer(dh *DoneHandler, msgType string, opts ...MsgMultiplexerOption) *MsgMultiplexer {
	id := uuid.NewV4().String()
	mp := &MsgMultiplexer{
		id:                          	id,
		inputChannels:               	NewMap(),
		transitionRegistryMap: 			NewSortedMap(),
		taggedToDeleteMap: 				NewSortedMap(),
		taggedToAddMap: 				NewSortedMap(),
		initProcesses: 					NewSortedMap(),
		doneHandler:                 	dh,
		startFrom: 						0,
		stopAt:							0,
		//inputChannelsLastRegLen:     	0,
		//inputChannelsEffectiveLen: 	0,
		inputChannelTransitionsInitLen: 0,
		preStreamMap:                	NewSortedMap(),
		stream:                      	make(chan (<-chan interface{})),
		debug:                       	"",
		lastInsertedPreStreamMapKey: 	0,
		lastSentPreStreamMapKey: 		0,
		//doneBridgeChan:          		make(chan interface{}),
		//dirty:                		false,
		state:                			Init,
		index:                			0,
		toStringIndex:        			"00000000000000000000",
		minIndexKey:          			"",
		MsgType:              			msgType,
		lock:                 			&sync.RWMutex{},
	}
	mp.getItemKeyFn = defaultGetItemKeyFn
	mp.transformFn = defaultTransformFn
	for _, opt := range opts {
		opt(mp)
	}
	go mp.doneRn()
	return mp
}
// doneRn - Checks when the MsgMultiplexer is done by listening to the closure of the DoneHandler.Done channel
func (mp *MsgMultiplexer) doneRn() {
	select {
	case <-mp.doneHandler.Done():
		mp.close()
	}

}
// MsgMultiplexerIndex - option to add the index value to the MsgMultiplexer
func MsgMultiplexerIndex(idx int64) MsgMultiplexerOption {
	return func(mp *MsgMultiplexer) {
		mp.index = idx
		mp.toStringIndex = IndexToString(idx)
	}
}
// MsgMultiplexerGetItemKeyFn - option to add a function to resolve the key value
// of an item of the channel to the MsgMultiplexer
func MsgMultiplexerGetItemKeyFn(fn func(v interface{}) int64) MsgMultiplexerOption {
	return func(mp *MsgMultiplexer) {
		mp.getItemKeyFn = fn
	}
}
// MsgMultiplexerTransformFn - option to add a function to transform the SortedMap output into
// the desired output structure to the MsgMultiplexer
func MsgMultiplexerTransformFn(fn func(mp *MsgMultiplexer, sm *SortedMap, correlationKey int64) interface{}) MsgMultiplexerOption {
	return func(mp *MsgMultiplexer) {
		mp.transformFn = fn
	}
}
// ID - retrieves the Id of the MsgMultiplexer
func (mp *MsgMultiplexer) ID() string {
	return mp.id
}
// Index - retrieves the Index of the MsgMultiplexer
func (mp *MsgMultiplexer) Index() int64 {
	return mp.index
}
// DoneHandler - retrieves the DoneHandler of the Multiplexer
func (mp *MsgMultiplexer) DoneHandler() *DoneHandler {
	return mp.doneHandler
}
// ToStringIndex - retrieves the ToStringIndex representation of the MsgMultiplexer
func (mp *MsgMultiplexer) ToStringIndex() string {
	return mp.toStringIndex
}
// Stop - stops the processor.
func (mp *MsgMultiplexer) Stop() {
	mp.lock.Lock()
	mp.stopIt()
	mp.lock.Unlock()

}
func (mp *MsgMultiplexer) stopIt() {
	mp.state = Stopped
	if mp.preStreamMap.Len()>0{
		mp.preStreamMap=nil
		mp.preStreamMap=NewSortedMap()
	}
}
func (mp *MsgMultiplexer) InitializeProcessors() {
	mp.lock.Lock()
	mp.inputChannelTransitionsInitLen=mp.inputChannels.Len()
	
	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.startFrom: " + fmt.Sprintf("%v",mp.startFrom) + "\n"
	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.inputChannelTransitionsInitLen: " + fmt.Sprintf("%v",mp.inputChannelTransitionsInitLen) + "\n"
	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.preStreamMap: " + fmt.Sprintf("%v",mp.preStreamMap) + "\n"
	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.initProcesses: " + fmt.Sprintf("%v",mp.initProcesses) + "\n"
	
	mp.lock.Unlock()
	for item := range mp.initProcesses.Iter() {
		mpi:=item.Value.(*MsgMultiplexerItem)
		go mp.processChannel(mpi.index, mpi.value, mpi.cancel)
	}
	mp.clearInitProcesses()
}
// Start - starts the MsgMultiplexer.
func (mp *MsgMultiplexer) Start() {
	mp.lock.Lock()
	mp.state = Processing
	mp.inputChannelTransitionsInitLen=mp.inputChannels.Len()
	mp.startFrom=mp.lastInsertedPreStreamMapKey

	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.startFrom: " + fmt.Sprintf("%v",mp.startFrom) + "\n"
	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.inputChannelTransitionsInitLen: " + fmt.Sprintf("%v",mp.inputChannelTransitionsInitLen) + "\n"
	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.preStreamMap: " + fmt.Sprintf("%v",mp.preStreamMap) + "\n"
	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.initProcesses: " + fmt.Sprintf("%v",mp.initProcesses) + "\n"
	
	mp.lock.Unlock()
	for item := range mp.initProcesses.Iter() {
		mpi:=item.Value.(*MsgMultiplexerItem)
		go mp.processChannel(mpi.index, mpi.value, mpi.cancel)
	}
	mp.clearInitProcesses()
}
// StopAfterCorrelationKey - stops the processor at a correlationKey.
func (mp *MsgMultiplexer) StopAfterCorrelationKey(cKey int64)  {
	mp.lock.Lock()
	mp.stopAt=cKey
	mp.lock.Unlock()
}
// Start - starts the MsgMultiplexer.
func (mp *MsgMultiplexer) StartAfterCorrelationKey(cKey int64) {
	mp.lock.Lock()
	mp.state = Processing
	mp.inputChannelTransitionsInitLen=mp.inputChannels.Len()
	mp.startFrom=cKey
	mp.stopAt=0
	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.startFrom: " + fmt.Sprintf("%v",mp.startFrom) + "\n"
	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.inputChannelTransitionsInitLen: " + fmt.Sprintf("%v",mp.inputChannelTransitionsInitLen) + "\n"
	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.preStreamMap: " + fmt.Sprintf("%v",mp.preStreamMap) + "\n"
	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.initProcesses: " + fmt.Sprintf("%v",mp.initProcesses) + "\n"
	
	mp.lock.Unlock()
	for item := range mp.initProcesses.Iter() {
		mpi:=item.Value.(*MsgMultiplexerItem)
		go mp.processChannel(mpi.index, mpi.value, mpi.cancel)
	}
	mp.clearInitProcesses()
}
func (mp *MsgMultiplexer) clearInitProcesses() {
	mp.lock.Lock()
	mp.initProcesses=nil
	mp.initProcesses=NewSortedMap()
	mp.lock.Unlock()
}
// Set - Registers a channel in the MsgMultiplexer, starts processing it
// and logs the length of the registered channels map.
func (mp *MsgMultiplexer) Set(key interface{}, value chan interface{}) {
	mp.add(key, value ,-1) 
}
// SetAtCorrelationKey - Registers a channel in the MsgMultiplexer at a correlation Key, starts processing it
// and logs the length of the registered channels map.
func (mp *MsgMultiplexer) SetAtCorrelationKey(key interface{}, value chan interface{},cKey int64) {
	mp.add(key, value ,cKey)
}
func (mp *MsgMultiplexer) add(key interface{}, value chan interface{},compKey int64) {
	mp.lock.Lock()
	mpi := &MsgMultiplexerItem{value: value, cancel: make(chan interface{}), index: key}
	mp.inputChannels.Set(key, mpi)
	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.lastSentPreStreamMapKey: " + fmt.Sprintf("%v",mp.lastSentPreStreamMapKey) + "\n"
	//mp.inputChannelsLastRegLen = mp.inputChannels.Len()
	//mp.debug += "Line 178 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Set-mp.inputChannels.Len(): " + fmt.Sprintf("%v",mp.inputChannels.Len()) + "\n"
	if mp.minIndexKey == "" || mp.isMinIndexKey(key) {
		mp.minIndexKey = key
	} 
	
	if mp.state == Processing {
		if compKey==-1{
			compKey=mp.lastSentPreStreamMapKey
			mp.taggedToAddMap.Set(key,compKey)
		}else {
			mp.taggedToAddMap.Set(key,compKey)
			mp.setTransitionRegistryItem(compKey,key,true)
			mp.reconcileTransitionRegistryItems(compKey)
		}
		mp.lock.Unlock()
		go mp.processChannel(mpi.index, mpi.value, mpi.cancel)
		return
	}
	mp.initProcesses.Set(key,mpi)
	mp.lock.Unlock()
	
}
func (mp *MsgMultiplexer) isMinIndexKey(key interface{}) bool {
	s := fmt.Sprintf("%v", key)
	s2 := fmt.Sprintf("%v", mp.minIndexKey)
	return s < s2
}
func (mp *MsgMultiplexer) setLastInsertedPreStreamMapKey(ckey int64) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	if ckey > mp.lastInsertedPreStreamMapKey {
		mp.lastInsertedPreStreamMapKey=ckey
	}
}
func (mp *MsgMultiplexer) getLastSentPreStreamMapKey()int64{
	mp.lock.Lock()
	defer mp.lock.Unlock()
	return mp.lastSentPreStreamMapKey
}
// Get - Retrieves a channel reqistered in the MsgMultiplexer by key
func (mp *MsgMultiplexer) Get(key interface{}) (chan interface{}, bool) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	v, ok := mp.inputChannels.Get(key)
	if !ok {
		return nil, false
	}
	return v.(*MsgMultiplexerItem).value, true
}
// Get - Retrieves a channel reqistered in the MsgMultiplexer by key
func (mp *MsgMultiplexer) getMsgMultiplexerItem(key interface{}) (*MsgMultiplexerItem, bool) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	v, ok := mp.inputChannels.Get(key)
	if !ok {
		return nil, false
	}
	return v.(*MsgMultiplexerItem), true
}
// DeleteAtCorrelationKey - Deletes a registered channel after a defined correlation Key from the MsgMultiplexer map of inputChannels
// and logs the length of the remaining registered channels map.
func (mp *MsgMultiplexer) DeleteAtCorrelationKey(index interface{}, cKey int64) {
	if c, ok := mp.getMsgMultiplexerItem(index); ok {
		mp.tagIndexToDelete(index,cKey)
		CloseChannel(c.cancel)
		//time.Sleep(1000 * time.Microsecond)
	}
}
func (mp *MsgMultiplexer) clean(key interface{}) {
	// mp.lock.Lock()
	// defer mp.lock.Unlock()
	mp.inputChannels.Delete(key)
	//mp.inputChannelsLastRegLen = mp.inputChannels.Len()

}
// Delete - Deletes a registered channel from the MsgMultiplexer map of inputChannels
// and logs the length of the remaining registered channels map.
func (mp *MsgMultiplexer) Delete(index interface{}) {
	if c, ok := mp.getMsgMultiplexerItem(index); ok {
		
		mp.tagIndexToDelete(index,-1)
		CloseChannel(c.cancel)
		//time.Sleep(1000 * time.Microsecond)
	}
}
func (mp *MsgMultiplexer) setTransitionRegistryItem(key int64,index interface{},add bool) {
	item, ok := mp.transitionRegistryMap.Get(key)
	var rm *ActionRegistryItem
	if ok {
		//mp.debug += "Line 248 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") setTransitionRegistryItem(add: " + fmt.Sprintf("%v",add) + ") - found at key: " + fmt.Sprintf("%v",key) + "\n"
		rm = item.(*ActionRegistryItem)
		if add {
			rm.indexesToAdd.Append(index)
		}else {
			rm.indexesToDelete.Append(index)
		}
		
	} else {
		s := NewSlice()
		a := NewSlice()
		if add {
			a.Append(index)
		}else {
			s.Append(index)
		}
		rm = &ActionRegistryItem{
			ckey: key,
			indexesToDelete: s,
			indexesToAdd: a,
			len: 0,
			elen: 0,
		}
	}
	//mp.debug += "Line 272 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") setTransitionRegistryItem(add: " + fmt.Sprintf("%v",add) + ")-item: " + fmt.Sprintf("%v",rm) + "\n"
	mp.transitionRegistryMap.Set(key,rm)
}
func (mp *MsgMultiplexer) reconcileTransitionRegistryItems(cKey int64) {
	if (mp.transitionRegistryMap.Len()>0){
		var item *ActionRegistryItem
		i:=0		
		for it:= range mp.transitionRegistryMap.Iter(){
			compItem:=it.Value.(*ActionRegistryItem)
			if (i==0){
				item=compItem
				item.len = mp.inputChannelTransitionsInitLen + item.indexesToAdd.Len()
				item.elen = item.len - item.indexesToDelete.Len()
				//mp.inputChannelsEffectiveLen=item.elen
			}else {
				compItem.len = item.elen + compItem.indexesToAdd.Len()
				compItem.elen = compItem.len - compItem.indexesToDelete.Len()
				item=compItem
			}
			//mp.debug += "Line 355 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") reconcileTransitionRegistryItems-item: " + fmt.Sprintf("%v",item) + ", compItem:" + fmt.Sprintf("%v",compItem) +"\n"
			i++
		}
	}
//	mp.debug += "Line 297 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") reconcileTransitionRegistryItems.transitionRegistryMap: " + fmt.Sprintf("%v",mp.transitionRegistryMap) +"\n"
}
func (mp *MsgMultiplexer) tagIndexToAdd(index interface{},key int64) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	if mp.state == Processing {
		mp.taggedToAddMap.Set(index,key)
	}
}
func (mp *MsgMultiplexer) tagIndexToDelete(index interface{},key int64) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	if mp.state == Processing {
		mp.taggedToDeleteMap.Set(index,key)
	}
}
func (mp *MsgMultiplexer) getTagIndexToAdd(index interface{})(int64,bool) {
	key,ok:=mp.taggedToAddMap.Get(index)
	if !ok {
		return int64(0), ok
	}
	return key.(int64),ok

}
func (mp *MsgMultiplexer) getTagIndexToDelete(index interface{})(int64,bool) {
	key,ok:=mp.taggedToDeleteMap.Get(index)
	if !ok {
		return int64(0), ok
	}
	return key.(int64),ok

}
func (mp *MsgMultiplexer) cleanAndCloseProcessor(cKey int64,key interface{}, value chan interface{}) bool{
	mp.lock.Lock()
	defer mp.lock.Unlock()
	dKey,ok:= mp.getTagIndexToDelete(key)
	if ok && (dKey == -1 || cKey == dKey) {
		mp.setTransitionRegistryItem(cKey,key,false)
		mp.taggedToDeleteMap.Delete(key)
		mp.reconcileTransitionRegistryItems(cKey)
		CloseChannel(value)
		mp.clean(key)
		return true
	}else if mp.state != Processing {
		CloseChannel(value)
		mp.clean(key)
		return true
	}
	return false
}
func (mp *MsgMultiplexer) startSendingProcessedItem(cKey int64,key interface{}) (bool, int){
	mp.lock.Lock()
	defer mp.lock.Unlock()
	lsKey := int64(0)
	tagIKey,ikeyOk:=mp.getTagIndexToAdd(key)
	if ikeyOk {
		lsKey = mp.lastSentPreStreamMapKey
		if cKey==tagIKey{
			mp.taggedToAddMap.Delete(key)
			return true, 1
		}else if cKey > lsKey {
			mp.setTransitionRegistryItem(cKey,key,true)
			mp.taggedToAddMap.Delete(key)
			mp.reconcileTransitionRegistryItems(cKey)
			//mp.debug += "Line 365 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Item: " + fmt.Sprintf("%v", key) + " starting at CKey: " + strconv.FormatInt(cKey, 10) + "\n"
			return true, 1
		}
		return false, 0
	} 
	return true, 1
}
// processChannel - Retrieves all the values of an inputChannel using range and when done it deletes it
// from the MsgMultiplexer map of inputChannels.
func (mp *MsgMultiplexer) processChannel(key interface{}, value chan interface{}, cancel chan interface{}) {
	cKey := int64(0)
	closed := false
	send:=true
	i:=0
	for item := range value {
		if !closed {
			
			cKey = mp.getItemKeyFn(item)
			mp.setLastInsertedPreStreamMapKey(cKey)
			select {
			case <-cancel:
				closed=mp.cleanAndCloseProcessor(cKey,key,value)
				mp.storeInPreStreamMap(cKey, key, item)
				mp.sendWhenReady(cKey)
			default:
				if i==0 {
					send,i=mp.startSendingProcessedItem(cKey,key) 
				}
				if send {
					mp.storeInPreStreamMap(cKey, key, item)
					mp.sendWhenReady(cKey)
				}
			}
		}
	}
	//fmt.Printf("We are out at cKey:%v\n", cKey)
	if _, ok := mp.Get(key); ok {
		mp.lock.Lock()
		mp.clean(key)
		mp.lock.Unlock()
	}

}
//storeInPreStreamMap - store stream item in PreStreamMap per correlationKey
func (mp *MsgMultiplexer) storeInPreStreamMap(cKey int64, index interface{}, v interface{}) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	if mp.state == Processing && cKey > mp.startFrom {
		item, ok := mp.preStreamMap.Get(cKey)
		var rm *SortedMap
		if ok {
			rm = item.(*SortedMap)
			rm.Set(index, v) // Sequence
		} else {
			rm = NewSortedMap()
			rm.Set(index, v) // Sequence
			mp.preStreamMap.Set(cKey, rm)
			mp.lastInsertedPreStreamMapKey = cKey
		}
	}
	//mp.debug += "Line 176 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") storeInPreStreamMap-mp.preStreamMap: " + fmt.Sprintf("%v",mp.preStreamMap) + ", index: "+ fmt.Sprintf("%v",key) + ", cKey: "+ fmt.Sprintf("%v",sequence) + "\n"
	
}
// close - Closes the MsgMultiplexer
func (mp *MsgMultiplexer) close() {
	fmt.Printf("Multiplexer - is closed\n")
}
func (mp *MsgMultiplexer) sendWhenReady(cKey int64) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	if mp.state == Processing && cKey > mp.startFrom {
		mp.dispatchIfReady(cKey,true)
	}
}
func (mp *MsgMultiplexer) getLenToCompareFromCeroItem(cKey int64,initLen int)(int,bool){
	cItemAtCero,_:=mp.transitionRegistryMap.GetSortedMapItemByIndex(0)
	cICero:=cItemAtCero.Value.(*ActionRegistryItem)
	if cKey < cICero.ckey {
	//	mp.debug += "Line 436 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") for item: cKey < cICero.ckey: " + fmt.Sprintf("%v",cICero.ckey) + ", cKey: " + fmt.Sprintf("%v",cKey) + ", l: " + fmt.Sprintf("%v",cICero.len - cICero.indexesToAdd.Len()) + "\n"			
		return cICero.len - cICero.indexesToAdd.Len(),false
	} else if cKey == cICero.ckey {
	//	mp.debug += "Line 439 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") for item: cKey == cICero.ckey: " + fmt.Sprintf("%v",cICero.ckey) + ", cKey: " + fmt.Sprintf("%v",cKey) + ", l: " + fmt.Sprintf("%v",cICero.len) + "\n"		
		return cICero.len,false
	}
	return cICero.len,true
}
func (mp *MsgMultiplexer) getLenToCompareFromNonCeroItems(cKey int64,initLen int,rMapLen int)(int){
	l:=initLen
	//Check if cKey > last item key
	cItemLast,_:=mp.transitionRegistryMap.GetSortedMapItemByIndex(rMapLen-1)
	cILast:=cItemLast.Value.(*ActionRegistryItem)
	if cKey > cILast.ckey {
	//	mp.debug += "Line 450 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") for item: cKey > cILast.ckey: " + fmt.Sprintf("%v",cILast.ckey) + ", cKey: " + fmt.Sprintf("%v",cKey) + ", l: " + fmt.Sprintf("%v",cILast.len) + "\n"		
		l=cILast.len
	}else {
		cancel := make(chan interface{})
		var compItem *ActionRegistryItem
		i:=0
		for item := range mp.transitionRegistryMap.IterWithCancel(cancel){
			arKey:=item.Key.(int64)
			arItem:=item.Value.(*ActionRegistryItem)
			if (cKey==arKey){
				//Check if cKey is one of the ActionRegistries
				l=arItem.len
			//	mp.debug += "Line 462 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") for item: cKey==arKey - arKey: " + fmt.Sprintf("%v",arKey) + ", cKey: " + fmt.Sprintf("%v",cKey) + ", arItem: " + fmt.Sprintf("%v",arItem) + ", l: " + fmt.Sprintf("%v",l) + "\n"		
				close(cancel)
				break
			}else if i>0 && cKey>compItem.ckey && cKey<arItem.ckey {
				//Check if cKey in between the ActionRegistries
				l=compItem.elen
			//	mp.debug += "Line 468 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") for item: cKey==arKey - arKey: " + fmt.Sprintf("%v",arKey) + ", cKey: " + fmt.Sprintf("%v",cKey) + ", compItem: " + fmt.Sprintf("%v",compItem) + ", l: " + fmt.Sprintf("%v",l) + "\n"		
				close(cancel)
				break
			}	
			i++
			compItem=arItem
		}	
	}
	return l
}
func (mp *MsgMultiplexer) getLenToCompare(cKey int64,initLen int)(int){
	l:=initLen
	greaterThanCeroKey:=false
	rMapLen:=mp.transitionRegistryMap.Len()
	if rMapLen>0{
		l,greaterThanCeroKey=mp.getLenToCompareFromCeroItem(cKey,initLen)
		if greaterThanCeroKey && rMapLen>1{
			l=mp.getLenToCompareFromNonCeroItems(cKey,initLen,rMapLen)
		}
	}
	return l
}
func (mp *MsgMultiplexer) canBeSent(cKey int64) bool{
	if (mp.taggedToAddMap.Len()>0){
		cFirstTaggedToAdd,_:=mp.taggedToAddMap.GetSortedMapItemByIndex(0)
		approxcKeyItemToAdd:=cFirstTaggedToAdd.Value.(int64)
		if (cKey>approxcKeyItemToAdd){
			//mp.debug += "Line 495 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") approxcKeyItemToAdd: " + strconv.FormatInt(approxcKeyItemToAdd, 10) + ", cKey: " + strconv.FormatInt(cKey, 10) + "\n"	
			return false
		}
	}
	return true
}
func (mp *MsgMultiplexer) sendPreviousToCurrentItems(cKey int64,cTransitionItemOk bool, parent bool)bool{
	pSent:=true
	if cTransitionItemOk {
		if (parent && mp.preStreamMap.Len() > 0){
		//	mp.debug += "Line 505 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") Parent (rItem.indexesToAdd.Len()>0) : " + fmt.Sprintf("%v",cKey) + "\n"	
			keys := mp.preStreamMap.GetKeys()
		//	mp.debug += "Line 507 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") mp.preStreamMap.GetKeys() : " + fmt.Sprintf("%v",keys) + "\n"	
			for _, ckey := range keys {
				if ckey.(int64)< cKey {
					pSent=mp.dispatchIfReady(ckey.(int64),false)
				//	mp.debug += "Line 511 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") childrensent: " + fmt.Sprintf("%v",childrensent) + " for key: " + fmt.Sprintf("%v",ckey) + "\n"	
					if !pSent {
						break
					}
				}
			}
		}
	}
	return pSent
}
func (mp *MsgMultiplexer) cleanEntries(cTransitionItemOk bool,cTransitionMapItem interface{},sm *SortedMap){
	if cTransitionItemOk {
		trItem:=cTransitionMapItem.(*ActionRegistryItem)
		if trItem.indexesToDelete.Len()>0 {
			for itemToRemove:= range trItem.indexesToDelete.Iter(){
				ir:=itemToRemove.Value.(int64)
				sm.Delete(ir)
			}
		}
	}
}
func (mp *MsgMultiplexer) sendPostToCurrentItems(cKey int64,trItem *ActionRegistryItem, parent bool){
	if parent && mp.preStreamMap.Len() > 0 && trItem.indexesToDelete.Len()>0{
		pSent:=true
		keys := mp.preStreamMap.GetKeys()
		for _, ckey := range keys {
			if ckey.(int64)>cKey {
			//	mp.debug += "Line 538 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") try send: " + fmt.Sprintf("%v",ckey) + "\n"	
				pSent=mp.dispatchIfReady(ckey.(int64),false)
				if !pSent {
					break
				}
			}
		}	
	}
}
func (mp *MsgMultiplexer) finalizeDispatch(cKey int64,cTransitionItemOk bool,cTransitionMapItem interface{}, parent bool){
	if cTransitionItemOk {
		trItem:=cTransitionMapItem.(*ActionRegistryItem)
		mp.inputChannelTransitionsInitLen=trItem.elen
		mp.transitionRegistryMap.Delete(cKey) 
		mp.sendPostToCurrentItems(cKey,trItem,parent)
	}
}
func (mp *MsgMultiplexer) dispatchIfReady(cKey int64,parent bool) bool {
	sent:=false
	if mp.state == Processing {
		m, ok := (mp.preStreamMap.Get(cKey))
		if ok && mp.canBeSent(cKey){
			sm := m.(*SortedMap)
			l := sm.Len()
			//	mp.debug += "Line 561 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") cKey: " + fmt.Sprintf("%v",cKey) + ", sm.len: " + fmt.Sprintf("%v",l) + ", mp.inputChannelTransitionsInitLen: " + fmt.Sprintf("%v",mp.inputChannelTransitionsInitLen) + "\n"	
			lenToCompare:=mp.getLenToCompare(cKey,mp.inputChannelTransitionsInitLen)
			if l == lenToCompare {
				cTransitionMapItem,cTransitionItemOk:=mp.transitionRegistryMap.Get(cKey)
				if mp.sendPreviousToCurrentItems(cKey,cTransitionItemOk, parent) {
					mp.cleanEntries(cTransitionItemOk,cTransitionMapItem,sm)
					//send to stream
					e := mp.transformFn(mp, sm, cKey)
					s := make(chan interface{}, 1)
					s <- e
					close(s)
					select {
					case <-mp.doneHandler.Done():
					case mp.stream <- s:
						mp.lastSentPreStreamMapKey=cKey
						sent=true
					//	mp.debug += "Line 577 -(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ") sent cKey: " + fmt.Sprintf("%v",cKey) + "- parent: " + fmt.Sprintf("%v",parent) + "\n"
					}
					//delete
					mp.preStreamMap.Delete(cKey)
					mp.finalizeDispatch(cKey,cTransitionItemOk,cTransitionMapItem,parent)
					if mp.stopAt>0 && mp.stopAt==mp.lastSentPreStreamMapKey {
						mp.stopIt()
					}
				}
			}
		}
	}
	return sent
}
// Iter iterates over the items in the Multiplexer
// Each item is sent over a channel, so that
// we can iterate over the it using the builtin range keyword
func (mp *MsgMultiplexer) Iter() chan interface{} {
	return Bridge(mp.doneHandler.Done(), mp.stream)
}
//PrintPreStreamMap - prints the preStreamMap
func (mp *MsgMultiplexer) PrintPreStreamMap() {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	mp.preStreamMap.Sort()
	fmt.Printf("preStreamMap:%v\n", mp.preStreamMap)
	l := mp.preStreamMap.Len()
	if l > 0 {
		for it := range mp.preStreamMap.Iter() {
			cKey:= it.Key.(int64)
			mop := it.Value.(*SortedMap)
			fmt.Printf("cKey:%v - mop.len:%v\n", cKey,mop.Len())
		}
	}
}
func (mp *MsgMultiplexer) PrintDebug() {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	fmt.Printf("debug:\n%v\n", mp.debug)
}
func (mp *MsgMultiplexer) String() string {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	return "Index: " + mp.toStringIndex + "ID: " + mp.ID()
}

/******************************************************************************
Plug and Play and Transformation functions
*******************************************************************************/

// defaultGetItemKeyFn - Gets the key value used by the comparator in order to group the Messages in the
// SortedMap. Can be overridden for a more generic implementation that is not Message Oriented
func defaultGetItemKeyFn(v interface{}) int64 {
	return v.(*Message).CorrelationKey
}

// MsgMultiplexerMessagePairGetItemKeyFn - Gets the key value used by the comparator in order to group the Messages in the
// SortedMap. Can be overridden for a more generic implementation that is not Message Oriented
func MsgMultiplexerMessagePairGetItemKeyFn(v interface{}) int64 {
	return v.(*MessagePair).CorrelationKey
}

// defaultTransformFn - Transforms the SortedMap output into an Message for future consumption as part of the
// output channel of the MsgMultiplexer. Can be overridden for a more generic implementation
func defaultTransformFn(mp *MsgMultiplexer, sm *SortedMap, correlationKey int64) interface{} {

	re := NewMessage(sm,
		mp.MsgType,
		MessageWithIndex(mp.index),
		MessageWithCorrelationKey(correlationKey),
	)

	return re
}

// MsgMultiplexerMessagePairTransformFn - Transforms the SortedMap output into an Message for future consumption as part of the
// output channel of the MsgMultiplexer. Can be overridden for a more generic implementation
func MsgMultiplexerMessagePairTransformFn(mp *MsgMultiplexer, sm *SortedMap, correlationKey int64) interface{} {

	inMsg, _ := sm.Get(mp.minIndexKey)

	in := inMsg.(*MessagePair).In

	out := NewMessage(sm,
		mp.MsgType,
		MessageWithIndex(mp.index),
		MessageWithCorrelationKey(correlationKey),
	)

	re := NewMessagePair(in, out)
	return re
}