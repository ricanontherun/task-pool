package queue

import (
	"sync"
)

// A blocking queue is a simple FIFO queue
// which blocks on Get() until an item is available.
type blockingQueue struct {
	lock   sync.RWMutex
	buffer []interface{}
}

func NewBlockingQueue() Queue {
	// TODO: Possible optimization, initialize buffer with size to prevent reallocation.
	return &blockingQueue{
		buffer: make([]interface{}, 0, 0),
	}
}

func (q *blockingQueue) Add(item interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.buffer = append(q.buffer, item)
}

func (q *blockingQueue) Get() interface{} {
	var item interface{}

	for {
		// TODO: Possible CPU hotspot
		q.lock.RLock()
		if len(q.buffer) == 0 {
			q.lock.RUnlock()
			continue
		} else {
			item = q.buffer[len(q.buffer)-1]
			q.lock.RUnlock()
			break
		}
	}

	q.lock.Lock()
	q.buffer = q.buffer[1:]
	q.lock.Unlock()

	return item
}

func (q *blockingQueue) Stats() {
}

func (q *blockingQueue) Length() int {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return len(q.buffer)
}
