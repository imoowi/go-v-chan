package vch

import (
	"sync"
)

type ChCounter struct {
	lock *sync.RWMutex
	// counter map[string]int
	count int
}

func NewChCounter() *ChCounter {
	_counter := &ChCounter{
		lock: &sync.RWMutex{},
		// counter: make(map[string]int),
		count: 0,
	}
	// _counter.counter[key] = 0
	return _counter
}
func (ct *ChCounter) ChCounter() int {
	// ct.lock.RLock()
	// defer ct.lock.RUnlock()
	// return ct.counter[key]
	return ct.count
}
func (ct *ChCounter) Incr() {
	// ct.lock.Lock()
	// defer ct.lock.Unlock()
	// ct.counter[key]++
	ct.count++
	// log.Println(`ch[`, key, `] ++, length=`, ct.counter[key])
}

func (ct *ChCounter) Decr() {
	// ct.lock.Lock()
	// defer ct.lock.Unlock()
	// ct.counter[key]--
	ct.count--
	// log.Println(`ch[`, key, `] --, length=`, ct.counter[key])
}
