package task

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool_SingleWorker(t *testing.T) {
	var tasksToGenerate int32 = 10

	wg := new(sync.WaitGroup)
	var counter int32

	config := NewConfig()
	config.SetTaskFunc(func(task Task) {
		defer func() {
	 		atomic.AddInt32(&counter, 1)
			wg.Done()
		}()
	})

	worker, err := NewWorkerPool(config)
	assert.Nil(t, err, "err should not be nil")

	for i := int32(0); i < tasksToGenerate; i++ {
		worker.AddTask(NewTask())
	}

	wg.Add(int(tasksToGenerate))
	worker.Start()

	wg.Wait()
	assert.Equal(t, tasksToGenerate, counter, fmt.Sprintf("%d task should have been completed", tasksToGenerate))
}

func TestWorkerPool_ManyWorkers(t *testing.T) {
	var tasksToGenerate int32 = 100

	wg := new(sync.WaitGroup)
	var counter int32

	config := NewConfig()
	config.SetConcurrency(5)
	config.SetTaskFunc(func(task Task) {
		defer func() {
			atomic.AddInt32(&counter, 1)
			wg.Done()
		}()
	})

	worker, err := NewWorkerPool(config)
	assert.Nil(t, err, "err should not be nil")

	for i := int32(0); i < tasksToGenerate; i++ {
		worker.AddTask(NewTask())
	}

	wg.Add(int(tasksToGenerate))
	worker.Start()

	wg.Wait()
	assert.Equal(t, tasksToGenerate, counter, fmt.Sprintf("%d task should have been completed", tasksToGenerate))
}

func TestWorkerPool_AddTask_Timeouts(t *testing.T) {
	var tasksToGenerate int32 = 100
	var timeoutsExpected = tasksToGenerate

	wg := new(sync.WaitGroup)
	var counter int32
	var timeoutCounter int32
	timeout := time.Millisecond * 100

	config := NewConfig()
	config.SetConcurrency(5)
	config.SetTaskFunc(func(task Task) {
		defer func() {
			atomic.AddInt32(&counter, 1)
			wg.Done()
		}()

		timer := time.NewTimer(time.Millisecond * 200)
		select {
		case <-timer.C: // shouldn't happen.
			t.Error("timer shouldn't have ended")
			t.FailNow()
		case <-task.Cancelled():
			atomic.AddInt32(&timeoutCounter, 1)
			break
		}
	})

	worker, err := NewWorkerPool(config)
	assert.Nil(t, err, "err should not be nil")

	for i := int32(0); i < tasksToGenerate; i++ {
		task := NewTask()
		task.SetTimeout(timeout)
		worker.AddTask(task)
	}

	wg.Add(int(tasksToGenerate))
	worker.Start()

	wg.Wait()
	assert.Equal(t, tasksToGenerate, counter, fmt.Sprintf("%d task should have been completed", tasksToGenerate))
	assert.Equal(t, timeoutsExpected, timeoutCounter, "All tasks should have timed out")
}

func TestWorkerPool_AddTask_Timeouts_B(t *testing.T) {
	var tasksToGenerate int32 = 100

	wg := new(sync.WaitGroup)
	var counter int32
	timeout := time.Millisecond * 100

	config := NewConfig()
	config.SetConcurrency(5)
	config.SetTaskFunc(func(task Task) {
		defer func() {
			atomic.AddInt32(&counter, 1)
			wg.Done()
		}()

		timer := time.NewTimer(time.Millisecond * 50)
		select {
		case <-timer.C: // expected.
			break
		case <-task.Cancelled():
			t.Error("task should have not timed out.")
			t.FailNow()
		}
	})

	worker, err := NewWorkerPool(config)
	assert.Nil(t, err, "err should not be nil")

	for i := int32(0); i < tasksToGenerate; i++ {
		task := NewTask()
		task.SetTimeout(timeout)
		worker.AddTask(task)
	}

	wg.Add(int(tasksToGenerate))
	worker.Start()

	wg.Wait()
	assert.Equal(t, tasksToGenerate, counter, fmt.Sprintf("%d task should have been completed", tasksToGenerate))
}