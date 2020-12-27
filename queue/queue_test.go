package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	queue := NewBlockingQueue()

	if queue.Length() != 0 {
		t.Error("queue should be empty.")
	}

	queue.Add(1)
	if queue.Length() != 1 {
		t.Error("queue should contain 1 item")
	}

	if item := queue.Get(); item.(int) != 1 {
		t.Error("queue item should have been 1")
	}

	if queue.Length() != 0 {
		t.Error("queue should be empty.")
	}

	for i := 0; i < 10; i++ {
		queue.Add(i)
	}

	if queue.Length() != 10 {
		t.Error("queue should contain 10 elements")
	}

	for i := 0; i < 10; i++ {
		if item := queue.Get(); item.(int) != i {
			t.Errorf("queue item should have equaled %d, was %d", i, item.(int))
		}
	}

	if queue.Length() != 0 {
		t.Error("queue should be empty")
	}

	things := make([]interface{}, 100)
	for i := 0; i < len(things); i++ {
		things[i] = 10
	}
	queue.AddMany(things)

	assert.EqualValues(t, 100, queue.Length())
}
