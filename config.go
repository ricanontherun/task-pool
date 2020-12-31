package task

type config struct {
	concurrency int
	debug       bool
	taskFunc    func(task Task)
}

func (config *config) TaskFunc() func(task Task) {
	return config.taskFunc
}

func (config *config) SetTaskFunc(taskFunc func(task Task)) {
	config.taskFunc = taskFunc
}

func (config *config) Concurrency() int {
	return config.concurrency
}

func (config *config) SetConcurrency(concurrency int) {
	config.concurrency = concurrency
}

func (config *config) Debug() bool {
	return config.debug
}

func (config *config) SetDebug(debug bool) {
	config.debug = debug
}

func NewConfig() *config {
	return &config{
		concurrency: 1,
		debug:       false,
		taskFunc:    func(task Task) {},
	}
}
