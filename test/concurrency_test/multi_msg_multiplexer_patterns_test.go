package concurrency_test

import (
	"fmt"

	concurrency "github.com/dcedyga/gooldi/concurrency"

	"time"
)

func (suite *Suite) Test06MultiMsgMultiplexer01WaitForAll() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(3 * time.Millisecond))
	defer dm.GetDoneFunc()()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	dh3 := dm.AddNewDoneHandler(3)
	//seqStream := make(chan (<-chan interface{}))
	//Create caster
	b1 := concurrency.NewBCaster(dh,
		"bcaster1",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b2 := concurrency.NewBCaster(dh1,
		"bcaster2",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b3 := concurrency.NewBCaster(dh2,
		"bcaster3",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	//Create Multiplexer
	mp := concurrency.NewMultiMsgMultiplexer(dh3, "MultiMsg", concurrency.MultiMsgMultiplexerWaitForAll(true), concurrency.MultiMsgMultiplexerIndex(1))

	mp.Start()
	//Create Processors

	mp.Set(b1.MsgType, b1.AddListener(dh))
	mp.Set(b2.MsgType, b2.AddListener(dh1))
	mp.Set(b3.MsgType, b3.AddListener(dh2))

	time.Sleep(300 * time.Microsecond)
	// Get all multiplexed messages

	go func() {
		i := 0
		for v := range mp.Iter() {
			_ = v
			item := v.(*concurrency.Message)
			fmt.Printf("Multiplexer-key:%v\n", item.ID)

			m := item.Message.(*concurrency.SortedMap)
			//fmt.Printf("m:%v\n", m)
			for it := range m.Iter() {
				mapit := it.Value.(*concurrency.SortedMap)
				for mit := range mapit.Iter() {
					e := mit.Value.(*concurrency.Message)
					fmt.Printf("bcaster-key:%v - msgkey:%v - e.OutMessage:%v\n", it.Key, mit.Key, e.Message)
				}

			}

			i++

		}
		fmt.Printf("Total Messages Extracted: %v,\n", i)
	}()

	// Broadcast Messages

	//wd, _ := mp.Get(1)
	go func() {
		for msgId := 0; msgId < 30; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b1", msgId),
				b1.MsgType,
			)
			b1.Broadcast(e)
			time.Sleep(500 * time.Microsecond)
		}

	}()
	go func() {
		for msgId := 0; msgId < 30; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b2", msgId),
				b2.MsgType,
			)
			b2.Broadcast(e)
			time.Sleep(1 * time.Millisecond)
		}

	}()
	go func() {
		for msgId := 0; msgId < 30; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b3", msgId),
				b3.MsgType,
			)
			b3.Broadcast(e)
			time.Sleep(500 * time.Microsecond)
		}

	}()
	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test06MultiMsgMultiplexer02WaitForAllBuffer() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	defer dm.GetDoneFunc()()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	dh3 := dm.AddNewDoneHandler(3)
	//seqStream := make(chan (<-chan interface{}))
	//Create caster
	b1 := concurrency.NewBCaster(dh,
		"bcaster1",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b2 := concurrency.NewBCaster(dh1,
		"bcaster2",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b3 := concurrency.NewBCaster(dh2,
		"bcaster3",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	//Create Multiplexer
	mp := concurrency.NewMultiMsgMultiplexer(dh3, "MultiMsg", concurrency.MultiMsgMultiplexerBufferSize(2), concurrency.MultiMsgMultiplexerIndex(1))

	mp.Start()
	//Create Processors

	mp.Set(b1.MsgType, b1.AddListener(dh))
	mp.Set(b2.MsgType, b2.AddListener(dh1))
	mp.Set(b3.MsgType, b3.AddListener(dh2))

	time.Sleep(300 * time.Microsecond)
	// Get all multiplexed messages

	go func() {
		i := 0
		for v := range mp.Iter() {
			_ = v
			item := v.(*concurrency.Message)
			fmt.Printf("Multiplexer-key:%v\n", item.ID)

			m := item.Message.(*concurrency.SortedMap)
			//fmt.Printf("m:%v\n", m)
			for it := range m.Iter() {
				mapit := it.Value.(*concurrency.SortedMap)
				fmt.Printf("bcaster-key:%v\n", it.Key)
				for mit := range mapit.Iter() {
					e := mit.Value.(*concurrency.Message)
					fmt.Printf("msgkey:%v - e.OutMessage:%v\n", mit.Key, e.Message)
				}

			}

			i++

		}
		fmt.Printf("Total Messages Extracted: %v,\n", i)
	}()

	// Broadcast Messages

	//wd, _ := mp.Get(1)
	go func() {
		for msgId := 0; msgId < 30; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b1", msgId),
				b1.MsgType,
			)
			b1.Broadcast(e)
			time.Sleep(500 * time.Microsecond)
		}

	}()
	go func() {
		for msgId := 0; msgId < 30; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b2", msgId),
				b2.MsgType,
			)
			b2.Broadcast(e)
			time.Sleep(1 * time.Millisecond)
		}

	}()
	go func() {
		for msgId := 0; msgId < 30; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b3", msgId),
				b3.MsgType,
			)
			b3.Broadcast(e)
			time.Sleep(500 * time.Microsecond)
		}

	}()

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test06MultiMsgMultiplexer03WaitForAllBM() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	defer dm.GetDoneFunc()()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	dh3 := dm.AddNewDoneHandler(3)
	//seqStream := make(chan (<-chan interface{}))
	//Create caster
	b1 := concurrency.NewBCaster(dh,
		"bcaster1",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b2 := concurrency.NewBCaster(dh1,
		"bcaster2",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b3 := concurrency.NewBCaster(dh2,
		"bcaster3",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	//Create Multiplexer
	mp := concurrency.NewMultiMsgMultiplexer(dh3, "MultiMsg", concurrency.MultiMsgMultiplexerWaitForAll(true), concurrency.MultiMsgMultiplexerIndex(1))

	mp.Start()
	//Create Processors

	mp.Set(b1.MsgType, b1.AddListener(dh))
	mp.Set(b2.MsgType, b2.AddListener(dh1))
	mp.Set(b3.MsgType, b3.AddListener(dh2))

	time.Sleep(300 * time.Microsecond)
	// Get all multiplexed messages

	go func() {
		i := 0
		for v := range mp.Iter() {
			_ = v

			i++

		}
		fmt.Printf("Total Messages Extracted: %v,\n", i)
	}()

	// Broadcast Messages

	//wd, _ := mp.Get(1)
	go func() {
		for msgId := 0; msgId < 300000; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b1", msgId),
				b1.MsgType,
			)
			b1.Broadcast(e)

		}

	}()
	go func() {
		for msgId := 0; msgId < 300000; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b2", msgId),
				b2.MsgType,
			)
			b2.Broadcast(e)

		}

	}()
	go func() {
		for msgId := 0; msgId < 300000; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b3", msgId),
				b3.MsgType,
			)
			b3.Broadcast(e)

		}

	}()

	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test06MultiMsgMultiplexer04WaitForAllWithTimer() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	defer dm.GetDoneFunc()()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	dh3 := dm.AddNewDoneHandler(3)
	//seqStream := make(chan (<-chan interface{}))
	//Create caster
	b1 := concurrency.NewBCaster(dh,
		"bcaster1",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b2 := concurrency.NewBCaster(dh1,
		"bcaster2",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b3 := concurrency.NewBCaster(dh2,
		"bcaster3",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	d := 200 * time.Millisecond
	//Create Multiplexer
	mp := concurrency.NewMultiMsgMultiplexer(dh3,
		"MultiMsg",
		concurrency.MultiMsgMultiplexerWaitForAll(true),
		concurrency.MultiMsgMultiplexerIndex(1),
		concurrency.MultiMsgMultiplexerSendPeriod(&d),
	)

	mp.Start()
	//Create Processors

	mp.Set(b1.MsgType, b1.AddListener(dh))
	mp.Set(b2.MsgType, b2.AddListener(dh1))
	mp.Set(b3.MsgType, b3.AddListener(dh2))

	time.Sleep(300 * time.Microsecond)
	// Get all multiplexed messages

	go func() {
		i := 0
		for v := range mp.Iter() {
			_ = v
			item := v.(*concurrency.Message)
			fmt.Printf("Multiplexer-key:%v\n", item.ID)

			m := item.Message.(*concurrency.SortedMap)
			//fmt.Printf("m:%v\n", m)
			for it := range m.Iter() {
				mapit := it.Value.(*concurrency.SortedMap)
				for mit := range mapit.Iter() {
					e := mit.Value.(*concurrency.Message)
					fmt.Printf("bcaster-key:%v - msgkey:%v - e.OutMessage:%v\n", it.Key, mit.Key, e.Message)
				}

			}

			i++

		}
		fmt.Printf("Total Messages Extracted: %v,\n", i)
	}()

	// Broadcast Messages

	//wd, _ := mp.Get(1)
	go func() {
		for msgId := 0; msgId < 300; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b1", msgId),
				b1.MsgType,
			)
			b1.Broadcast(e)
			time.Sleep(500 * time.Microsecond)
		}

	}()
	go func() {
		for msgId := 0; msgId < 300; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b2", msgId),
				b2.MsgType,
			)
			b2.Broadcast(e)
			time.Sleep(1 * time.Millisecond)
		}

	}()
	go func() {
		for msgId := 0; msgId < 300; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b3", msgId),
				b3.MsgType,
			)
			b3.Broadcast(e)
			time.Sleep(500 * time.Microsecond)
		}

	}()
	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test06MultiMsgMultiplexer05WaitForAllWithTimerBM() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	defer dm.GetDoneFunc()()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	dh3 := dm.AddNewDoneHandler(3)
	//seqStream := make(chan (<-chan interface{}))
	//Create caster
	b1 := concurrency.NewBCaster(dh,
		"bcaster1",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b2 := concurrency.NewBCaster(dh1,
		"bcaster2",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b3 := concurrency.NewBCaster(dh2,
		"bcaster3",
		//concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	d := 1 * time.Microsecond
	//Create Multiplexer
	mp := concurrency.NewMultiMsgMultiplexer(dh3,
		"MultiMsg",
		concurrency.MultiMsgMultiplexerWaitForAll(true),
		concurrency.MultiMsgMultiplexerIndex(1),
		concurrency.MultiMsgMultiplexerSendPeriod(&d),
	)

	mp.Start()
	//Create Processors

	mp.Set(b1.MsgType, b1.AddListener(dh))
	mp.Set(b2.MsgType, b2.AddListener(dh1))
	mp.Set(b3.MsgType, b3.AddListener(dh2))

	time.Sleep(300 * time.Microsecond)
	// Get all multiplexed messages

	go func() {
		i := 0
		for v := range mp.Iter() {
			_ = v

			i++

		}
		fmt.Printf("Total Messages Extracted: %v,\n", i)
	}()

	// Broadcast Messages

	//wd, _ := mp.Get(1)
	go func() {
		for msgId := 0; msgId < 300000; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b1", msgId),
				b1.MsgType,
			)
			b1.Broadcast(e)

		}

	}()
	go func() {
		for msgId := 0; msgId < 300000; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b2", msgId),
				b2.MsgType,
			)
			b2.Broadcast(e)

		}

	}()
	go func() {
		for msgId := 0; msgId < 300000; msgId++ {
			e := concurrency.NewMessage(
				fmt.Sprintf("%v From b3", msgId),
				b3.MsgType,
			)
			b3.Broadcast(e)

		}

	}()
	time.Sleep(1000 * time.Millisecond)

}
