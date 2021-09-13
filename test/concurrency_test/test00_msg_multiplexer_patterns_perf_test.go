package concurrency_test

import (
	concurrency "github.com/dcedyga/gooldi/concurrency"

	"time"
)

func (suite *Suite) Test00MsgMultiplexer01ProcessorsBM2Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 2)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer02ProcessorsBM4Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 4)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer03ProcessorsBM8Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 8)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer04ProcessorsBM16Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 16)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer05ProcessorsBM32Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 32)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer06ProcessorsBM64Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 64)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer07ProcessorsBM128Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 128)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer08ProcessorsBM512Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 512)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test00MsgMultiplexer09ProcessorsBM1024Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 1024)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test00MsgMultiplexer10ProcessorsBM2048Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 2048)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer11ProcessorsBM10000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 10000)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer12ProcessorsBM20000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 20000)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer13ProcessorsBM30000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 30000)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer14ProcessorsBM40000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 40000)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer15ProcessorsBM50000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 50000)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer16ProcessorsBM60000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 60000)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer17ProcessorsBM70000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 70000)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer18ProcessorsBM80000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 80000)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test00MsgMultiplexer19ProcessorsBM90000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(100 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessors(mp, b, dh1, 90000)
	//mp start
	mp.Start()
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
