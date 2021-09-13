package concurrency_test

import (
	"fmt"

	concurrency "github.com/dcedyga/gooldi/concurrency"

	"time"
)

func customGetKeyFn(v interface{}) int64 {
	return v.(*CustomMsg).Key
}

func customTransformFn(mp *concurrency.MsgMultiplexer, sm *concurrency.SortedMap, correlationKey int64) interface{} {

	resultSlice := concurrency.NewSlice()
	var msgType string
	for item := range sm.Iter() {
		e := item.Value.(*CustomMsg)
		if msgType == "" {
			msgType = e.MsgType
		}
		resultSlice.Append(e)
	}

	re := &CustomMsg{
		Key:     time.Now().UnixNano(),
		MsgType: mp.MsgType,
		Msg:     resultSlice,
		Index:   mp.Index(),
	}
	return re
}

type CustomMsg struct {
	Key     int64
	Msg     interface{}
	MsgType string
	Index   int64
}

func customProcess1(pr *concurrency.Processor, input interface{}, params ...interface{}) interface{} {
	msg := input.(*CustomMsg)
	id := params[0]
	e := &CustomMsg{
		Key:     msg.Key,
		Msg:     fmt.Sprintf("From CustomMsg: Client %d got message: %v", id, msg.Msg),
		MsgType: msg.MsgType,
		Index:   int64(id.(int)),
	}
	return e
}

func customProcess2(pr *concurrency.Processor, input interface{}, params ...interface{}) interface{} {
	msg := input.(*CustomMsg)
	id := params[0]
	i := msg.Msg.(int) * 2
	e := &CustomMsg{
		Key:     msg.Key,
		Msg:     fmt.Sprintf("From CustomMsg: This is a different client that multiplies by 2 -> ClientId %v,  message: %v, value: %v", id, msg.Msg, i),
		MsgType: msg.MsgType,
		Index:   int64(id.(int)),
	}
	return e

}

func customProcessTransformFn(p *concurrency.Processor, input interface{}, result interface{}) interface{} {
	msg := input.(*CustomMsg)

	if result != nil {
		return result
	}

	e := &CustomMsg{
		Key:     msg.Key,
		Msg:     nil,
		MsgType: msg.MsgType,
		Index:   p.Index(),
	}

	return e
}

func (suite *Suite) Test06MsgMultiplexer01CustomMsg() {
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
	mp := concurrency.NewMsgMultiplexer(dh2, b.MsgType,
		concurrency.MsgMultiplexerIndex(1),
		concurrency.MsgMultiplexerGetItemKeyFn(customGetKeyFn),
		concurrency.MsgMultiplexerTransformFn(customTransformFn),
	)

	w2 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 2),
		dh1,
		concurrency.ProcessorWithIndex(2),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(customProcessTransformFn),
	)
	mp.Set(w2.Index(), w2.OutputChannel())
	go w2.Process(customProcess2, 2)

	w1 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 1),
		dh1,
		concurrency.ProcessorWithIndex(1),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(customProcessTransformFn),
	)
	mp.Set(w1.Index(), w1.OutputChannel())
	go w1.Process(customProcess1, 1)

	go func() {
		i := 0
		for v := range mp.Iter() {
			item := v.(*CustomMsg)
			s := item.Msg.(*concurrency.Slice)
			fmt.Printf("-------\n")
			for pev := range s.Iter() {
				pe := pev.Value.(*CustomMsg)
				fmt.Printf("key(%v), MsgType(%v), Index(%v) - CustomMsg.Msg: %v \n", pe.Key, pe.MsgType, pe.Index, pe.Msg)
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
			e := &CustomMsg{
				Key:     time.Now().UnixNano(),
				Msg:     msgId,
				MsgType: b.MsgType,
				Index:   0,
			}

			b.Broadcast(e)

		}

	}()
	time.Sleep(1000 * time.Millisecond)
}

//concurrency.CorrelatedMessage
