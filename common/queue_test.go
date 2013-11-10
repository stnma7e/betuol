package common

import "testing"

//import "fmt"

func Test_IsEmpty(t *testing.T) {

	q := &IntQueue{}

	if !q.IsEmpty() {
		t.Error("Queue has no items, but IsEmpty() returned false.")
	}
}

func Test_AddIsNotEmpty(t *testing.T) {

	q := &IntQueue{}

	q.Queue(1)

	if q.IsEmpty() {
		t.Error("Stack should not be empty.")
	}

	value, error := q.Dequeue()

	if error != nil || value != 1 {
		t.Error("Unexpected value when dequeuing item from stack")
	}

	if !q.IsEmpty() {
		t.Error("Queue should be empty.")
	}
}

func Test_AddMultipleItems(t *testing.T) {

	q := &IntQueue{}

	q.Queue(1)
	q.Queue(15)
	q.Queue(1998)

	value1, error := q.Dequeue()

	t.Log("Item dequeued is ", value1)

	if error != nil {
		t.Error("Unexpected error dequeuing item.")
	} else if value1 != 1 {
		t.Error("Unexpected value when dequeuing item.")
	}

	value2, error := q.Dequeue()

	t.Log("Item dequeued is ", value2)

	if error != nil {
		t.Error("Unexpected error dequeuing item.")
	} else if value2 != 15 {
		t.Error("Unexpected value when dequeuing item.")
	}

	value3, error := q.Dequeue()

	t.Log("Item dequeued is ", value3)

	if error != nil {
		t.Error("Unexpected error dequeuing item.")
	} else if value3 != 1998 {
		t.Error("Unexpected value when dequeuing item.")
	}

	if !q.IsEmpty() {
		t.Error("Queue should be empty.")
	}
}

func Test_DeQueueWithNoItems(t *testing.T) {

	q := &IntQueue{}
	value, error := q.Dequeue()

	if !(error != nil && value == 0) {
		t.Error("Expected empty stack error and 0 for value.", error, value)
	}

}

func Test_WalkQueue(t *testing.T) {

	q := &IntQueue{}
	q.Queue(133)
	q.Queue(999)
	q.WalkQueue()
	q.WalkQueue()

	// TODO : make WalkQueue testable.
	// Make walkQueue return string
	// test that values of inseted numbers match with the output of
	// walkQueue
}

func Test_OverDequeue(t *testing.T) {

	q := &IntQueue{}
	q.Queue(133)
	q.Queue(999)

	q.Dequeue()
	q.Dequeue()
	value, error := q.Dequeue()

	if !(error != nil && value == 0) {
		t.Error("Expected empty stack error and 0 for value.", error, value)
	}

}
