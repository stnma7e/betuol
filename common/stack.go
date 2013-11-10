package common

type StackNode struct {
	forward *StackNode
	value   interface{}
}

type Stack struct {
	head, tail *StackNode
	Size       int
}

func (stk *Stack) Top() interface{} {
	return stk.tail.value
}

func (stk *Stack) Push(val interface{}) {
	sn := StackNode{
		stk.tail,
		val,
	}
	stk.tail = &sn
	if sn.forward == nil {
		stk.head = &sn
	}
	stk.Size++
}

func (stk *Stack) Pop() interface{} {
	if stk.tail == nil {
		return nil
	}
	tmp := stk.tail
	stk.tail = stk.tail.forward
	stk.Size--
	return tmp.value
}

func (stk *Stack) IsEmpty() bool {
	return stk.tail == nil || stk.head == nil
}
