package concurrency_test

import (
	"fmt"
	concurrency "github.com/dcedyga/gooldi/concurrency"
	"runtime"

	"time"
)

// func (suite *Suite) Test02Patterns01OrGoRoutine() {
// 	a, b, c, d, e := sendAfter(2*time.Second), sendAfter(5*time.Second), sendAfter(1*time.Second), sendAfter(3*time.Second), sendAfter(4*time.Second)
// 	for i := 0; i < 5; i++ {
// 		val := <-concurrency.Or(a, b, c, d, e)
// 		fmt.Printf("done after %v\n", val)
// 	}

// }

// //
// func sendAfter(after time.Duration) <-chan interface{} {
// 	c := make(chan interface{})
// 	go func() {
// 		//defer close(c)
// 		time.Sleep(after)
// 		fmt.Printf("After:%v\n", after)
// 		c <- fmt.Sprintf("Send After: %v\n", after)
// 	}()
// 	return c
// }

type Human struct {
	Name string
	Age  int
}

//Repeat
func (suite *Suite) Test02Patterns02RepeatGoRoutine() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	i := 0
	for item := range concurrency.Repeat(dh.Done(), &Human{Name: "Jon", Age: 23}) {
		fmt.Printf("%v\n", item)
		i++
		if i == 5 {
			break
		}
	}
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("%v\n", dh.Err())
				return
			}
		}
	}()
}

func (suite *Suite) Test02Patterns03TakeRepeatGoRoutine() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	for item := range concurrency.Take(dh.Done(), concurrency.Repeat(dh.Done(), &Human{Name: "Moñi", Age: 34}), 5) {
		fmt.Printf("%v\n", item)
	}
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("%v\n", dh.Err())
				return
			}
		}
	}()
}

func giveMeNow() interface{} {

	time.Sleep(100 * time.Microsecond)
	return time.Now()
}
func (suite *Suite) Test02Patterns04TakeRepeatFnGoRoutine() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()

	for item := range concurrency.Take(dh.Done(), concurrency.RepeatFn(dh.Done(), giveMeNow), 5) {
		fmt.Printf("%v\n", item)
	}
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("%v\n", dh.Err())
				return
			}
		}
	}()
}

func (suite *Suite) Test02Patterns05ToString() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	h := &Human{Name: "Jon", Age: 23}
	var message string
	for token := range concurrency.ToString(dh.Done(), concurrency.Take(dh.Done(), concurrency.Repeat(dh.Done(), h.Name, ":", h.Age), 3)) {
		message += token
	}
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("%v\n", dh.Err())
				return
			}
		}
	}()
	fmt.Printf("message: %s...\n", message)
}

func getTime(done <-chan interface{}) <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)

		select {
		case <-done:
			return
		case stream <- giveMeNow():
		}
	}()
	return stream
}

//FanIn: We can use it to gather state from all the subscriptions, timers, services, clients on a best effort bases
//and use this state as the entry point for the plan phase in a sense->plan->act pattern
func (suite *Suite) Test02Patterns06FanIn() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()

	start := time.Now()
	concurrentThreads := runtime.NumCPU()
	fmt.Printf("Spinning up %d concurrent Threads.\n", concurrentThreads)
	threads := make([]<-chan interface{}, concurrentThreads)

	for i := 0; i < concurrentThreads; i++ {
		threads[i] = getTime(dh.Done())
	}
	for t := range concurrency.Take(dh.Done(), concurrency.FanIn(dh.Done(), threads...), concurrentThreads) {
		fmt.Printf("\t%v\n", t)
		//_ = t
	}
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("%v\n", dh.Err())
				return
			}
		}
	}()
	fmt.Printf("took: %v\n", time.Since(start))
}
func (suite *Suite) Test02Patterns07FanInRec() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	start := time.Now()
	concurrentThreads := runtime.NumCPU()
	fmt.Printf("Spinning up %d concurrent Threads.\n", concurrentThreads)
	threads := make([]<-chan interface{}, concurrentThreads)

	for i := 0; i < concurrentThreads; i++ {
		threads[i] = getTime(dh.Done())
	}

	for t := range concurrency.Take(dh.Done(), concurrency.FanInRec(dh.Done(), threads...), concurrentThreads) {
		//_ = t
		fmt.Printf("\t%v\n", t)
	}
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("%v\n", dh.Err())
				return
			}
		}
	}()
	fmt.Printf("took: %v\n", time.Since(start))

}

// Split in n channels
func (suite *Suite) Test02Patterns08Route() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	out := concurrency.Route(dh.Done(), concurrency.Take(dh.Done(), concurrency.Repeat(dh.Done(), 1, 2, 3, 4, 5), 10), 5)
	for val1 := range out[0] {
		fmt.Printf("out1: %v, out2: %v, out3: %v, out4: %v, out5: %v\n", val1, <-out[1], <-out[2], <-out[3], <-out[4])
	}
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("%v\n", dh.Err())
				return
			}
		}
	}()
}

func (suite *Suite) Test02Patterns09Route() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	out := concurrency.Route(dh.Done(), concurrency.Take(dh.Done(), concurrency.Repeat(dh.Done(), 1, 2, 3, 4, 5), 5), 5)
	for index, item := range out {
		go func(index int, item chan interface{}) {
			for val := range item {
				fmt.Printf("Split - out[%v]: %v\n", index, val)
			}
		}(index, item)
	}
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("%v\n", dh.Err())
				return
			}
		}
	}()
	time.Sleep(time.Second)
}

func (suite *Suite) Test02Patterns10Bridge() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i*i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}

	for v := range concurrency.Bridge(dh.Done(), genVals()) {
		fmt.Printf("%v ", v)
	}
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("%v\n", dh.Err())
				return
			}
		}
	}()
}
