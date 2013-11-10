package common

import (
	"fmt"
)

type queueNode struct {
	backward *queueNode
	value    interface{}
}

type Queue struct {
	head, tail *queueNode
	Size       int
}

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

func (stk *Queue) IsEmpty() bool {
	if stk.head == nil || stk.tail == nil {
		return true
	}
	return false
}

func (stk *Queue) Peek() interface{} {
	return stk.head.value
}

func (stk *Queue) WalkQueue() {
	stk.walkQueue(stk.head)
}
func (stk *Queue) walkQueue(node *queueNode) {
	if node != nil {
		fmt.Println(node, "\b,")
		if node == stk.tail {
			fmt.Println()
			return
		}
		stk.walkQueue(node.backward)
	}
}

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
