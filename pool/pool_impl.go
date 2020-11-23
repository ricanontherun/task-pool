package pool

import (
	"log"
	"task-pool/pool/queue"
	"time"
)

type workerPool struct {
	config        Config
	queue         queue.Queue
	workerReady   chan bool
	taskAvailable chan Task
}

func worker(id int, readyChan chan<- bool, workChan <-chan Task, onComplete func()) {
	// Workers should be ready immediately.
	readyChan <- true

	for {
		log.Printf("worker %d waiting for tasks\n", id)
		task := <-workChan
		log.Printf("worker %d received task, work for %d seconds\n", id, task.Sleep)
		time.Sleep(time.Second * time.Duration(task.Sleep))
		log.Printf("worker %d finished!\n", id)

		onComplete()

		// alert dispatcher we're workerReady for more taskAvailable
		readyChan <- true
	}
}

func NewWorkerPool(config Config) WorkerPool {
	workerPool := &workerPool{
		config:        config,
		workerReady:   make(chan bool, config.Concurrency),
		taskAvailable: make(chan Task, config.Concurrency),
		queue:         queue.NewBlockingQueue(),
	}

	// Spawn the worker threads.
	for i := 0; i < config.Concurrency; i++ {
		go worker(i, workerPool.workerReady, workerPool.taskAvailable, config.OnTaskComplete)
	}

	// task dispatch thread.
	// dispatch tasks to workers are both become available.
	go func() {
		for {
			// block until a task is available.
			nextTask := workerPool.queue.Get().(Task)

			// block until a worker is workerReady.
			<-workerPool.workerReady

			// Make this task available to the worker.
			workerPool.taskAvailable <- nextTask
		}
	}()

	return workerPool
}

func (pool *workerPool) AddTask(task Task) {
	pool.queue.Add(task)
}

func (pool *workerPool) RemainingTasks() int {
	return pool.queue.Length()
}
