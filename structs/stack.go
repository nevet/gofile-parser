package stack

type Stack struct {
	top  *Element
	size int
}

type Element struct {
	value interface{}
	next  *Element
}

func (stack *Stack) Empty() bool {
	return stack.size == 0
}

func (stack *Stack) Size() int {
	return stack.size
}

func (stack *Stack) Top() (value interface{}) {
	return stack.top.value
}

func (stack *Stack) Push(val interface{}) {
	stack.top = &Element{value: val, next: stack.top}
	stack.size++
}

func (stack *Stack) Pop() {
	if size > 0 {
		stack.top = stack.top.next
		size--
	}
}
