package queue

type Queue interface {
	Add(item interface{})
	AddMany(item []interface{})
	Get() interface{}
	Length() int
	Stats()
}
