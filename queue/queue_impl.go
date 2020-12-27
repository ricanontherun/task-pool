package queue

import (
	"sync"
)

type atomicBlockingQueue struct {
	lock     sync.RWMutex
	buffer   []interface{}
	buffered AtomicBool
}

func NewBlockingQueue() Queue {
	return &atomicBlockingQueue{
		buffer:   make([]interface{}, 0, 0),
		buffered: NewAtomicBool(false),
	}
}

func (C *atomicBlockingQueue) Add(item interface{}) {
	C.lock.Lock()
	defer C.lock.Unlock()

	C.buffer = append(C.buffer, item)
	C.buffered.Set(true)
}

func (C *atomicBlockingQueue) AddMany(items []interface{}) {
	C.lock.Lock()
	defer C.lock.Unlock()

	for _, item := range items {
		C.buffer = append(C.buffer, item)
	}

	C.buffered.Set(true)
}

func (C *atomicBlockingQueue) Get() interface{} {
	for {
		if C.buffered.Get() {
			C.lock.RLock()
			item := C.buffer[0]
			C.lock.RUnlock()

			C.lock.Lock()
			C.buffer = C.buffer[1:]
			C.lock.Unlock()

			C.lock.RLock()
			C.buffered.Set(len(C.buffer) > 0)
			C.lock.RUnlock()

			return item
		}
	}
}

func (C *atomicBlockingQueue) Length() int {
	C.lock.RLock()
	defer C.lock.RUnlock()

	return len(C.buffer)
}

func (C *atomicBlockingQueue) Stats() {
}
