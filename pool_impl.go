package task

import (
	"context"
	"github.com/ricanontherun/task/queue"
	"go.uber.org/zap"
	"sync/atomic"
)

type workerPool struct {
	config         *config
	queue          queue.Queue
	workerReady    chan bool
	taskAvailable  chan Task
	tasksCompleted uint64
	logger         *zap.Logger
}

func (pool *workerPool) worker(workerId int) {
	pool.workerReady <- true
	pool.logger.Info("worker initialized, ready", zap.Int("workerId", workerId))

	for {
		pool.logger.Debug("worker waiting for next task", zap.Int("workerId", workerId))
		task := <-pool.taskAvailable
		pool.logger.Debug("worker received task", zap.Int("workerId", workerId))

		timeout := task.GetTimeout()
		if timeout.Nanoseconds() == 0 {
			pool.logger.Debug("starting task without timeout", zap.Int("workerId", workerId))
			pool.config.taskFunc(task)
		} else {
			pool.logger.Debug("starting request with timeout", zap.Duration("timeout", timeout))
			taskCtx, cancel := context.WithTimeout(context.Background(), timeout)
			doneChannel := make(chan bool, 1)

			go func(task Task, done chan bool) {
				defer func() {
					done <- true
				}()
				pool.config.taskFunc(task)
			}(task, doneChannel)

			// Wait for the task to finish naturally or timeout.
			select {
			case <-doneChannel:
				pool.logger.Debug("task with timeout finished in time")
			case <-taskCtx.Done():
				task.Cancel()
				pool.logger.Debug("task timed out", zap.Duration("timeout", timeout))
			}

			cancel()
		}

		pool.logger.Debug("worker completed task", zap.Int("workerId", workerId))
		atomic.AddUint64(&pool.tasksCompleted, 1)
		pool.workerReady <- true
	}
}

func NewWorkerPool(config *config) (WorkerPool, error) {
	workerPool := &workerPool{
		config:        config,
		workerReady:   make(chan bool, config.concurrency),
		taskAvailable: make(chan Task, config.concurrency),
		queue:         queue.NewBlockingQueue(),
	}

	var logger *zap.Logger
	var loggerErr error
	if config.debug {
		logger, loggerErr = zap.NewDevelopment()
	} else {
		logger, loggerErr = zap.NewProduction()
	}

	if loggerErr != nil {
		return nil, loggerErr
	} else {
		workerPool.logger = logger
	}

	return workerPool, nil
}

func (pool *workerPool) AddTask(task Task) {
	pool.queue.Add(task)
}

func (pool *workerPool) AddTasks(tasks []Task) {
	iTasks := make([]interface{}, len(tasks))
	for i := range tasks {
		iTasks[i] = tasks[i]
	}
	pool.queue.AddMany(iTasks)
}

// not thread safe.
func (pool *workerPool) Stats() Stats {
	queued := uint64(pool.queue.Length())
	completed := atomic.LoadUint64(&pool.tasksCompleted)

	return Stats{
		TasksCompleted: completed,
		TasksQueued:    queued,
		TasksAdded:     completed + queued,
	}
}

func (pool *workerPool) Start() {
	// Spawn the worker threads
	for i := 0; i < pool.config.concurrency; i++ {
		go pool.worker(i)
	}

	// task dispatch thread
	// dispatch tasks to workers are both become available
	go func() {
		for {
			// block until a task is available
			nextTask := pool.queue.Get().(Task)

			// block until a worker is ready
			<-pool.workerReady

			// dispatch task to worker
			pool.taskAvailable <- nextTask
		}
	}()
}
