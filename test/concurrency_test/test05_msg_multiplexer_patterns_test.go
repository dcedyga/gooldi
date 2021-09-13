package concurrency_test

import (
	"fmt"

	concurrency "github.com/dcedyga/gooldi/concurrency"

	"time"
)

func (suite *Suite) Test05MsgMultiplexer01FilteredProcessors() {
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
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w := createMessagePairProcessors(mp, b, dh1, 4)

	fmt.Printf("MsgMultiplexer - ID:%v, Index: %v\n", mp.ID(), mp.Index())
	wc, _ := mp.Get(w.GetItemAtIndex(0).(*concurrency.Processor).Index())
	fmt.Printf("MsgMultiplexer.bcaster1:%v\n", wc)
	//Create Filter
	filter := concurrency.NewFilter(
		fmt.Sprintf("filter%v", 5),
		dh1,
		concurrency.FilterWithIndex(5),
		concurrency.FilterWithInputChannel(mp.Iter()),
		concurrency.FilterTransformFn(filterTransformFn),
	)

	fmt.Printf("filter - ID: %v, Index: %v, InputChannel:%v, HasValidInputChan:%v\n", filter.ID(), filter.Index(), filter.InputChannel(), filter.HasValidInputChan())

	go filter.Filter(filterFn)
	filter.Stop()
	filter.Start()
	//mp start
	mp.Start()
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
func (suite *Suite) Test05MsgMultiplexer02Processors() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w := createProcessors(mp, b, dh1, 4)
	//mp start
	mp.Start()
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

	processor1 := w.GetItemAtIndex(3).(*concurrency.Processor)
	// Broadcast Messages
	broadcastWithStop(b, 30, 10, 15, processor1)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer03DeleteProcessor() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w := createProcessors(mp, b, dh1, 4)
	//mp start
	mp.Start()
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

	processor1 := w.GetItemAtIndex(2).(*concurrency.Processor)
	// Broadcast Messages
	broadcastWithDeleteAfterCKey(b, 30, 25, mp, processor1)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer04Delete3Processors() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w := createProcessors(mp, b, dh1, 8)
	//mp start
	mp.Start()
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

	processor1 := w.GetItemAtIndex(2).(*concurrency.Processor)
	processor2 := w.GetItemAtIndex(4).(*concurrency.Processor)
	processor3 := w.GetItemAtIndex(6).(*concurrency.Processor)
	// Broadcast Messages
	broadcastWithDeleteAfterCKey3(b, 30, 5,12,25, mp, processor1,processor2,processor3)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer05DeleteDecoupledProcessors() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w := createProcessors(mp, b, dh1, 16)
	//mp start
	mp.Start()
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

	processor1 := w.GetItemAtIndex(1).(*concurrency.Processor)
	processor3 := w.GetItemAtIndex(3).(*concurrency.Processor)
	processor5 := w.GetItemAtIndex(5).(*concurrency.Processor)
	processor7 := w.GetItemAtIndex(7).(*concurrency.Processor)
	processor9 := w.GetItemAtIndex(9).(*concurrency.Processor)
	processor11 := w.GetItemAtIndex(11).(*concurrency.Processor)
	processor13 := w.GetItemAtIndex(13).(*concurrency.Processor)
	go func() {
		time.Sleep(500 * time.Microsecond)
		mp.Delete(processor1.Index())
		//time.Sleep(1000 * time.Microsecond)
		mp.Delete(processor3.Index())
		// time.Sleep(100 * time.Microsecond)
		mp.Delete(processor5.Index())
		// time.Sleep(100 * time.Microsecond)
		mp.Delete(processor7.Index())
		// time.Sleep(100 * time.Microsecond)
		mp.Delete(processor9.Index())
		// time.Sleep(100 * time.Microsecond)
		mp.Delete(processor11.Index())
		// time.Sleep(100 * time.Microsecond)
		mp.Delete(processor13.Index())
	}()
	// Broadcast Messages
	broadcast(b, 30)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer06DeleteDecoupledSlowProcessors() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w := createProcessors(mp, b, dh1, 16)
	ws := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 16),
		dh1,
		concurrency.ProcessorWithIndex(int64(16)),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
	)
	mp.Set(ws.Index(), ws.OutputChannel())
	go ws.Process(slowprocess9, 16)
	//mp start
	mp.Start()
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

	processor1 := w.GetItemAtIndex(1).(*concurrency.Processor)
	processor3 := w.GetItemAtIndex(3).(*concurrency.Processor)
	processor5 := w.GetItemAtIndex(5).(*concurrency.Processor)
	processor7 := w.GetItemAtIndex(7).(*concurrency.Processor)
	processor9 := w.GetItemAtIndex(9).(*concurrency.Processor)
	processor11 := w.GetItemAtIndex(11).(*concurrency.Processor)
	processor13 := w.GetItemAtIndex(13).(*concurrency.Processor)
	go func() {
		time.Sleep(500 * time.Microsecond)
		mp.Delete(processor1.Index())
		//time.Sleep(1000 * time.Microsecond)
		mp.Delete(processor3.Index())
		// time.Sleep(100 * time.Microsecond)
		mp.Delete(processor5.Index())
		// time.Sleep(100 * time.Microsecond)
		mp.Delete(processor7.Index())
		// time.Sleep(100 * time.Microsecond)
		mp.Delete(processor9.Index())
		// time.Sleep(100 * time.Microsecond)
		mp.Delete(processor11.Index())
		// time.Sleep(100 * time.Microsecond)
		mp.Delete(processor13.Index())
		time.Sleep(10 * time.Millisecond)
		mp.Delete(ws.Index())
	}()
	// Broadcast Messages
	broadcast(b, 30)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer07AddDecoupledProcessor() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessorsFrom(mp, b, dh1, 0, 4)
	//mp start
	mp.Start()
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
	go func() {
		time.Sleep(500 * time.Microsecond)
		createProcessorsFrom(mp, b, dh1, 4, 1)
	}()
	// Broadcast Messages
	broadcast(b, 20)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer08Add4DecoupledProcessor() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessorsFrom(mp, b, dh1, 0, 4)
	//mp start
	mp.Start()
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
	go func() {
		time.Sleep(500 * time.Microsecond)
		createProcessorsFrom(mp, b, dh1, 4, 3)
	}()
	// Broadcast Messages
	broadcast(b, 30)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer09AddAndDeleteProcessor() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w := createProcessors(mp, b, dh1, 4)
	//mp start
	mp.Start()
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

	processor1 := w.GetItemAtIndex(2).(*concurrency.Processor)
	// Broadcast Messages
	broadcastWithAddDelete(b, 30, 25, mp, processor1,dh1)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer10AddDeleteDecoupled() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w := createProcessors(mp, b, dh1, 8)
	//mp start
	mp.Start()
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

	processor1 := w.GetItemAtIndex(1).(*concurrency.Processor)
	processor3 := w.GetItemAtIndex(3).(*concurrency.Processor)
	processor5 := w.GetItemAtIndex(5).(*concurrency.Processor)
	
	go func() {
		time.Sleep(500 * time.Microsecond)
		mp.Delete(processor1.Index())
		//time.Sleep(1000 * time.Microsecond)
		createProcessorsFrom(mp, b, dh1, 8, 1)
		mp.Delete(processor3.Index())
		// time.Sleep(100 * time.Microsecond)
		createProcessorsFrom(mp, b, dh1, 9, 1)
		mp.Delete(processor5.Index())
		// time.Sleep(100 * time.Microsecond)
		createProcessorsFrom(mp, b, dh1, 10, 1)
	}()

	// Broadcast Messages
	broadcast(b, 30)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer11AddDeleteDecoupledWithTimeDiff() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w := createProcessors(mp, b, dh1, 8)
	//mp start
	mp.Start()
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

	processor1 := w.GetItemAtIndex(1).(*concurrency.Processor)
	processor3 := w.GetItemAtIndex(3).(*concurrency.Processor)
	processor5 := w.GetItemAtIndex(5).(*concurrency.Processor)
	
	go func() {
		time.Sleep(500 * time.Microsecond)
		mp.Delete(processor1.Index())
		time.Sleep(200 * time.Microsecond)
		createProcessorsFrom(mp, b, dh1, 8, 1)
		time.Sleep(200 * time.Microsecond)
		mp.Delete(processor3.Index())
		time.Sleep(200 * time.Microsecond)
		createProcessorsFrom(mp, b, dh1, 9, 1)
		time.Sleep(200 * time.Microsecond)
		mp.Delete(processor5.Index())
		time.Sleep(200 * time.Microsecond)
		createProcessorsFrom(mp, b, dh1, 10, 1)
	}()

	// Broadcast Messages
	broadcast(b, 50)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer12AddAtCorrelationKey() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	//w := createProcessors(mp, b, dh1, 4)
	createProcessors(mp, b, dh1, 4)
	//mp start
	mp.Start()
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
	broadcastWithAddAtCKey(b, 20, 18,4,dh1, mp) 

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer13DeleteAddAtSameCorrelationKey() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w := createProcessors(mp, b, dh1, 4)
	//mp start
	mp.Start()
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
	// processor1 := concurrency.NewProcessor(
	// 	fmt.Sprintf("worker%v", 4),
	// 	dh1,
	// 	concurrency.ProcessorWithIndex(int64(4)),
	// 	concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
	// 	//concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
	// )
	// Broadcast Messages
	processor1 := w.GetItemAtIndex(1).(*concurrency.Processor)

	broadcastWithDeleteAddAtCKeys(b, 20, 18,18,4,dh1, mp,processor1) 

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer14DeleteAddAtDiffCorrelationKey() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w := createProcessors(mp, b, dh1, 4)
	//mp start
	mp.Start()
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
	// processor1 := concurrency.NewProcessor(
	// 	fmt.Sprintf("worker%v", 4),
	// 	dh1,
	// 	concurrency.ProcessorWithIndex(int64(4)),
	// 	concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
	// 	//concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
	// )
	// Broadcast Messages
	processor1 := w.GetItemAtIndex(1).(*concurrency.Processor)

	broadcastWithDeleteAddAtCKeys(b, 20, 17,18,4,dh1, mp,processor1) 

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer15StopStart() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessorsFrom(mp, b, dh1, 0, 4)
	//mp start
	mp.Start()
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
	go func() {
		time.Sleep(300 * time.Microsecond)
		mp.Stop()
		time.Sleep(200 * time.Microsecond)
		mp.Start()
	}()
	// Broadcast Messages
	broadcast(b, 30)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer16StopDeleteStart() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w :=createProcessorsFrom(mp, b, dh1, 0, 4)
	//mp start
	mp.Start()
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
	processor1 := w.GetItemAtIndex(2).(*concurrency.Processor)
	go func() {
		time.Sleep(200 * time.Microsecond)
		mp.Stop()
		time.Sleep(200 * time.Microsecond)
		mp.Delete(processor1.Index())
		time.Sleep(200 * time.Microsecond)
		mp.Start()
	}()
	// Broadcast Messages
	broadcast(b, 40)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer17StopAddStart() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessorsFrom(mp, b, dh1, 0, 4)
	//mp start
	mp.Start()
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
	//processor1 := w.GetItemAtIndex(2).(*concurrency.Processor)
	go func() {
		time.Sleep(200 * time.Microsecond)
		mp.Stop()
		time.Sleep(200 * time.Microsecond)
		createProcessorsFrom(mp, b, dh, 4, 1)
		time.Sleep(200 * time.Microsecond)
		mp.Start()
	}()
	// Broadcast Messages
	broadcast(b, 50)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer18StopDeleteAddStart() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w :=createProcessorsFrom(mp, b, dh1, 0, 4)
	//mp start
	mp.Start()
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
	processor1 := w.GetItemAtIndex(2).(*concurrency.Processor)
	go func() {
		time.Sleep(200 * time.Microsecond)
		mp.Stop()
		time.Sleep(200 * time.Microsecond)
		mp.Delete(processor1.Index())
		createProcessorsFrom(mp, b, dh, 4, 1)
		time.Sleep(200 * time.Microsecond)
		mp.Start()
	}()
	// Broadcast Messages
	broadcast(b, 50)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer19StopAfterCorrelationnKey() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessorsFrom(mp, b, dh1, 0, 4)
	//mp start
	mp.Start()
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
	broadcastWitStopAfterCKeys(b, 10,5, mp)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer20StopStartAfterCorrelationnKeys() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessorsFrom(mp, b, dh1, 0, 4)
	//mp start
	mp.Start()
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
	broadcastWitStopStartAfterCKeys(b, 20,4,9, mp)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer21StartAfterCorrelationnKey() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessorsFrom(mp, b, dh1, 0, 4)
	//mp start
	mp.InitializeProcessors()
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
	broadcastWitStartAfterCKeys(b, 10,4, mp)

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer22StopAddDeleteStartAfterCorrelationnKeys() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w:=createProcessorsFrom(mp, b, dh1, 0, 4)
	//mp start
	mp.Start()
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
	processor1 := w.GetItemAtIndex(2).(*concurrency.Processor)
	// Broadcast Messages
	broadcastWitStopAddDeleteStartAfterCKeys(b, 20,4, 9,4,dh1,mp,processor1) 
	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer23StopDeleteStartAfterCorrelationnKeys() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	w:=createProcessorsFrom(mp, b, dh1, 0, 4)
	//mp start
	mp.Start()
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
	processor1 := w.GetItemAtIndex(2).(*concurrency.Processor)
	// Broadcast Messages
	broadcastWitStopDeleteStartAfterCKeys(b, 20,4, 9,mp,processor1) 
	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MsgMultiplexer24StopAddStartAfterCorrelationnKeys() {
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
	mp := concurrency.NewMsgMultiplexer(dh2,
		b.MsgType, concurrency.MsgMultiplexerIndex(1),
	)
	defer doneWithBCaster(dm, mp, b)
	//Create Processors
	createProcessorsFrom(mp, b, dh1, 0, 4)
	//mp start
	mp.Start()
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
	broadcastWitStopAddStartAfterCKeys(b, 20,4, 9,4,dh1,mp) 
	time.Sleep(1000 * time.Millisecond)

}