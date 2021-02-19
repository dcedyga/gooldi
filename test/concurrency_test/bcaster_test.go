package concurrency_test

import (
	"fmt"
	concurrency "github.com/dcedyga/goold/concurrency"

	"time"

	uuid "github.com/satori/go.uuid"
)

// Broadcaster
func clientFunc(id int, b *concurrency.BCaster, dh *concurrency.DoneHandler) {

	msgCh := b.AddListener(dh)
	for msg := range msgCh {
		fmt.Printf("Client %d got message: %v\n", id, msg.(*concurrency.Event).OutMessage.Message)
	}
}
func clientFunc2(id int, b *concurrency.BCaster, dh *concurrency.DoneHandler) {
	wClose := make(chan interface{})
	defer concurrency.CloseChannel(wClose)
	msgCh := b.AddListener(dh)
	for msg := range msgCh {
		i := msg.(*concurrency.Event).OutMessage.Message.(int) * 2
		fmt.Printf("This is a different client that multiplies by 2 -> ClientId %v,  message: %v, value: %v\n", id, msg.(*concurrency.Event).OutMessage.Message, i)
	}
}
func clientFunc3(id int, b *concurrency.BCaster, dh *concurrency.DoneHandler) {
	wClose := make(chan interface{})
	defer concurrency.CloseChannel(wClose)
	msgCh := b.AddListener(dh)
	for msg := range msgCh {
		i := msg.(*concurrency.Event).OutMessage.Message.(int) * msg.(*concurrency.Event).OutMessage.Message.(int)
		fmt.Printf("This is the final client that squares the message -> ClientId %v,  message: %v, value: %v\n", id, msg.(*concurrency.Event).OutMessage.Message, i)
	}
}
func (suite *Suite) Test03BCaster01ThreeClients() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)

	b := concurrency.NewBCaster(dh,
		"bcaster",
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	go clientFunc(1, b, dh1)
	go clientFunc2(2, b, dh1)
	go clientFunc3(3, b, dh1)
	time.Sleep(1 * time.Millisecond)
	// Start publishing messages:
	go func() {
		for msgId := 0; ; msgId++ {
			e := &concurrency.Message{ID: uuid.NewV4().String(), Message: msgId, TimeInNano: time.Now().UnixNano(), MsgType: b.MsgType}

			b.Broadcast(e)
			if msgId == 8000 {
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("We are done from layer 0 bcaster -> %v: %v\n", dh.ID(), dh.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dh1.Done():
				fmt.Printf("We are done from layer 1 clients -> %v: %v\n", dh1.ID(), dh1.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dm.Done():
				fmt.Printf("We are done from DoneManager -> %v: %v\n", dm.ID(), dm.Err())
				return
			}
		}
	}()
	time.Sleep(1 * time.Millisecond)
	dm.GetDoneFunc()()
}

func (suite *Suite) Test03BCaster02SegregatedProcessors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	defer dm.GetDoneFunc()()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)

	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
		concurrency.BCasterTransformFn(concurrency.BCasterEventTransformFn),
	)

	//Create Processors
	w2 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 2),
		dh1,
		concurrency.ProcessorWithSequence(2),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
	)
	go w2.Process(process2)

	w3 := concurrency.NewProcessor(
		fmt.Sprintf("worker%v", 3),
		dh1,
		concurrency.ProcessorWithSequence(3),
		concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
		concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
	)
	go w3.Process(process3, 3)

	time.Sleep(300 * time.Microsecond)
	// Get all multiplexed messages
	go func() {
		i := 0
		for v := range w2.OutputChannel() {
			_ = v

			//fmt.Printf("Processor - %v, ProcessedMsg: %v\n", w2.ID, v.(concurrency.Event).ProcessedEvent)

			i++
		}
		fmt.Printf("Total Messages from Processor %v: %v,\n", w2.ID(), i)
	}()

	go func() {
		i := 0
		for v := range w3.OutputChannel() {
			_ = v

			//fmt.Printf("Processor - %v, ProcessedMsg: %v\n", w3.ID, v.(concurrency.Event).ProcessedEvent)

			i++
			//}

		}
		fmt.Printf("Total Messages from Processor %v: %v,\n", w3.ID(), i)
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
