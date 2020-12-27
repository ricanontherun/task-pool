package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAtomicBool(t *testing.T) {
	ab1 := NewAtomicBool(false)
	assert.EqualValues(t, false, ab1.Get())
	ab1.Set(true)
	assert.EqualValues(t, true, ab1.Get())

	ab2 := NewAtomicBool(true)
	assert.EqualValues(t, true, ab2.Get())
	ab2.Set(false)
	assert.EqualValues(t, false, ab2.Get())
}
