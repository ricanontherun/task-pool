package task

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestTimeouts(t *testing.T) {
	// Single task
	// Wait for it with a wait group
	// write to some global variable within the task
	// Timeouts should prevent task from being
	var wg sync.WaitGroup
	pool, err := NewWorkerPool(NewConfigBuilder().SetTaskFunc(func(task Task) {
		fmt.Println("in task function")

		select {
		case <-task.Cancelled():
			fmt.Println("I've been cancelled")
			break
		case <-time.NewTicker(time.Second * 10).C:
			fmt.Println("my 2 second work is finished")
			break
		}

		wg.Done()
	}).Build())

	assert.Nil(t, err, "err should be nil")
	wg.Add(1)

	newTask := NewTask()
	newTask.SetTimeout(time.Second * 5)
	pool.AddTask(newTask)
	pool.Start()
	wg.Wait()
}

func TestThing(t *testing.T) {
	cancelledChannel := make(chan bool, 1)

	go func() {
		time.Sleep(time.Second * 1)
		// Golang allows for channels to be written to and not consumed.
		// as this doesn't result in a deadlock.
		// Golang goes not allow channels to consumed from and NOT written to.
		// This is literally the cause of deadlocks. A thread waiting for something
		// that cannot happen.
		cancelledChannel <- true
	}()

	fmt.Println("ok then")
}
