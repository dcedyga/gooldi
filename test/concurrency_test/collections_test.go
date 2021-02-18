package concurrency_test

import (
	"fmt"
	concurrency "goold/concurrency"
	"time"
)

func (suite *Suite) Test06Collections01SortedMapIterWithCancel() {

	m := concurrency.NewSortedMap()
	cancel := make(chan interface{})
	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}
	for item := range m.IterWithCancel(cancel) {
		fmt.Printf("%v:%v\n", item.Key, item.Value)
		if item.Key.(int) == 5 {
			close(cancel)
			break
		}
	}

	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test06Collections02SortedMapIterWithCancelNoExit() {

	m := concurrency.NewSortedMap()
	cancel := make(chan interface{})
	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}
	for item := range m.IterWithCancel(cancel) {
		fmt.Printf("%v:%v\n", item.Key, item.Value)
	}

	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test06Collections03SortedMapIterWithCancelExitDeleteAndReplay() {

	m := concurrency.NewSortedMap()
	cancel := make(chan interface{})
	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}
	keyToDelete := -1
	for item := range m.IterWithCancel(cancel) {
		if item.Key.(int) == 5 {
			keyToDelete = item.Key.(int)
			close(cancel)
			break
		}
	}

	m.Delete(keyToDelete)
	for item := range m.Iter() {
		fmt.Printf("%v:%v\n", item.Key, item.Value)
	}
	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test06Collections04MapIterWithCancel() {

	m := concurrency.NewMap()
	cancel := make(chan interface{})
	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}
	for item := range m.IterWithCancel(cancel) {
		fmt.Printf("%v:%v\n", item.Key, item.Value)
		if item.Key.(int) == 5 {
			close(cancel)
			break
		}
	}

	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test06Collections05MapIterWithCancelNoExit() {

	m := concurrency.NewMap()
	cancel := make(chan interface{})
	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}
	for item := range m.IterWithCancel(cancel) {
		fmt.Printf("%v:%v\n", item.Key, item.Value)
	}

	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test06Collections06MapIterWithCancelExitDeleteAndReplay() {

	m := concurrency.NewMap()
	cancel := make(chan interface{})
	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}
	keyToDelete := -1
	for item := range m.IterWithCancel(cancel) {
		if item.Key.(int) == 5 {
			keyToDelete = item.Key.(int)
			close(cancel)
			break
		}
	}

	m.Delete(keyToDelete)
	for item := range m.Iter() {
		fmt.Printf("%v:%v\n", item.Key, item.Value)
	}
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test06Collections07SliceIterWithCancel() {

	s := concurrency.NewSlice()
	cancel := make(chan interface{})
	s.Append("Juan")
	s.Append("Pepi")
	s.Append("David")
	s.Append("Ayleen")
	s.Append("Lila")
	s.Append("Freddy")
	s.Append("Moncho")
	s.Append("Zac")
	s.Append("Caty")
	s.Append("Tom")
	for item := range s.IterWithCancel(cancel) {
		fmt.Printf("%v:%v\n", item.Index, item.Value)
		if item.Index == 5 {
			close(cancel)
			break
		}
	}

	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test06Collections08SliceIterWithCancelNoExit() {

	s := concurrency.NewSlice()
	cancel := make(chan interface{})
	s.Append("Juan")
	s.Append("Pepi")
	s.Append("David")
	s.Append("Ayleen")
	s.Append("Lila")
	s.Append("Freddy")
	s.Append("Moncho")
	s.Append("Zac")
	s.Append("Caty")
	s.Append("Tom")
	for item := range s.IterWithCancel(cancel) {
		fmt.Printf("%v:%v\n", item.Index, item.Value)
	}

	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test06Collections09SliceIterWithCancelExitDeleteAndReplay() {

	s := concurrency.NewSlice()
	cancel := make(chan interface{})
	s.Append("Juan")
	s.Append("Pepi")
	s.Append("David")
	s.Append("Ayleen")
	s.Append("Lila")
	s.Append("Freddy")
	s.Append("Moncho")
	s.Append("Zac")
	s.Append("Caty")
	s.Append("Tom")
	keyToDelete := -1
	for item := range s.IterWithCancel(cancel) {
		if item.Index == 5 {
			keyToDelete = item.Index
			close(cancel)
			break
		}
	}

	s.RemoveItemAtIndex(keyToDelete)
	for item := range s.Iter() {
		fmt.Printf("%v:%v\n", item.Index, item.Value)
	}
	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test06Collections10SortedSliceIterWithCancel() {

	s := concurrency.NewSortedSlice()
	cancel := make(chan interface{})
	s.Append("Juan")
	s.Append("Pepi")
	s.Append("David")
	s.Append("Ayleen")
	s.Append("Lila")
	s.Append("Freddy")
	s.Append("Moncho")
	s.Append("Zac")
	s.Append("Caty")
	s.Append("Tom")
	for item := range s.IterWithCancel(cancel) {
		fmt.Printf("%v:%v\n", item.Index, item.Value)
		if item.Index == 5 {
			close(cancel)
			break
		}
	}

	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test06Collections11SortedSliceIterWithCancelNoExit() {

	s := concurrency.NewSortedSlice()
	cancel := make(chan interface{})
	s.Append("Juan")
	s.Append("Pepi")
	s.Append("David")
	s.Append("Ayleen")
	s.Append("Lila")
	s.Append("Freddy")
	s.Append("Moncho")
	s.Append("Zac")
	s.Append("Caty")
	s.Append("Tom")
	for item := range s.IterWithCancel(cancel) {
		fmt.Printf("%v:%v\n", item.Index, item.Value)
	}

	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test06Collections12SortedSliceIterWithCancelExitDeleteAndReplay() {

	s := concurrency.NewSortedSlice()
	cancel := make(chan interface{})
	s.Append("Juan")
	s.Append("Pepi")
	s.Append("David")
	s.Append("Ayleen")
	s.Append("Lila")
	s.Append("Freddy")
	s.Append("Moncho")
	s.Append("Zac")
	s.Append("Caty")
	s.Append("Tom")
	keyToDelete := -1
	for item := range s.IterWithCancel(cancel) {
		if item.Index == 5 {
			keyToDelete = item.Index
			close(cancel)
			break
		}
	}

	s.RemoveItemAtIndex(keyToDelete)
	for item := range s.Iter() {
		fmt.Printf("%v:%v\n", item.Index, item.Value)
	}
	time.Sleep(1000 * time.Millisecond)
}
