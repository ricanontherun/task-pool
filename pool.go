package task

type WorkerPool interface {
	AddTask(task Task)
	AddTasks(tasks []Task)
	Start()
	Stats() Stats
}
