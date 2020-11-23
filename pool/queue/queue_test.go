package queue

import (
	"testing"
)

func TestQueue(t *testing.T) {
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
}