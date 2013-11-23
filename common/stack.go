package common

type stackNode struct {
	forward *stackNode
	value   interface{}
}

// Stack implements a LIFO stack similar to C++ std::stack
type Stack struct {
	head, tail *stackNode
	Size       int
}

// Top returns the value at the top of the stack without removing it from the stack.
func (stk *Stack) Top() interface{} {
	return stk.tail.value
}

// Push adds a data structure to the top of the stack.
func (stk *Stack) Push(val interface{}) {
	sn := stackNode{
		stk.tail,
		val,
	}
	stk.tail = &sn
	if sn.forward == nil {
		stk.head = &sn
	}
	stk.Size++
}

// Pop removes and returns the top data structure of the stack.
func (stk *Stack) Pop() interface{} {
	if stk.tail == nil {
		return nil
	}
	tmp := stk.tail
	stk.tail = stk.tail.forward
	stk.Size--
	return tmp.value
}

// IsEmpty returns true if there are no data structures in the stack.
// If the same number of data structures have been pushed as poped, then IsEmpty will return true.
func (stk *Stack) IsEmpty() bool {
	return stk.tail == nil || stk.head == nil
}
