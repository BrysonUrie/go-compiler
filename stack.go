package main

import "fmt"

// Stack defines a generic stack
type Stack[T any] struct {
	items []T
}

// Push adds an element to the stack
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top element of the stack
func (s *Stack[T]) Pop() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, fmt.Errorf("stack is empty")
	}

	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]

	return top, nil
}

// Peek returns the top element without removing it
func (s *Stack[T]) Peek() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, fmt.Errorf("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

// IsEmpty checks if the stack is empty
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}
