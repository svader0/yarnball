package stack

import "fmt"

/*
	A simple stack implementation for use in the Yarnball interpreter.
*/

type Stack struct {
	items []int
}

func New() *Stack {
	return &Stack{
		items: []int{},
	}
}

func (s *Stack) Push(item int) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() (int, error) {
	if len(s.items) == 0 {
		return 0, fmt.Errorf("stack underflow")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}

func (s *Stack) Peek() (int, bool) {
	if len(s.items) == 0 {
		return 0, false // Stack is empty
	}
	return s.items[len(s.items)-1], true
}

// PeekAt returns the value depth elements from the top (0 = top).
func (s *Stack) PeekAt(depth int) (int, error) {
	if depth < 0 || depth >= len(s.items) {
		return 0, fmt.Errorf("stack underflow")
	}
	return s.items[len(s.items)-1-depth], nil
}

// Roll moves the value at depth to the top (0 = no-op).
func (s *Stack) Roll(depth int) error {
	if depth < 0 || depth >= len(s.items) {
		return fmt.Errorf("stack underflow")
	}
	if depth == 0 {
		return nil
	}
	idx := len(s.items) - 1 - depth
	val := s.items[idx]
	copy(s.items[idx:], s.items[idx+1:])
	s.items[len(s.items)-1] = val
	return nil
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Size() int {
	return len(s.items)
}

func (s *Stack) Clear() {
	s.items = []int{}
}
