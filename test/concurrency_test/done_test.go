package concurrency_test

import (
	"fmt"
	concurrency "goold/concurrency"
	"time"
)

func (suite *Suite) Test01Done01DoneHandler() {
	dh := concurrency.NewDoneHandler()
	defer dh.GetDoneFunc()()
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("We are done: %v\n", dh.Err())
				return
			}
		}
	}()

	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test01Done02DoneHandlerWithTimeout() {

	dh := concurrency.NewDoneHandler(concurrency.DoneHandlerWithTimeout(100 * time.Millisecond))

	for {
		select {
		case <-dh.Done():
			fmt.Printf("We are done: %v\n", dh.Err())
			return
		}
	}

}

func (suite *Suite) Test01Done03DoneManager() {

	dm := concurrency.NewDoneManager()
	dh := concurrency.NewDoneHandler()
	dm.AddDoneHandler(dh, 0)

	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("We are done: %v\n", dh.Err())
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
	dm.GetDoneFunc()()
	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test01Done04DoneManagerMultiLayerNoDelay() {

	dm := concurrency.NewDoneManager()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)

	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("We are done on layer 0 -> %v: %v\n", dh.ID(), dh.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dh1.Done():
				fmt.Printf("We are done on layer 1 -> %v: %v\n", dh1.ID(), dh1.Err())
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
	dm.GetDoneFunc()()
	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test01Done05DoneManagerSingleLayer() {

	dm := concurrency.NewDoneManager()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(0)
	dh2 := dm.AddNewDoneHandler(0)

	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("We are done from layer 0 item1 -> %v: %v\n", dh.ID(), dh.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dh1.Done():
				fmt.Printf("We are done from layer 0 item2 -> %v: %v\n", dh1.ID(), dh1.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dh2.Done():
				fmt.Printf("We are done from layer 0 item3 -> %v: %v\n", dh2.ID(), dh2.Err())
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
	dm.GetDoneFunc()()
	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test01Done06DoneManagerWithdelay() {

	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)

	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("We are done with delay - on layer 0 -> %v: %v\n", dh.ID(), dh.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dh1.Done():
				fmt.Printf("We are done with delay - on layer 1 -> %v: %v\n", dh1.ID(), dh1.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dm.Done():
				fmt.Printf("We are done with delay - on DoneManager -> %v: %v\n", dm.ID(), dm.Err())
				return
			}
		}
	}()
	dm.GetDoneFunc()()
	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test01Done07DoneManagerWithdelayAndTimeout() {

	dm := concurrency.NewDoneManager(
		concurrency.DoneManagerWithDelay(1*time.Millisecond),
		concurrency.DoneManagerWithTimeout(100*time.Millisecond),
	)
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("We are done with delay and timeout - on layer 0 -> %v: %v\n", dh.ID(), dh.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dh1.Done():
				fmt.Printf("We are done with delay and timeout - on layer 1 -> %v: %v\n", dh1.ID(), dh1.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dm.Done():
				fmt.Printf("We are done with delay and timeout - on DoneManager -> %v: %v\n", dm.ID(), dm.Err())
				return
			}
		}
	}()

	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test01Done08DoneManagerWithdelayTimeoutAndDoneHandlerTimeout() {

	dm := concurrency.NewDoneManager(
		concurrency.DoneManagerWithDelay(1*time.Millisecond),
		concurrency.DoneManagerWithTimeout(100*time.Millisecond),
	)
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1, concurrency.DoneHandlerWithTimeout(10*time.Millisecond))
	dh2 := dm.AddNewDoneHandler(2)
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("We are done with delay and timeout - on layer 0 -> %v: %v\n", dh.ID(), dh.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dh1.Done():
				fmt.Printf("We are done with delay and timeout - on layer 1 -> %v: %v\n", dh1.ID(), dh1.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dh2.Done():
				fmt.Printf("We are done with delay and timeout - on layer 2 -> %v: %v\n", dh2.ID(), dh2.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dm.Done():
				fmt.Printf("We are done with delay and timeout - on DoneManager -> %v: %v\n", dm.ID(), dm.Err())
				return
			}
		}
	}()

	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test01Done09DoneManagerGetDoneHandler() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)

	gdh, layer, founded := dm.GetDoneHandler(concurrency.QueryDoneHandlerWithKey(dh.ID()))
	fmt.Printf("Found DoneHandler: %v,%v,%v\n", gdh, layer, founded)
	go func() {
		for {
			select {
			case <-dh.Done():
				fmt.Printf("We are done with delay - on layer 0 -> %v: %v\n", dh.ID(), dh.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dh1.Done():
				fmt.Printf("We are done with delay - on layer 1 -> %v: %v\n", dh1.ID(), dh1.Err())
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-dm.Done():
				fmt.Printf("We are done with delay - on DoneManager -> %v: %v\n", dm.ID(), dm.Err())
				return
			}
		}
	}()
	dm.GetDoneFunc()()
	time.Sleep(1000 * time.Millisecond)
}
