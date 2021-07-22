package concurrency_test

import (
	"fmt"

	concurrency "github.com/dcedyga/gooldi/concurrency"

	"time"
)

func (suite *Suite) Test04MsgMultiplexer01FilteredProcessors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	defer dm.GetDoneFunc()()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
		concurrency.MsgMultiplexerTransformFn(concurrency.MsgMultiplexerMessagePairTransformFn),
		concurrency.MsgMultiplexerGetItemKeyFn(concurrency.MsgMultiplexerMessagePairGetItemKeyFn),
	)
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	w := createMessagePairProcessors(mp, b, dh1, 4)
	//Create Filter
	filter := concurrency.NewFilter(
		fmt.Sprintf("filter%v", 5),
		dh1,
		concurrency.FilterWithIndex(5),
		concurrency.FilterWithInputChannel(mp.Iter()),
	)
	go filter.Filter(filterFn)
	// Print output
	go func() {
		i := 0
		for v := range filter.OutputChannel() {
			_ = v
			item := v.(*concurrency.MessagePair)
			s := item.Out.Message.(*concurrency.SortedMap)
			fmt.Printf("-------\n")
			for pe := range s.Iter() {
				pr := pe.Value.(*concurrency.MessagePair).Out
				fmt.Printf("CorrelationKey(%v) - Index(%v) - TimeInNano(%v) - ID(%v) - MsgType(%v) - Message: %v \n", pr.CorrelationKey, pr.Index, pr.TimeInNano, pr.ID, pr.MsgType, pr.Message)
			}
			i++

		}
		fmt.Printf("Total Messages Filtered: %v,\n", i)
	}()

	// Broadcast Messages
	broadcastWithStop(b, 30, 10, 15, w.GetItemAtIndex(3).(*concurrency.Processor))

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test04MsgMultiplexer02Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	defer dm.GetDoneFunc()()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	w := createProcessors(mp, b, dh1, 4)
	// Print output
	go func() {
		i := 0
		for v := range mp.Iter() {
			_ = v
			item := v.(*concurrency.Message)
			s := item.Message.(*concurrency.SortedMap)
			fmt.Printf("-------\n")
			for pe := range s.Iter() {
				pr := pe.Value.(*concurrency.Message)
				fmt.Printf("CorrelationKey(%v) - Index(%v) - TimeInNano(%v) - ID(%v) - MsgType(%v) - Message: %v \n", pr.CorrelationKey, pr.Index, pr.TimeInNano, pr.ID, pr.MsgType, pr.Message)
			}
			i++

		}
		fmt.Printf("Total Messages Multiplexed: %v,\n", i)
	}()

	// Broadcast Messages
	broadcastWithStop(b, 30, 10, 15, w.GetItemAtIndex(3).(*concurrency.Processor))

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test04MsgMultiplexer03ProcessorsBM2Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 2)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer04ProcessorsBM4Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 4)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer05ProcessorsBM8Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 8)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer06ProcessorsBM16Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 16)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer07ProcessorsBM32Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 32)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer08ProcessorsBM64Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 64)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer09ProcessorsBM128Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 128)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer10ProcessorsBM512Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 512)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer11ProcessorsBM1024Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 1024)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer12ProcessorsBM2048Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 2048)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer13ProcessorsBM10000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 10000)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer14ProcessorsBM20000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 20000)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer15ProcessorsBM30000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 30000)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer16ProcessorsBM40000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 40000)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer17ProcessorsBM50000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 50000)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer18ProcessorsBM60000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 60000)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer19ProcessorsBM70000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 70000)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer20ProcessorsBM80000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 80000)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer21ProcessorsBM90000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
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
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 90000)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test04MsgMultiplexer22ProcessorsBM100000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(50 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType, concurrency.MsgMultiplexerIndex(1))
	defer func() {
		dm.GetDoneFunc()()
		mp.PrintPreStreamMap()
	}()
	//Create Processors
	createProcessors(mp, b, dh1, 100000)
	// Print Result goroutine
	printMultiplexedResult(mp)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
