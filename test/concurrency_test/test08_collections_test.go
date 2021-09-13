package concurrency_test

import (
	"fmt"
	"time"

	concurrency "github.com/dcedyga/gooldi/concurrency"
)

func (suite *Suite) Test08Collections01SortedMapIterWithCancel() {

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

func (suite *Suite) Test08Collections02SortedMapIterWithCancelNoExit() {

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

func (suite *Suite) Test08Collections03SortedMapIterWithCancelExitDeleteAndReplay() {

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

func (suite *Suite) Test08Collections04SortedMapGetItemAtIndex() {

	m := concurrency.NewSortedMap()
	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}
	if item, ok := m.GetSortedMapItemByIndex(5); ok {
		fmt.Printf("key:%v,item:%v\n", item.Key, item.Value)
	}

}

func (suite *Suite) Test08Collections05MapIterWithCancel() {

	m := concurrency.NewMap()
	cancel := make(chan interface{})
	for i := 0; i < 10; i++ {
		go m.Set(i, i)
	}
	for item := range m.IterWithCancel(cancel) {
		fmt.Printf("%v:%v\n", item.Key, item.Value)

		if item.Key.(int) == 5 {
			close(cancel)
			key, _ := m.GetKeyByItem(item.Value)
			fmt.Printf("GetKeyByItem:%v:%v\n", key, item.Value)

			break
		}
	}

	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test08Collections06MapIterWithCancelNoExit() {

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

func (suite *Suite) Test08Collections07MapIterWithCancelExitDeleteAndReplay() {

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
func (suite *Suite) Test08Collections08SliceIterWithCancel() {

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

func (suite *Suite) Test08Collections09SliceIterWithCancelNoExit() {

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

func (suite *Suite) Test08Collections10SliceIterWithCancelExitDeleteAndReplay() {

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
func (suite *Suite) Test08Collections11SortedSliceIterWithCancel() {

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

func (suite *Suite) Test08Collections12SortedSliceIterWithCancelNoExit() {

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

func (suite *Suite) Test08Collections13SortedSliceIterWithCancelExitDeleteAndReplay() {

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

func (suite *Suite) Test08Collections14MapIterByKeys() {

	m := concurrency.NewMap()
	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}
	keys := m.GetKeys()
	for k := range keys {
		it, _ := m.Get(k)
		key, _ := m.GetKeyByItem(it)
		fmt.Printf("%v:%v\n", key, it)
	}

	time.Sleep(1000 * time.Millisecond)
}

func (suite *Suite) Test08Collections15SortedMapIterByKeys() {

	m := concurrency.NewSortedMap()
	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}
	keys := m.GetKeys()
	for _, idx := range keys {
		k, _ := m.GetKeyByIndex(idx.(int))
		it, _ := m.Get(k)
		key, _ := m.GetKeyByItem(it)
		fmt.Printf("%v:%v\n", key, it)
	}

	time.Sleep(1000 * time.Millisecond)
}
func (suite *Suite) Test08Collections16SortedSliceIndexOf() {
	s := concurrency.NewSlice()
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

	fmt.Printf("len:%v,cap:%v,indexOf(Tom):%v,GetItemAtIndex(5):%v\n", s.Len(), s.Cap(), s.IndexOf("Tom"), s.GetItemAtIndex(5))

}
func (suite *Suite) Test08Collections17SliceIndexOf() {
	s := concurrency.NewSortedSlice()
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

	fmt.Printf("len:%v,cap:%v,indexOf(Tom):%v,GetItemAtIndex(5):%v\n", s.Len(), s.Cap(), s.IndexOf("Tom"), s.GetItemAtIndex(5))

}
