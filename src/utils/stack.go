package utils

import "fmt"

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(data T) {
	s.items = append(s.items, data)
}

func (s *Stack[T]) Pop() {
	if s.IsEmpty() {
		return
	}
	s.items = s.items[:len(s.items)-1]
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack[T]) Peek() (T, error) {
	if s.IsEmpty() {
		return s.items[0], fmt.Errorf("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

func (s *Stack[T]) Get(i int) (T, error) {
	if s.IsEmpty() {
		return s.items[0], fmt.Errorf("stack is empty")
	}
	return s.items[i], nil
}
