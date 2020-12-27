package task

import (
	"sync"
	"time"
)

type task struct {
	data      map[string]interface{}
	timeout   time.Duration
	cancelled chan bool
	lock      sync.RWMutex
}

type Task interface {
	// Remove this in favor of casts and sets.
	SetTimeout(ttl time.Duration) Task
	GetTimeout() time.Duration
	Cancelled() <-chan bool
	Cancel()
}

func NewTask() Task {
	return &task{
		data:      make(map[string]interface{}),
		cancelled: make(chan bool, 1),
	}
}

func (task *task) SetTimeout(ttl time.Duration) Task {
	task.timeout = ttl
	return task
}

func (task *task) GetTimeout() time.Duration {
	return task.timeout
}

func (task *task) Cancelled() <-chan bool {
	task.lock.RLock()
	defer task.lock.RUnlock()

	return task.cancelled
}

func (task *task) Cancel() {
	task.lock.Lock()
	defer task.lock.Unlock()

	task.cancelled <- true
	close(task.cancelled)
}
