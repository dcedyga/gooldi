package concurrency_test

import (
	"fmt"

	concurrency "github.com/dcedyga/gooldi/concurrency"

	"time"
)

func (suite *Suite) Test07MultiMsgMultiplexer01WaitForAll() {
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
	mp := concurrency.NewMultiMsgMultiplexer(dh3,
		"MultiMsg",
		concurrency.MultiMsgMultiplexerWaitForAll(true),
		concurrency.MultiMsgMultiplexerIndex(1),
		concurrency.MultiMsgMultiplexerItemKeyFn(multiMsgGetItemKey),
		concurrency.MultiMsgMultiplexerTransformFn(multiMsgTransformFn),
	)

	//Create Processors

	mp.Set(b1.MsgType, b1.AddListener(dh))
	mp.Set(b2.MsgType, b2.AddListener(dh1))
	mp.Set(b3.MsgType, b3.AddListener(dh2))

	fmt.Printf("multiMsgMultiplexer - ID:%v, Index: %v\n", mp.ID(), mp.Index())
	bc, _ := mp.Get("bcaster1")
	fmt.Printf("multiMsgMultiplexer.bcaster1:%v\n", bc)

	//time.Sleep(300 * time.Microsecond)
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
		i := 0
		exit := false
		for msgId := 0; msgId < 30; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b1", msgId),
					b3.MsgType,
				)
				b1.Broadcast(e)
				i++
			}
			if exit {
				break
			}
			time.Sleep(500 * time.Microsecond)
		}
		fmt.Printf("Total Messages Broadcasted - bcaster1: %v,\n", i)
	}()
	go func() {
		i := 0
		exit := false
		for msgId := 0; msgId < 30; msgId++ {
			select {
			case <-dh1.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b2", msgId),
					b3.MsgType,
				)
				b2.Broadcast(e)
				i++
			}
			if exit {
				break
			}
			time.Sleep(1000 * time.Microsecond)
		}
		fmt.Printf("Total Messages Broadcasted - bcaster2: %v,\n", i)
	}()
	go func() {
		i := 0
		exit := false
		for msgId := 0; msgId < 30; msgId++ {

			select {
			case <-dh2.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b3", msgId),
					b3.MsgType,
				)
				b3.Broadcast(e)
				i++
			}
			if exit {
				break
			}
			time.Sleep(500 * time.Microsecond)
		}
		fmt.Printf("Total Messages Broadcasted - bcaster3: %v,\n", i)
	}()
	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test07MultiMsgMultiplexer02WaitForAllBuffer() {
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

	//Create Processors

	mp.Set(b1.MsgType, b1.AddListener(dh))
	mp.Set(b2.MsgType, b2.AddListener(dh1))
	mp.Set(b3.MsgType, b3.AddListener(dh2))

	//time.Sleep(300 * time.Microsecond)
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
		i := 0
		exit := false
		for msgId := 0; msgId < 30; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b1", msgId),
					b3.MsgType,
				)
				b1.Broadcast(e)
				i++
			}
			if exit {
				break
			}
			time.Sleep(500 * time.Microsecond)
		}
		fmt.Printf("Total Messages Broadcasted - bcaster1: %v,\n", i)
	}()
	go func() {
		i := 0
		exit := false
		for msgId := 0; msgId < 30; msgId++ {
			select {
			case <-dh1.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b2", msgId),
					b3.MsgType,
				)
				b2.Broadcast(e)
				i++
			}
			if exit {
				break
			}
			time.Sleep(1000 * time.Microsecond)
		}
		fmt.Printf("Total Messages Broadcasted - bcaster2: %v,\n", i)
	}()
	go func() {
		i := 0
		exit := false
		for msgId := 0; msgId < 30; msgId++ {

			select {
			case <-dh2.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b3", msgId),
					b3.MsgType,
				)
				b3.Broadcast(e)
				i++
			}
			if exit {
				break
			}
			time.Sleep(500 * time.Microsecond)
		}
		fmt.Printf("Total Messages Broadcasted - bcaster3: %v,\n", i)
	}()

	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test07MultiMsgMultiplexer03WaitForAllBM() {
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

	//Create Processors

	mp.Set(b1.MsgType, b1.AddListener(dh))
	mp.Set(b2.MsgType, b2.AddListener(dh1))
	mp.Set(b3.MsgType, b3.AddListener(dh2))

	//time.Sleep(300 * time.Microsecond)
	// Get all multiplexed messages

	printMultMsgMultiplexedResult(mp)

	// Broadcast Messages

	//wd, _ := mp.Get(1)
	go func() {
		i := 0
		exit := false
		for msgId := 0; msgId < 300000; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b1", msgId),
					b3.MsgType,
				)
				b1.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted - bcaster1: %v,\n", i)
	}()
	go func() {
		i := 0
		exit := false
		for msgId := 0; msgId < 300000; msgId++ {
			select {
			case <-dh1.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b2", msgId),
					b3.MsgType,
				)
				b2.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted - bcaster2: %v,\n", i)
	}()
	go func() {
		i := 0
		exit := false
		for msgId := 0; msgId < 300000; msgId++ {

			select {
			case <-dh2.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b3", msgId),
					b3.MsgType,
				)
				b3.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted - bcaster3: %v,\n", i)
	}()

	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test07MultiMsgMultiplexer04WaitForAllWithTimer() {
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

	//Create Processors

	mp.Set(b1.MsgType, b1.AddListener(dh))
	mp.Set(b2.MsgType, b2.AddListener(dh1))
	mp.Set(b3.MsgType, b3.AddListener(dh2))

	//time.Sleep(300 * time.Microsecond)
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
		i := 0
		exit := false
		for msgId := 0; msgId < 300; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b1", msgId),
					b3.MsgType,
				)
				b1.Broadcast(e)
				i++
			}
			if exit {
				break
			}
			time.Sleep(500 * time.Microsecond)
		}
		fmt.Printf("Total Messages Broadcasted - bcaster1: %v,\n", i)
	}()
	go func() {
		i := 0
		exit := false
		for msgId := 0; msgId < 300; msgId++ {
			select {
			case <-dh1.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b2", msgId),
					b3.MsgType,
				)
				b2.Broadcast(e)
				i++
			}
			if exit {
				break
			}
			time.Sleep(1000 * time.Microsecond)
		}
		fmt.Printf("Total Messages Broadcasted - bcaster2: %v,\n", i)
	}()
	go func() {
		i := 0
		exit := false
		for msgId := 0; msgId < 300; msgId++ {

			select {
			case <-dh2.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b3", msgId),
					b3.MsgType,
				)
				b3.Broadcast(e)
				i++
			}
			if exit {
				break
			}
			time.Sleep(500 * time.Microsecond)
		}
		fmt.Printf("Total Messages Broadcasted - bcaster3: %v,\n", i)
	}()
	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test07MultiMsgMultiplexer05WaitForAllWithTimerBM() {
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

	//Create Processors

	mp.Set(b1.MsgType, b1.AddListener(dh))
	mp.Set(b2.MsgType, b2.AddListener(dh1))
	mp.Set(b3.MsgType, b3.AddListener(dh2))

	//time.Sleep(300 * time.Microsecond)
	// Get all multiplexed messages

	printMultMsgMultiplexedResult(mp)
	// Broadcast Messages

	//wd, _ := mp.Get(1)
	go func() {
		i := 0
		exit := false
		for msgId := 0; msgId < 500000; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b1", msgId),
					b3.MsgType,
				)
				b1.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted - bcaster1: %v,\n", i)
	}()
	go func() {
		i := 0
		exit := false
		for msgId := 0; msgId < 500000; msgId++ {
			select {
			case <-dh1.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b2", msgId),
					b3.MsgType,
				)
				b2.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted - bcaster2: %v,\n", i)
	}()
	go func() {
		i := 0
		exit := false
		for msgId := 0; msgId < 500000; msgId++ {

			select {
			case <-dh2.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					fmt.Sprintf("%v From b3", msgId),
					b3.MsgType,
				)
				b3.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted - bcaster3: %v,\n", i)
	}()
	time.Sleep(1000 * time.Millisecond)

}
