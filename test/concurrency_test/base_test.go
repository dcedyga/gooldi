package concurrency_test

import (
	"fmt"
	"testing"

	concurrency "github.com/dcedyga/gooldi/concurrency"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

//
func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

//
func (suite *Suite) SetupSuite() {}

//
func (suite *Suite) SetupTest() {}

//
func (suite *Suite) TearDownSuite() {

}

// Broadcaster
func clientFunc(id int, b *concurrency.BCaster, dh *concurrency.DoneHandler) {
	msgCh := b.AddListener(dh)
	for msg := range msgCh {
		fmt.Printf("Client %d got message: %v\n", id, msg.(*concurrency.Message).Message)
	}
}
func clientFunc2(id int, b *concurrency.BCaster, dh *concurrency.DoneHandler) {
	// wClose := make(chan interface{})
	// defer concurrency.CloseChannel(wClose)
	msgCh := b.AddListener(dh)
	for msg := range msgCh {
		i := msg.(*concurrency.Message).Message.(int) * 2
		fmt.Printf("This is a different client that multiplies by 2 -> ClientId %v,  message: %v, value: %v\n", id, msg.(*concurrency.Message).Message, i)
	}
}
func clientFunc3(id int, b *concurrency.BCaster, dh *concurrency.DoneHandler) {
	// wClose := make(chan interface{})
	// defer concurrency.CloseChannel(wClose)
	msgCh := b.AddListener(dh)
	for msg := range msgCh {
		i := msg.(*concurrency.Message).Message.(int) * msg.(*concurrency.Message).Message.(int)
		fmt.Printf("This is the final client that squares the message -> ClientId %v,  message: %v, value: %v\n", id, msg.(*concurrency.Message).Message, i)
	}
}

func process1(pr *concurrency.Processor, input interface{}, params ...interface{}) interface{} {
	msg := input.(*concurrency.Message)
	id := params[0]
	return concurrency.NewMessage(fmt.Sprintf("ClientId %v: echo -> message: %v, value: %v", id, msg.Message, msg.Message),
		msg.MsgType,
		concurrency.MessageWithCorrelationKey(msg.CorrelationKey),
		concurrency.MessageWithIndex(pr.Index()),
	)
}

// func process2(pr *concurrency.Processor, input interface{}, params ...interface{}) interface{} {
// 	msg := input.(*concurrency.Message)
// 	id := 2
// 	i := msg.Message.(int) * 2
// 	return concurrency.NewMessage(fmt.Sprintf("This is a different client that multiplies by 2 -> ClientId %v,  message: %v, value: %v", id, msg.Message, i),
// 		msg.MsgType,
// 		concurrency.MessageWithCorrelationKey(msg.CorrelationKey),
// 	)

// }
func process2(pr *concurrency.Processor, input interface{}, params ...interface{}) interface{} {
	msg := input.(*concurrency.Message)
	id := params[0]
	i := msg.Message.(int) * 2
	return concurrency.NewMessage(fmt.Sprintf("ClientId %v: multiplies by 2 -> message: %v, value: %v", id, msg.Message, i),
		msg.MsgType,
		concurrency.MessageWithCorrelationKey(msg.CorrelationKey),
		concurrency.MessageWithIndex(pr.Index()),
	)

}
func process3(pr *concurrency.Processor, input interface{}, params ...interface{}) interface{} {
	msg := input.(*concurrency.Message)
	id := params[0]
	i := msg.Message.(int) * msg.Message.(int)
	return concurrency.NewMessage(fmt.Sprintf("ClientId %v: squares the message -> message: %v, value: %v", id, msg.Message, i),
		msg.MsgType,
		concurrency.MessageWithCorrelationKey(msg.CorrelationKey),
		concurrency.MessageWithIndex(pr.Index()),
	)
}
func process4(pr *concurrency.Processor, input interface{}, params ...interface{}) interface{} {
	msg := input.(*concurrency.Message)
	id := params[0]
	i := msg.Message.(int) * 3
	return concurrency.NewMessage(fmt.Sprintf("ClientId %v: multiplies by 3 -> message: %v, value: %v", id, msg.Message, i),
		msg.MsgType,
		concurrency.MessageWithCorrelationKey(msg.CorrelationKey),
		concurrency.MessageWithIndex(pr.Index()),
	)
}

func process5(pr *concurrency.Processor, input interface{}, params ...interface{}) interface{} {
	msg := input.(*concurrency.Message)
	id := params[0]
	i := msg.Message.(int) * msg.Message.(int) * msg.Message.(int)
	return concurrency.NewMessage(fmt.Sprintf("ClientId %v: to the power of 3 -> message: %v, value: %v", id, msg.Message, i),
		msg.MsgType,
		concurrency.MessageWithCorrelationKey(msg.CorrelationKey),
		concurrency.MessageWithIndex(pr.Index()),
	)
}

func process6(pr *concurrency.Processor, input interface{}, params ...interface{}) interface{} {
	msg := input.(*concurrency.Message)
	id := params[0]
	i := msg.Message.(int) * 4
	return concurrency.NewMessage(fmt.Sprintf("ClientId %v: multiplies by 4 -> message: %v, value: %v", id, msg.Message, i),
		msg.MsgType,
		concurrency.MessageWithCorrelationKey(msg.CorrelationKey),
		concurrency.MessageWithIndex(pr.Index()),
	)
}

func process7(pr *concurrency.Processor, input interface{}, params ...interface{}) interface{} {
	msg := input.(*concurrency.Message)
	id := params[0]
	i := msg.Message.(int) * msg.Message.(int) * msg.Message.(int) * msg.Message.(int)
	return concurrency.NewMessage(fmt.Sprintf("ClientId %v: to the power of 4 -> message: %v, value: %v", id, msg.Message, i),
		msg.MsgType,
		concurrency.MessageWithCorrelationKey(msg.CorrelationKey),
		concurrency.MessageWithIndex(pr.Index()),
	)
}

func process8(pr *concurrency.Processor, input interface{}, params ...interface{}) interface{} {
	msg := input.(*concurrency.Message)
	id := params[0]
	i := msg.Message.(int) * 5
	return concurrency.NewMessage(fmt.Sprintf("ClientId %v: multiplies by 5 -> message: %v, value: %v", id, msg.Message, i),
		msg.MsgType,
		concurrency.MessageWithCorrelationKey(msg.CorrelationKey),
		concurrency.MessageWithIndex(pr.Index()),
	)
}

func createMessagePairProcessors(mp *concurrency.MsgMultiplexer, b *concurrency.BCaster, dh1 *concurrency.DoneHandler, numWorkers int) *concurrency.SortedSlice {
	s := concurrency.NewSortedSlice()
	for i := 0; i < numWorkers; i++ {
		w := concurrency.NewProcessor(
			fmt.Sprintf("worker%v", i),
			dh1,
			concurrency.ProcessorWithIndex(int64(i)),
			concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
			concurrency.ProcessorTransformFn(concurrency.ProcessorMessagePairTransformFn),
		)
		if mp != nil {
			mp.Set(w.Index(), w.OutputChannel())
		}

		if i%8 == 0 {
			go w.Process(process8, i)
		} else if i%7 == 0 {
			go w.Process(process7, i)
		} else if i%6 == 0 {
			go w.Process(process6, i)
		} else if i%5 == 0 {
			go w.Process(process5, i)
		} else if i%4 == 0 {
			go w.Process(process4, i)
		} else if i%3 == 0 {
			go w.Process(process3, i)
		} else if i%2 == 0 {
			go w.Process(process2, i)
		} else {
			go w.Process(process1, i)
		}

		s.Append(w)
	}
	return s
}

func createProcessors(mp *concurrency.MsgMultiplexer, b *concurrency.BCaster, dh1 *concurrency.DoneHandler, numWorkers int) *concurrency.SortedSlice {
	s := concurrency.NewSortedSlice()
	for i := 0; i < numWorkers; i++ {
		w := concurrency.NewProcessor(
			fmt.Sprintf("worker%v", i),
			dh1,
			concurrency.ProcessorWithIndex(int64(i)),
			concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
			//concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
		)
		if mp != nil {
			mp.Set(w.Index(), w.OutputChannel())
		}

		if i%8 == 0 {
			go w.Process(process8, i)
		} else if i%7 == 0 {
			go w.Process(process7, i)
		} else if i%6 == 0 {
			go w.Process(process6, i)
		} else if i%5 == 0 {
			go w.Process(process5, i)
		} else if i%4 == 0 {
			go w.Process(process4, i)
		} else if i%3 == 0 {
			go w.Process(process3, i)
		} else if i%2 == 0 {
			go w.Process(process2, i)
		} else {
			go w.Process(process1, i)
		}
		s.Append(w)
	}
	return s
}

func getTotalMessagesFromProcessors(ss *concurrency.SortedSlice, c chan int) {
	j := 0
	go func(s *concurrency.SortedSlice) {
		for item := range s.Iter() {
			p := item.Value.(*concurrency.Processor)
			go func(w *concurrency.Processor) {
				i := 0
				for v := range w.OutputChannel() {
					_ = v
					//fmt.Printf("Processor - %v, ProcessedMsg: %v\n", w.ID(), v.(*concurrency.Message).Message)
					i++
				}
				c <- i
			}(p)
		}
	}(ss)

	go func() {
		for it := range c {
			j = j + it
		}
		fmt.Printf("Total Messages Extracted: %v,\n", j)
	}()

}

func broadcast( b *concurrency.BCaster,numMsg int){
	go func() {
		for msgId := 0; msgId < numMsg; msgId++ {
			e := concurrency.NewMessage(
				msgId,
				b.MsgType,
			)
			b.Broadcast(e)

		}

	}()
}
func broadcastWithStop( b *concurrency.BCaster,numMsg int,stopAt int, startAt int, w *concurrency.Processor){
	go func() {
		for msgId := 0; msgId < numMsg; msgId++ {
			if msgId == stopAt {
				w.Stop()
			} else if msgId == startAt {
				w.Start()
			}
			e := concurrency.NewMessage(
				msgId,
				b.MsgType,
			)

			b.Broadcast(e)

		}

	}()
}

func printMultiplexedResult(mp *concurrency.MsgMultiplexer){
	go func() {
		i := 0
		j := 0
		for v := range mp.Iter() {
			//mop := v.(*concurrency.SortedMap)
			msg := v.(*concurrency.Message)

			mop := msg.Message.(*concurrency.SortedMap)
			j++
			i = i + mop.Len()

		}
		fmt.Printf("Total Messages Extracted: %v,\n", j)
		fmt.Printf("Total Messages Multiplexed: %v,\n", i)
	}()
}

func filterFn(f *concurrency.Filter, input interface{}, params ...interface{}) bool {
	msg := input.(*concurrency.MessagePair)

	i := msg.In.Message.(int) % 2

	return i == 0

}
