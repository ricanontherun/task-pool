package pool

type Config struct {
	Concurrency    int
	OnTaskComplete func()
}
