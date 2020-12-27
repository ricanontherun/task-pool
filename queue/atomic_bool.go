package queue

import "sync/atomic"

type atomicBool struct {
	value int32
}

type AtomicBool interface {
	Get() bool
	Set(bool)
}

func NewAtomicBool(initialValue bool) AtomicBool {
	return &atomicBool{
		value: boolToInt(initialValue),
	}
}

func (a *atomicBool) Get() bool {
	return intToBool(atomic.LoadInt32(&a.value))
}

func (a *atomicBool) Set(newValue bool) {
	atomic.StoreInt32(&a.value, boolToInt(newValue))
}

// 0 -> false, else -> true
func intToBool(value int32) bool {
	if value == 0 {
		return false
	} else {
		return true
	}
}

// true -> 1, else -> 0
func boolToInt(value bool) int32 {
	if value == true {
		return 1
	} else {
		return 0
	}
}
