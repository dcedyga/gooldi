package concurrency_test

import (
	"fmt"
	concurrency "gool/concurrency"

	"time"

	uuid "github.com/satori/go.uuid"
)

func (suite *Suite) Test05MultiMsgMultiplexer01WaitForAll() {
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
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b2 := concurrency.NewBCaster(dh1,
		"bcaster2",
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b3 := concurrency.NewBCaster(dh2,
		"bcaster3",
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	//Create Multiplexer
	mp := concurrency.NewMultiMsgMultiplexer(dh3, "MultiMsg", concurrency.MultiMsgMultiplexerWaitForAll(true), concurrency.MultiMsgMultiplexerSequence(1))

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
			item := v.(*concurrency.Event)
			fmt.Printf("Multiplexer-key:%v\n", item.OutMessage.ID)

			m := item.OutMessage.Message.(*concurrency.SortedMap)
			//fmt.Printf("m:%v\n", m)
			for it := range m.Iter() {
				mapit := it.Value.(*concurrency.SortedMap)
				for mit := range mapit.Iter() {
					e := mit.Value.(*concurrency.Event)
					fmt.Printf("bcaster-key:%v - msgkey:%v - e.InitMessage:%v - e.OutMessage:%v\n", it.Key, mit.Key, e.InitMessage.Message, e.OutMessage.Message)
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

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b1", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b1.MsgType}

			b1.Broadcast(e)
			time.Sleep(500 * time.Microsecond)
		}

	}()
	go func() {
		for msgId := 0; msgId < 30; msgId++ {

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b2", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b2.MsgType}

			b2.Broadcast(e)
			time.Sleep(1 * time.Millisecond)
		}

	}()
	go func() {
		for msgId := 0; msgId < 30; msgId++ {

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b3", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b3.MsgType}

			b3.Broadcast(e)
			time.Sleep(500 * time.Microsecond)
		}

	}()
	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test05MultiMsgMultiplexer02WaitForAllBuffer() {
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
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b2 := concurrency.NewBCaster(dh1,
		"bcaster2",
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b3 := concurrency.NewBCaster(dh2,
		"bcaster3",
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	//Create Multiplexer
	mp := concurrency.NewMultiMsgMultiplexer(dh3, "MultiMsg", concurrency.MultiMsgMultiplexerBufferSize(2), concurrency.MultiMsgMultiplexerSequence(1))

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
			item := v.(*concurrency.Event)
			fmt.Printf("Multiplexer-key:%v\n", item.OutMessage.ID)

			m := item.OutMessage.Message.(*concurrency.SortedMap)
			//fmt.Printf("m:%v\n", m)
			for it := range m.Iter() {
				mapit := it.Value.(*concurrency.SortedMap)
				fmt.Printf("bcaster-key:%v\n", it.Key)
				for mit := range mapit.Iter() {
					e := mit.Value.(*concurrency.Event)
					fmt.Printf("msgkey:%v - e.InitMessage:%v - e.OutMessage:%v\n", mit.Key, e.InitMessage.Message, e.OutMessage.Message)
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

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b1", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b1.MsgType}

			b1.Broadcast(e)
			time.Sleep(500 * time.Microsecond)
		}

	}()
	go func() {
		for msgId := 0; msgId < 30; msgId++ {

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b2", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b2.MsgType}

			b2.Broadcast(e)
			time.Sleep(1 * time.Millisecond)
		}

	}()
	go func() {
		for msgId := 0; msgId < 30; msgId++ {

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b3", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b3.MsgType}

			b3.Broadcast(e)
			time.Sleep(500 * time.Microsecond)
		}

	}()
	time.Sleep(1000 * time.Millisecond)

}
func (suite *Suite) Test05MultiMsgMultiplexer03WaitForAllBM() {
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
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b2 := concurrency.NewBCaster(dh1,
		"bcaster2",
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b3 := concurrency.NewBCaster(dh2,
		"bcaster3",
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	//Create Multiplexer
	mp := concurrency.NewMultiMsgMultiplexer(dh3, "MultiMsg", concurrency.MultiMsgMultiplexerWaitForAll(true), concurrency.MultiMsgMultiplexerSequence(1))

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

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b1", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b1.MsgType}

			b1.Broadcast(e)

		}

	}()
	go func() {
		for msgId := 0; msgId < 300000; msgId++ {

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b2", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b2.MsgType}

			b2.Broadcast(e)

		}

	}()
	go func() {
		for msgId := 0; msgId < 300000; msgId++ {

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b3", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b3.MsgType}

			b3.Broadcast(e)

		}

	}()
	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test05MultiMsgMultiplexer04WaitForAllWithTimer() {
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
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b2 := concurrency.NewBCaster(dh1,
		"bcaster2",
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b3 := concurrency.NewBCaster(dh2,
		"bcaster3",
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	d := 200 * time.Millisecond
	//Create Multiplexer
	mp := concurrency.NewMultiMsgMultiplexer(dh3,
		"MultiMsg",
		concurrency.MultiMsgMultiplexerWaitForAll(true),
		concurrency.MultiMsgMultiplexerSequence(1),
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
			item := v.(*concurrency.Event)
			fmt.Printf("Multiplexer-key:%v\n", item.OutMessage.ID)

			m := item.OutMessage.Message.(*concurrency.SortedMap)
			//fmt.Printf("m:%v\n", m)
			for it := range m.Iter() {
				mapit := it.Value.(*concurrency.SortedMap)
				for mit := range mapit.Iter() {
					e := mit.Value.(*concurrency.Event)
					fmt.Printf("bcaster-key:%v - msgkey:%v - e.InitMessage:%v - e.OutMessage:%v\n", it.Key, mit.Key, e.InitMessage.Message, e.OutMessage.Message)
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

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b1", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b1.MsgType}

			b1.Broadcast(e)
			time.Sleep(500 * time.Microsecond)
		}

	}()
	go func() {
		for msgId := 0; msgId < 300; msgId++ {

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b2", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b2.MsgType}

			b2.Broadcast(e)
			time.Sleep(1 * time.Millisecond)
		}

	}()
	go func() {
		for msgId := 0; msgId < 300; msgId++ {

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b3", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b3.MsgType}

			b3.Broadcast(e)
			time.Sleep(500 * time.Microsecond)
		}

	}()
	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test05MultiMsgMultiplexer05WaitForAllWithTimerBM() {
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
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b2 := concurrency.NewBCaster(dh1,
		"bcaster2",
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)
	b3 := concurrency.NewBCaster(dh2,
		"bcaster3",
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	d := 1 * time.Microsecond
	//Create Multiplexer
	mp := concurrency.NewMultiMsgMultiplexer(dh3,
		"MultiMsg",
		concurrency.MultiMsgMultiplexerWaitForAll(true),
		concurrency.MultiMsgMultiplexerSequence(1),
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

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b1", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b1.MsgType}

			b1.Broadcast(e)
		}

	}()
	go func() {
		for msgId := 0; msgId < 300000; msgId++ {

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b2", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b2.MsgType}

			b2.Broadcast(e)
		}

	}()
	go func() {
		for msgId := 0; msgId < 300000; msgId++ {

			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: fmt.Sprintf("%v From b3", msgId), TimeInNano: time.Now().UnixNano(), MsgType: b3.MsgType}

			b3.Broadcast(e)
		}

	}()
	time.Sleep(1000 * time.Millisecond)

}
