package concurrency_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
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
func (suite *Suite) TearDownSuite() {}
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
func slowprocess9(pr *concurrency.Processor, input interface{}, params ...interface{}) interface{} {
	msg := input.(*concurrency.Message)
	id := params[0]
	i := msg.Message.(int) * 15
	time.Sleep(5*time.Millisecond)
	return concurrency.NewMessage(fmt.Sprintf("ClientId %v: multiplies by 15 but is slow -> message: %v, value: %v", id, msg.Message, i),
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
func createProcessorsFrom(mp *concurrency.MsgMultiplexer, b *concurrency.BCaster, dh1 *concurrency.DoneHandler, indexFrom, numWorkers int) *concurrency.SortedSlice {
	s := concurrency.NewSortedSlice()
	for i := indexFrom; i < (numWorkers + indexFrom); i++ {
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
func getTotalMessagesFromProcessors(ss *concurrency.SortedSlice, c chan interface{}) {

	go func(s *concurrency.SortedSlice, c chan interface{}) {
		var wg sync.WaitGroup
		for item := range s.Iter() {
			p := item.Value.(*concurrency.Processor)
			wg.Add(1)
			go func(w *concurrency.Processor, wg *sync.WaitGroup) {
				defer wg.Done()
				i := 0
				for v := range w.OutputChannel() {
					_ = v
					//fmt.Printf("Processor - %v, ProcessedMsg: %v\n", w.ID(), v.(*concurrency.Message).Message)
					i++
				}
				c <- i
			}(p, &wg)
		}
		wg.Wait()
		close(c)
	}(ss, c)

	go func(c chan interface{}) {
		j := 0
		for it := range c {
			j = j + it.(int)
		}
		fmt.Printf("Total Messages Extracted: %v,\n", j)
	}(c)

}
func broadcast(b *concurrency.BCaster, numMsg int) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					msgId,
					b.MsgType,
				)
				b.Broadcast(e)
				i++
			}
			if exit {
				break
			}
		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func broadcastWithStop(b *concurrency.BCaster, numMsg int, stopAt int, startAt int, w *concurrency.Processor) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:
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
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func broadcastWithDeleteAfterCKey(b *concurrency.BCaster, numMsg int, stopAt int, mp *concurrency.MsgMultiplexer, w *concurrency.Processor) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:

				e := concurrency.NewMessage(
					msgId,
					b.MsgType,
				)
				if msgId == stopAt {
					mp.DeleteAtCorrelationKey(w.Index(), e.CorrelationKey)
				}
				b.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func broadcastWithDeleteAfterCKey3(b *concurrency.BCaster, numMsg int, delAt1,delAt2,delAt3 int, mp *concurrency.MsgMultiplexer, w1,w2,w3 *concurrency.Processor) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:

				e := concurrency.NewMessage(
					msgId,
					b.MsgType,
				)
				if msgId == delAt1 {
					mp.DeleteAtCorrelationKey(w1.Index(), e.CorrelationKey)
				}
				if msgId == delAt2 {
					mp.DeleteAtCorrelationKey(w2.Index(), e.CorrelationKey)
				}
				if msgId == delAt3 {
					mp.DeleteAtCorrelationKey(w3.Index(), e.CorrelationKey)
				}
				b.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func broadcastWithDeleteAddAtCKeys(b *concurrency.BCaster, numMsg int,deleteAt int, addAt int,numWorker int, dh1 *concurrency.DoneHandler,mp *concurrency.MsgMultiplexer,w1 *concurrency.Processor) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:

				e := concurrency.NewMessage(
					msgId,
					b.MsgType,
				)
				if msgId == deleteAt {
					mp.DeleteAtCorrelationKey(w1.Index(), e.CorrelationKey)
				}
				if msgId == addAt {
					w := concurrency.NewProcessor(
						fmt.Sprintf("worker%v", numWorker),
						dh1,
						concurrency.ProcessorWithIndex(int64(numWorker)),
						concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
						//concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
					)
					if mp != nil {
						mp.SetAtCorrelationKey(w.Index(), w.OutputChannel(),e.CorrelationKey)
					}
					if numWorker%8 == 0 {
						go w.Process(process8, numWorker)
					} else if numWorker%7 == 0 {
						go w.Process(process7, numWorker)
					} else if numWorker%6 == 0 {
						go w.Process(process6, numWorker)
					} else if numWorker%5 == 0 {
						go w.Process(process5, numWorker)
					} else if numWorker%4 == 0 {
						go w.Process(process4, numWorker)
					} else if numWorker%3 == 0 {
						go w.Process(process3, numWorker)
					} else if numWorker%2 == 0 {
						go w.Process(process2, numWorker)
					} else {
						go w.Process(process1, numWorker)
					}
				}
				b.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func broadcastWithAddAtCKey(b *concurrency.BCaster, numMsg int, addAt int,numWorker int, dh1 *concurrency.DoneHandler,mp *concurrency.MsgMultiplexer) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:

				e := concurrency.NewMessage(
					msgId,
					b.MsgType,
				)
				if msgId == addAt {
					w := concurrency.NewProcessor(
						fmt.Sprintf("worker%v", numWorker),
						dh1,
						concurrency.ProcessorWithIndex(int64(numWorker)),
						concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
						//concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
					)
					if mp != nil {
						mp.SetAtCorrelationKey(w.Index(), w.OutputChannel(),e.CorrelationKey)
					}
					if numWorker%8 == 0 {
						go w.Process(process8, numWorker)
					} else if numWorker%7 == 0 {
						go w.Process(process7, numWorker)
					} else if numWorker%6 == 0 {
						go w.Process(process6, numWorker)
					} else if numWorker%5 == 0 {
						go w.Process(process5, numWorker)
					} else if numWorker%4 == 0 {
						go w.Process(process4, numWorker)
					} else if numWorker%3 == 0 {
						go w.Process(process3, numWorker)
					} else if numWorker%2 == 0 {
						go w.Process(process2, numWorker)
					} else {
						go w.Process(process1, numWorker)
					}
				}
				b.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func broadcastWitStartAfterCKeys(b *concurrency.BCaster, numMsg int, startAt int,mp *concurrency.MsgMultiplexer) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					msgId,
					b.MsgType,
				)
				if msgId == startAt {
					mp.StartAfterCorrelationKey(e.CorrelationKey)
				}
				b.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func broadcastWitStopAfterCKeys(b *concurrency.BCaster, numMsg int,stopAt int, mp *concurrency.MsgMultiplexer) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:

				e := concurrency.NewMessage(
					msgId,
					b.MsgType,
				)
				if msgId == stopAt {
					mp.StopAfterCorrelationKey(e.CorrelationKey)
				}
			
				b.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func broadcastWitStopStartAfterCKeys(b *concurrency.BCaster, numMsg int,stopAt int, startAt int,mp *concurrency.MsgMultiplexer) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:

				e := concurrency.NewMessage(
					msgId,
					b.MsgType,
				)
				if msgId == stopAt {
					mp.StopAfterCorrelationKey(e.CorrelationKey)
				}
				if msgId == startAt {
					mp.StartAfterCorrelationKey(e.CorrelationKey)
				}
				b.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func broadcastWitStopAddDeleteStartAfterCKeys(b *concurrency.BCaster, numMsg int,stopAt int, startAt int,numWorker int,dh1 *concurrency.DoneHandler,mp *concurrency.MsgMultiplexer,w *concurrency.Processor) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:

				e := concurrency.NewMessage(
					msgId,
					b.MsgType,
				)
				if msgId == stopAt {
					mp.StopAfterCorrelationKey(e.CorrelationKey)
				}
				if msgId == (stopAt+1) {
					mp.DeleteAtCorrelationKey(w.Index(),e.CorrelationKey)
					
				}
				if msgId == startAt {
					w1 := concurrency.NewProcessor(
						fmt.Sprintf("worker%v", numWorker),
						dh1,
						concurrency.ProcessorWithIndex(int64(numWorker)),
						concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
						//concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
					)
					if mp != nil {
						mp.Set(w1.Index(), w1.OutputChannel())
					}
					if numWorker%8 == 0 {
						go w1.Process(process8, numWorker)
					} else if numWorker%7 == 0 {
						go w1.Process(process7, numWorker)
					} else if numWorker%6 == 0 {
						go w1.Process(process6, numWorker)
					} else if numWorker%5 == 0 {
						go w1.Process(process5, numWorker)
					} else if numWorker%4 == 0 {
						go w1.Process(process4, numWorker)
					} else if numWorker%3 == 0 {
						go w1.Process(process3, numWorker)
					} else if numWorker%2 == 0 {
						go w1.Process(process2, numWorker)
					} else {
						go w1.Process(process1, numWorker)
					}
					mp.StartAfterCorrelationKey(e.CorrelationKey)
				}
				b.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func broadcastWitStopDeleteStartAfterCKeys(b *concurrency.BCaster, numMsg int,stopAt int, startAt int,mp *concurrency.MsgMultiplexer,w *concurrency.Processor) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:

				e := concurrency.NewMessage(
					msgId,
					b.MsgType,
				)
				if msgId == stopAt {
					mp.StopAfterCorrelationKey(e.CorrelationKey)
				}
				if msgId == (stopAt+1) {
					mp.DeleteAtCorrelationKey(w.Index(),e.CorrelationKey)
				}
				if msgId == startAt {
					mp.StartAfterCorrelationKey(e.CorrelationKey)
				}
				b.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func broadcastWitStopAddStartAfterCKeys(b *concurrency.BCaster, numMsg int,stopAt int, startAt int,numWorker int,dh1 *concurrency.DoneHandler,mp *concurrency.MsgMultiplexer) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:

				e := concurrency.NewMessage(
					msgId,
					b.MsgType,
				)
				if msgId == stopAt {
					mp.StopAfterCorrelationKey(e.CorrelationKey)
				}
				if msgId == startAt {
					w1 := concurrency.NewProcessor(
						fmt.Sprintf("worker%v", numWorker),
						dh1,
						concurrency.ProcessorWithIndex(int64(numWorker)),
						concurrency.ProcessorWithInputChannel(b.AddListener(dh1)),
						//concurrency.ProcessorTransformFn(concurrency.ProcessorEventTransformFn),
					)
					if mp != nil {
						mp.Set(w1.Index(), w1.OutputChannel())
					}
					if numWorker%8 == 0 {
						go w1.Process(process8, numWorker)
					} else if numWorker%7 == 0 {
						go w1.Process(process7, numWorker)
					} else if numWorker%6 == 0 {
						go w1.Process(process6, numWorker)
					} else if numWorker%5 == 0 {
						go w1.Process(process5, numWorker)
					} else if numWorker%4 == 0 {
						go w1.Process(process4, numWorker)
					} else if numWorker%3 == 0 {
						go w1.Process(process3, numWorker)
					} else if numWorker%2 == 0 {
						go w1.Process(process2, numWorker)
					} else {
						go w1.Process(process1, numWorker)
					}
					mp.StartAfterCorrelationKey(e.CorrelationKey)
				}
				b.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func broadcastWithAddDelete(b *concurrency.BCaster, numMsg int, changeAt int, mp *concurrency.MsgMultiplexer, w *concurrency.Processor, dh *concurrency.DoneHandler) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:

				e := concurrency.NewMessage(
					msgId,
					b.MsgType,
				)
				if msgId == changeAt {
					mp.DeleteAtCorrelationKey(w.Index(), e.CorrelationKey)
					createProcessorsFrom(mp, b, dh, 4, 1)
				}
				b.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func broadcastWithDeleteAfterCKeyDecoupled(b *concurrency.BCaster, numMsg int, stopAt int, mp *concurrency.MsgMultiplexer, w *concurrency.Processor) {
	go func() {
		i := 0
		exit := false
		dh := b.DoneHandler()
		for msgId := 0; msgId < numMsg; msgId++ {
			select {
			case <-dh.Done():
				exit = true
			default:
				e := concurrency.NewMessage(
					msgId,
					b.MsgType,
				)
				b.Broadcast(e)
				i++
			}
			if exit {
				break
			}

		}
		fmt.Printf("Total Messages Broadcasted: %v,\n", i)
	}()
}
func printMultiplexedResult(mp *concurrency.MsgMultiplexer) {
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
func done(dm *concurrency.DoneManager, mp *concurrency.MsgMultiplexer) {
	dm.GetDoneFunc()()
	//time.Sleep(600 * time.Millisecond)
	mp.PrintPreStreamMap()
	//mp.PrintDebug()
}
func doneWithBCaster(dm *concurrency.DoneManager, mp *concurrency.MsgMultiplexer, b *concurrency.BCaster) {
	dm.GetDoneFunc()()
	//time.Sleep(600 * time.Millisecond)
	mp.PrintDebug()
	//fmt.Printf("PrintPreStreamMap turn\n")
	mp.PrintPreStreamMap()

}
func printMultMsgMultiplexedResult(mp *concurrency.MultiMsgMultiplexer) {
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
func filterTransformFn(f *concurrency.Filter, input interface{}) interface{} {
	return input
}
func multiMsgGetItemKey(v interface{}) int64 {
	return v.(*concurrency.Message).CorrelationKey
}
// defaultMultiMsgTransformFn - Transforms the SortedMap output into a Message for future consumption as part of the
// output channel of the MultiMsgMultiplexer. Can be overridden for a more generic implementation
// Oriented
func multiMsgTransformFn(mp *concurrency.MultiMsgMultiplexer, r *concurrency.SortedMap) interface{} {
	m := concurrency.NewSortedMap()
	for item := range r.Iter() {
		it := item.Value.(*concurrency.SortedMap).Clone()
		m.Set(item.Key, it)
	}
	rmsg := concurrency.NewMessage(m,
		mp.MsgType,
		concurrency.MessageWithIndex(mp.Index()),
	)

	return rmsg
}
