package common

import "fmt"

type intQueueNode struct {
	backward *intQueueNode
	value int
}
type IntQueue struct {
	head, tail *intQueueNode
	Size int
}

func (stk *IntQueue) Queue(num int) {
	node := intQueueNode{stk.head, num}
	var tmp *intQueueNode = stk.tail
	stk.tail = &node
	if tmp != nil {
		tmp.backward = stk.tail
	}
	if (stk.head == nil) {
		stk.head = &node
	}
	stk.Size++
}
func (stk *IntQueue) Dequeue() (int,error) {
	if stk.head != nil {
		tmp := stk.head
		if stk.head == stk.tail {
			stk.head = nil
		} else {
			stk.head = stk.head.backward
		}
		stk.tail.backward = stk.head
		stk.Size--
		return tmp.value,nil
	}
	return 0,fmt.Errorf("stack empty")
}

func (stk *IntQueue) IsEmpty() bool {
	if stk.head == nil || stk.tail == nil {
		return true
	}
	return false
}

func (stk *IntQueue) Peek() int {
	return stk.head.value
}

func (stk *IntQueue) WalkQueue() {
	stk.walkQueue(stk.head)
}
func (stk *IntQueue) walkQueue(node *intQueueNode) {
	if node != nil {
		fmt.Println(node,"\b,")
		if node == stk.tail {
			fmt.Println()
			return
		}
		stk.walkQueue(node.backward)
	}
}

func (stk *IntQueue) ToList() []int {
	node := stk.head
	list := make([]int, stk.Size)
	for i := range list {
		if node == nil {
			break
		}

		list[i] = node.value
		node = node.backward
	}

	return list
}