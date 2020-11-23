package queue

type Queue interface {
	Add(item interface{})
	Get() interface{}
	Length() int
	Stats()
}
