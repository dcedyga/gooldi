package concurrency_test

import (
	"fmt"
	"runtime"

	concurrency "github.com/dcedyga/gooldi/concurrency"

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
				stream <- i * i
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

func (suite *Suite) Test02Patterns11AsChan() {
	primes := [6]int{2, 3, 5, 7, 11, 13}

	for i := range concurrency.AsChan(primes) {
		fmt.Printf("prime:%v\n", i)
	}

}
func (suite *Suite) Test02Patterns12Or() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-concurrency.Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Millisecond),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}

func (suite *Suite) Test02Patterns12OrJust2() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-concurrency.Or(
		sig(5*time.Minute),
		sig(1*time.Millisecond),
	)
	fmt.Printf("done after %v", time.Since(start))
}
func giveMeNowParam(param interface{}) interface{} {

	fmt.Printf("param:%v\n", param)
	time.Sleep(100 * time.Microsecond)
	return time.Now()
}
func giveMeNowParams(params ...interface{}) interface{} {

	fmt.Printf("params[0]:%v,params[1]:%v\n", params[0], params[1])
	time.Sleep(100 * time.Microsecond)
	return time.Now()
}
func (suite *Suite) Test02Patterns13TakeRepeatFnParamGoRoutine() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	for j := 0; j < 5; j++ {
		for item := range concurrency.Take(dh.Done(), concurrency.RepeatParamFn(dh.Done(), giveMeNowParam, j), 5) {
			fmt.Printf("%v\n", item)
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

func (suite *Suite) Test02Patterns14TakeRepeatFnParamsGoRoutine() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	for j := 0; j < 5; j++ {
		for item := range concurrency.Take(dh.Done(), concurrency.RepeatParamsFn(dh.Done(), giveMeNowParams, j, j*j), 5) {
			fmt.Printf("%v\n", item)
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
func doneFn() {
	fmt.Printf("Before ending: from doneFn\n")

}
func doneParamsFn(params ...interface{}) {
	for _, p := range params {
		fmt.Printf("Before ending - from doneParamsFn, param: %v \n", p)
	}
}
func doneChanParamFn(param chan interface{}) {
	for p := range param {
		fmt.Printf("Before ending - from doneParamsFn, param: %v \n", p)
	}
}
func doneChanParamsFn(params ...chan interface{}) {
	for _, p := range params {
		go func(c chan interface{}) {
			for i := range c {
				fmt.Printf("Before ending - from doneParamsFn, param: %v \n", i)
			}
		}(p)
	}
}
func (suite *Suite) Test02Patterns15OrDoneFn() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	c := concurrency.OrDoneFn(dh.Done(), make(chan interface{}), doneFn)
	fmt.Printf("Channel:%v\n", c)
}
func (suite *Suite) Test02Patterns16OrDoneParamsFn() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	c := concurrency.OrDoneParamsFn(dh.Done(), make(chan interface{}), doneParamsFn, 1, 2)
	fmt.Printf("Channel:%v\n", c)
}
func (suite *Suite) Test02Patterns17OrDoneChanParamFn() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	c1 := make(chan interface{})
	c := concurrency.OrDoneChanParamFn(dh.Done(), make(chan interface{}), doneChanParamFn, c1)
	go func() {
		for i := 0; i < 5; i++ {
			c1 <- i
		}
	}()
	fmt.Printf("Channel:%v\n", c)
}

func (suite *Suite) Test02Patterns17OrDoneChanParamsFn() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	c1 := make(chan interface{})
	c2 := make(chan interface{})
	c := concurrency.OrDoneChanParamsFn(dh.Done(), make(chan interface{}), doneChanParamsFn, c1, c2)
	go func() {
		for i := 0; i < 5; i++ {
			c1 <- i
		}
	}()
	go func() {
		for i := 5; i < 10; i++ {
			c2 <- i
		}
	}()
	fmt.Printf("Channel:%v\n", c)
}
