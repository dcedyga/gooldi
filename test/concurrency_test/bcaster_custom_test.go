package concurrency_test

import (
	"fmt"

	concurrency "github.com/dcedyga/gooldi/concurrency"
	"github.com/stretchr/testify/assert"

	"time"
)

type CustomEvent struct {
	InitMessage *concurrency.Message
	OutMessage  *concurrency.Message
}

func customProcess3(pr *concurrency.Processor, input interface{}, params ...interface{}) interface{} {
	event := input.(*CustomEvent)
	//id := params[0]
	i := event.OutMessage.Message.(int) * event.OutMessage.Message.(int)
	return concurrency.NewMessage(i,
		event.OutMessage.MsgType,
		concurrency.MessageWithIndex(event.OutMessage.Index),
		concurrency.MessageWithCorrelationKey(event.OutMessage.CorrelationKey),
	)
}

func BCasterCustomEventTransformFn(b *concurrency.BCaster, input interface{}) interface{} {
	var msg *concurrency.Message
	if input != nil {
		msg = input.(*concurrency.Message)
	}
	e := &CustomEvent{
		InitMessage: msg,
		OutMessage:  msg,
	}
	return e
}

func ProcessorCustomEventTransformFn(p *concurrency.Processor, input interface{}, result interface{}) interface{} {
	var event *CustomEvent
	var resultMsg *concurrency.Message
	if input != nil {
		event = input.(*CustomEvent)
	}
	if result != nil {
		resultMsg = result.(*concurrency.Message)
	}

	r := &CustomEvent{
		InitMessage: event.InitMessage,
		OutMessage:  resultMsg,
	}
	return r
}

func (suite *Suite) Test04BCaster01CustomProcessors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	defer dm.GetDoneFunc()()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)

	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		concurrency.BCasterTransformFn(BCasterCustomEventTransformFn),
	)

	//Create Processors
	w2 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 2),
		dh1,
		concurrency.ProcessorWithIndex(2),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(ProcessorCustomEventTransformFn),
	)
	go w2.Process(customProcess3, 2)

	time.Sleep(300 * time.Microsecond)
	// Get all multiplexed messages
	go func() {
		i := 0
		for v := range w2.OutputChannel() {
			_ = v

			//fmt.Printf("Processor - %v, ProcessedMsg: %v\n", w2.ID(), v.(*CustomEvent).OutMessage)

			i++
		}
		assert.Greater(suite.T(), i, 0)
		fmt.Printf("Total Messages from Processor %v: %v,\n", w2.ID(), i)
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
			e := concurrency.NewMessage(
				msgId,
				b.MsgType,
			)

			b.Broadcast(e)

		}

	}()
	time.Sleep(1000 * time.Millisecond)

}
