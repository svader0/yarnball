package stack

import "fmt"

// Stack represents a simple stack structure for managing integers.
type Stack struct {
	items []int
}

// New creates a new Stack instance.
func New() *Stack {
	return &Stack{
		items: []int{},
	}
}

// Push adds an item to the top of the stack.
func (s *Stack) Push(item int) {
	s.items = append(s.items, item)
}

// Pop removes and returns the item from the top of the stack.
func (s *Stack) Pop() (int, error) {
	if len(s.items) == 0 {
		return 0, fmt.Errorf("stack underflow")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}

// Peek returns the item at the top of the stack without removing it.
func (s *Stack) Peek() (int, bool) {
	if len(s.items) == 0 {
		return 0, false // Stack is empty
	}
	return s.items[len(s.items)-1], true
}

// IsEmpty checks if the stack is empty.
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of items in the stack.
func (s *Stack) Size() int {
	return len(s.items)
}

// Clear removes all items from the stack.
func (s *Stack) Clear() {
	s.items = []int{}
}
