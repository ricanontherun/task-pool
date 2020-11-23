package pool

type WorkerPool interface {
	AddTask(task Task)
	RemainingTasks() int
}
