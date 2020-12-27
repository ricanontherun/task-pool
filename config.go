package task

type config struct {
	Concurrency int
	Debug       bool
	TaskFunc    func(task Task)
}

type configBuilder struct {
	concurrency int
	taskFunc    func(task Task)
	debug       bool
}

func (cb *configBuilder) SetConcurrency(concurrency int) *configBuilder {
	cb.concurrency = concurrency
	return cb
}

func (cb *configBuilder) SetDebug(debug bool) *configBuilder {
	cb.debug = debug
	return cb
}

func (cb *configBuilder) SetTaskFunc(tf func(task Task)) *configBuilder {
	cb.taskFunc = tf
	return cb
}

func (cb *configBuilder) Build() config {
	return config{
		Concurrency: cb.concurrency,
		Debug:       cb.debug,
		TaskFunc:    cb.taskFunc,
	}
}

func NewConfigBuilder() *configBuilder {
	cb := &configBuilder{}
	cb.concurrency = 1
	cb.debug = false
	cb.taskFunc = func(task Task) {}
	return cb
}
