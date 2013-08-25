package common

import "fmt"

type IntQueueNode struct {
	backward *IntQueueNode
	value int
}
type IntQueue struct {
	head, tail *IntQueueNode
}

func (stk *IntQueue) Queue(num int) {
	node := IntQueueNode{stk.head, num}
	var tmp *IntQueueNode = stk.tail
	stk.tail = &node
	if tmp != nil {
		tmp.backward = stk.tail
	}
	if (stk.head == nil) {
		stk.head = &node
	}

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

func (stk *IntQueue) WalkQueue() {
	stk.walkQueue(stk.head)
}
func (stk *IntQueue) walkQueue(node *IntQueueNode) {
	if node != nil {
		fmt.Println(node,"\b,")
		if node == stk.tail {
			fmt.Println()
			return
		}
		stk.walkQueue(node.backward)
	}
	return
}

func (stk *IntQueue) Size() uint {
	node := stk.head
	var i uint
	for ; ; i++ {
		if node == nil {
			break
		}
	}

	return i
}

func (stk *IntQueue) ToList() []int {
	node := stk.head
	list := make([]int, stk.Size())
	for i := range list {
		if node == nil {
			break
		}

		list[i] = node.value
	}

	return list
}