package common

import (
	"fmt"
)

type queueNode struct {
	backward *queueNode
	value    interface{}
}

// FIFO queue based on the C++ std::queue class.
type Queue struct {
	head, tail *queueNode
	Size       int
}

// Queue adds a data structure to the back of the queue
func (stk *Queue) Queue(a interface{}) {
	node := queueNode{stk.head, a}
	var tmp *queueNode = stk.tail
	stk.tail = &node
	if tmp != nil {
		tmp.backward = stk.tail
	}
	if stk.head == nil {
		stk.head = &node
	}
	stk.Size++
}

// Dequeue removes the least recently added data structure: the structure at the front of the queue.
// After returning the data structure, the queue's head pointer is updated to the data structure immediately following the one just removed.
func (stk *Queue) Dequeue() (interface{}, error) {
	if stk.head != nil {
		tmp := stk.head
		if stk.head == stk.tail {
			stk.head = nil
		} else {
			stk.head = stk.head.backward
		}
		stk.tail.backward = stk.head
		stk.Size--
		return tmp.value, nil
	}
	return 0, fmt.Errorf("stack empty")
}

// IsEmpty returns true if there are no data structures in the queue.
// This will be true if the same number of items queued have been dequeued.
func (stk *Queue) IsEmpty() bool {
	if stk.head == nil || stk.tail == nil {
		return true
	}
	return false
}

// Peek will return the value at the front of the queue: the next item to be removed.
// However, Peek will not remove the data structure from the queue.
// It can be seen again after calling the function.
func (stk *Queue) Peek() interface{} {
	return stk.head.value
}

// Array walks through the queue and appends each element to an array to be returned.
func (stk *Queue) Array() []interface{} {
	node := stk.head
	list := make([]interface{}, stk.Size)
	for i := range list {
		if node == nil {
			break
		}

		list[i] = node.value
		node = node.backward
	}

	return list
}
