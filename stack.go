package main

type Stack[T any] []T

func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack[T]) Push(elem T) {
	*s = append(*s, elem)
}

func (s *Stack[T]) Pop() (any, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}
