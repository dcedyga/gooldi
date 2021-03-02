package concurrency_test

import (
	"fmt"

	concurrency "github.com/dcedyga/gooldi/concurrency"

	"time"

	uuid "github.com/satori/go.uuid"
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
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, concurrency.MsgMultiplexerSequence(1))

	mp.Start()
	//Create Processors

	w2 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 2),
		dh1,
		concurrency.ProcessorWithSequence(2),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
	)
	mp.Set(w2.Sequence(), w2.OutputChannel())
	go w2.Process(process2)

	w3 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 3),
		dh1,
		concurrency.ProcessorWithSequence(3),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
	)
	mp.Set(w3.Sequence(), w3.OutputChannel())
	go w3.Process(process3, 3)

	w1 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 1),
		dh1,
		concurrency.ProcessorWithSequence(1),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
	)
	mp.Set(w1.Sequence(), w1.OutputChannel())
	go w1.Process(process1, 1)

	w4 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 4),
		dh1,
		concurrency.ProcessorWithSequence(4),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
	)
	mp.Set(w4.Sequence(), w4.OutputChannel())
	go w4.Process(process4, 4)

	filter := concurrency.NewFilter(
		fmt.Sprintf("filter%v", 5),
		dh1,
		concurrency.FilterWithSequence(5),
		concurrency.FilterWithInputChannel(mp.Iter()),
		concurrency.FilterTransformFn(concurrency.FilterEventTransformFn),
	)
	go filter.Filter(filterFn)

	time.Sleep(300 * time.Microsecond)
	// Get all multiplexed messages

	go func() {
		i := 0
		for v := range filter.OutputChannel() {
			_ = v
			item := v.(*concurrency.Event)
			s := item.OutMessage.Message.(*concurrency.Slice)
			is := item.InMessageSequence.GetItemAtIndex(0).(*concurrency.Slice)
			// fmt.Printf("Msg: %v - ProcessedEvents:%v \n", item.InitEvent.Message, item.ProcessedEvents.Len())
			for pe := range s.Iter() {
				index := pe.Index
				im := is.GetItemAtIndex(index).(concurrency.Slice)
				pr := pe.Value.(*concurrency.Message)
				var outmsg interface{}
				if pr.Message != nil {
					outmsg = pr.Message
				}

				eType := pr.MsgType
				//fmt.Printf("Sequence:%v - ProcessedEvent: %v - ProcessedMsg:%v \n", pr.Sequence, pr.ProcessedEvent, pr.InitEvent)
				fmt.Printf("eType:%v - inEventSeq(%v).Message:%v - msg:%v\n", eType, index, im.GetItemAtIndex(0).(*concurrency.Message).Message, outmsg)
			}
			i++

		}
		fmt.Printf("Total Messages Extracted: %v,\n", i)
	}()

	// Broadcast Messages

	//wd, _ := mp.Get(1)
	go func() {
		for msgId := 0; msgId < 30; msgId++ {
			if msgId == 10 {
				w2.Stop()
			} else if msgId == 15 {
				w2.Start()
			}
			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: msgId, TimeInNano: time.Now().UnixNano(), MsgType: b.MsgType}

			b.Broadcast(e)

		}

	}()
	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test04MsgMultiplexer02ProcessorsBM() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	defer dm.GetDoneFunc()()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2, concurrency.MsgMultiplexerSequence(1))

	mp.Start()
	//Create Processors

	w2 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 2),
		dh1,
		concurrency.ProcessorWithSequence(2),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
	)
	mp.Set(w2.Sequence(), w2.OutputChannel())
	go w2.Process(process2)

	w3 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 3),
		dh1,
		concurrency.ProcessorWithSequence(3),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
	)
	mp.Set(w3.Sequence(), w3.OutputChannel())
	go w3.Process(process3, 3)

	w1 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 1),
		dh1,
		concurrency.ProcessorWithSequence(1),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
	)
	mp.Set(w1.Sequence(), w1.OutputChannel())
	go w1.Process(process1, 1)

	w4 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 4),
		dh1,
		concurrency.ProcessorWithSequence(4),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
	)
	mp.Set(w4.Sequence(), w4.OutputChannel())
	go w4.Process(process4, 4)

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
			if msgId == 10 {
				w2.Stop()
			} else if msgId == 15 {
				w2.Start()
			}
			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: msgId, TimeInNano: time.Now().UnixNano(), MsgType: b.MsgType}

			b.Broadcast(e)

		}

	}()
	time.Sleep(1000 * time.Millisecond)

}

var globalI int = 0

func customGetLastRegisteredKey() int64 {
	return time.Now().UnixNano()
}

func customGetKeyFn(v interface{}) int64 {
	return v.(*CustomMsg).Key
}

func customGetSequenceFn(v interface{}) interface{} {
	return v.(*CustomMsg).Sequence
}

func customTransformFn(mp *concurrency.MsgMultiplexer, sm *concurrency.SortedMap) interface{} {
	globalI++
	fmt.Printf("globalI:%v\n", globalI)
	resultSlice := concurrency.NewSlice()
	var msgType string
	for item := range sm.Iter() {
		e := item.Value.(*CustomMsg)
		if msgType == "" {
			msgType = e.MsgType
		}
		resultSlice.Append(e.Msg)
	}

	re := &CustomMsg{
		Key:      time.Now().UnixNano(),
		MsgType:  msgType,
		Msg:      resultSlice,
		Sequence: mp.Sequence(),
	}
	return re
}

type CustomMsg struct {
	Key      int64
	Msg      interface{}
	MsgType  string
	Sequence interface{}
}

func customProcess1(input interface{}, params ...interface{}) interface{} {
	msg := input.(*CustomMsg)
	id := params[0]
	e := &CustomMsg{
		Key:      msg.Key,
		Msg:      fmt.Sprintf("From CustomMsg: Client %d got message: %v", id, msg.Msg),
		MsgType:  msg.MsgType,
		Sequence: id,
	}
	return e
}

func customProcess2(input interface{}, params ...interface{}) interface{} {
	msg := input.(*CustomMsg)
	id := params[0]
	i := msg.Msg.(int) * 2
	e := &CustomMsg{
		Key:      msg.Key,
		Msg:      fmt.Sprintf("From CustomMsg: This is a different client that multiplies by 2 -> ClientId %v,  message: %v, value: %v", id, msg.Msg, i),
		MsgType:  msg.MsgType,
		Sequence: id,
	}
	return e

}

func customProcessTransformFn(p *concurrency.Processor, input interface{}, result interface{}) interface{} {
	msg := input.(*CustomMsg)

	if result != nil {
		return result
	}

	e := &CustomMsg{
		Key:      msg.Key,
		Msg:      nil,
		MsgType:  msg.MsgType,
		Sequence: p.Sequence(),
	}

	return e
}

func (suite *Suite) Test04MsgMultiplexer03CustomMsg() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	defer dm.GetDoneFunc()()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	dh2 := dm.AddNewDoneHandler(2)

	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)

	//Create Multiplexer
	mp := concurrency.NewMsgMultiplexer(dh2,
		concurrency.MsgMultiplexerSequence(1),
		concurrency.MsgMultiplexerGetChannelItemKeyFn(customGetKeyFn),
		concurrency.MsgMultiplexerGetChannelItemSequenceFn(customGetSequenceFn),
		concurrency.MsgMultiplexerGetLastRegKeyFn(customGetLastRegisteredKey),
		concurrency.MsgMultiplexerTransformFn(customTransformFn),
	)

	mp.Start()

	w2 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 2),
		dh1,
		concurrency.ProcessorWithSequence(2),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(customProcessTransformFn),
	)
	mp.Set(w2.Sequence(), w2.OutputChannel())
	go w2.Process(customProcess2, 2)

	w1 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 1),
		dh1,
		concurrency.ProcessorWithSequence(1),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(customProcessTransformFn),
	)
	mp.Set(w1.Sequence(), w1.OutputChannel())
	go w1.Process(customProcess1, 1)

	go func() {
		i := 0
		for v := range mp.Iter() {
			item := v.(*CustomMsg)
			s := item.Msg.(*concurrency.Slice)
			for pe := range s.Iter() {
				fmt.Printf("key(%v), MsgType(%v), Sequence(%v) - CustomMsg.Msg: %v \n", item.Key, item.MsgType, item.Sequence, pe.Value)
			}
			i++

		}
		fmt.Printf("Total Messages Extracted: %v,\n", i)
	}()

	go func() {
		for msgId := 0; msgId < 30; msgId++ {
			if msgId == 10 {
				w2.Stop()
			} else if msgId == 15 {
				w2.Start()
			}
			e := &CustomMsg{Key: time.Now().UnixNano(), Msg: msgId, MsgType: b.MsgType, Sequence: 0}

			b.Broadcast(e)

		}

	}()
	time.Sleep(1000 * time.Millisecond)
}
