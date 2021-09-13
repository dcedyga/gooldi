package concurrency_test

import (
	"fmt"

	concurrency "github.com/dcedyga/gooldi/concurrency"

	"time"
)

func (suite *Suite) Test03BCaster01ThreeClients() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)

	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	go clientFunc(1, b, dh1)
	go clientFunc2(2, b, dh1)
	go clientFunc3(3, b, dh1)
	time.Sleep(1 * time.Millisecond)
	// Start publishing messages:
	broadcast(b, 8000)
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

func (suite *Suite) Test03BCaster02Segregated4Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	c := make(chan interface{})
	defer func() {
		dm.GetDoneFunc()()
	}()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	// Create Processors
	s := createProcessors(nil, b, dh1, 4)
	// Get all multiplexed messages
	go getTotalMessagesFromProcessors(s, c)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)

}

func (suite *Suite) Test03BCaster03Segregated8Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	c := make(chan interface{})
	defer func() {
		dm.GetDoneFunc()()
	}()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	// Create Processors
	s := createProcessors(nil, b, dh1, 8)
	// Get all multiplexed messages
	go getTotalMessagesFromProcessors(s, c)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test03BCaster04Segregated16Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	c := make(chan interface{})
	defer func() {
		dm.GetDoneFunc()()
	}()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	// Create Processors
	s := createProcessors(nil, b, dh1, 16)
	// Get all multiplexed messages
	go getTotalMessagesFromProcessors(s, c)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test03BCaster05Segregated32Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	c := make(chan interface{})
	defer func() {
		dm.GetDoneFunc()()
	}()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	// Create Processors
	s := createProcessors(nil, b, dh1, 32)
	// Get all multiplexed messages
	go getTotalMessagesFromProcessors(s, c)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test03BCaster06Segregated64Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	c := make(chan interface{})
	defer func() {
		dm.GetDoneFunc()()
	}()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	// Create Processors
	s := createProcessors(nil, b, dh1, 64)
	// Get all multiplexed messages
	go getTotalMessagesFromProcessors(s, c)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test03BCaster07Segregated128Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	c := make(chan interface{})
	defer func() {
		dm.GetDoneFunc()()
	}()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	// Create Processors
	s := createProcessors(nil, b, dh1, 128)
	// Get all multiplexed messages
	go getTotalMessagesFromProcessors(s, c)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test03BCaster08Segregated256Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	c := make(chan interface{})
	defer func() {
		dm.GetDoneFunc()()
	}()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	// Create Processors
	s := createProcessors(nil, b, dh1, 256)
	// Get all multiplexed messages
	go getTotalMessagesFromProcessors(s, c)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test03BCaster09Segregated512Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	c := make(chan interface{})
	defer func() {
		dm.GetDoneFunc()()
	}()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	// Create Processors
	s := createProcessors(nil, b, dh1, 512)
	// Get all multiplexed messages
	go getTotalMessagesFromProcessors(s, c)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test03BCaster10Segregated1024Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	c := make(chan interface{})
	defer func() {
		dm.GetDoneFunc()()
	}()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	// Create Processors
	s := createProcessors(nil, b, dh1, 1024)
	// Get all multiplexed messages
	go getTotalMessagesFromProcessors(s, c)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test03BCaster11Segregated2048Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	c := make(chan interface{})
	defer func() {
		dm.GetDoneFunc()()
	}()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	// Create Processors
	s := createProcessors(nil, b, dh1, 2048)
	// Get all multiplexed messages
	go getTotalMessagesFromProcessors(s, c)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test03BCaster12Segregated10000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	c := make(chan interface{})
	defer func() {
		dm.GetDoneFunc()()
	}()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	// Create Processors
	s := createProcessors(nil, b, dh1, 10000)
	// Get all multiplexed messages
	go getTotalMessagesFromProcessors(s, c)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test03BCaster13Segregated20000Processors() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	c := make(chan interface{})
	defer func() {
		dm.GetDoneFunc()()
	}()
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)
	//Create caster
	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	// Create Processors
	s := createProcessors(nil, b, dh1, 20000)
	// Get all multiplexed messages
	go getTotalMessagesFromProcessors(s, c)
	// Broadcast Messages
	broadcast(b, 500000)
	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test03BCaster14ThreeClientsRemoveListener() {
	dm := concurrency.NewDoneManager(concurrency.DoneManagerWithDelay(1 * time.Millisecond))
	dh := dm.AddNewDoneHandler(0)
	dh1 := dm.AddNewDoneHandler(1)

	b := concurrency.NewBCaster(dh,
		"bcaster",
	)
	fmt.Printf("BCaster.ID():%v\n", b.ID())
	go clientFunc(1, b, dh1)
	go clientFunc2(2, b, dh1)
	go clientFunc3(3, b, dh1)
	time.Sleep(1 * time.Millisecond)

	l := b.AddListener(dh)
	b.RemoveListener(l)
	// Start publishing messages:
	broadcast(b, 8000)

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
