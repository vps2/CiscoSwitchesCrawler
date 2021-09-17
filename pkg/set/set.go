package set

import (
	"fmt"
)

type Set[T comparable] struct {
	container map[T]nothing
}

type nothing struct{}

func New[T comparable]() *Set[T] {
	return &Set[T]{
		container: make(map[T]nothing),
	}
}

func (s *Set[T]) Add(val T) {
	s.container[val] = nothing{}
}

func (s *Set[T]) Remove(val T) {
	delete(s.container, val)
}

func (s *Set[T]) Clear() {
	s.container = map[T]nothing{}
}

func (s *Set[T]) Has(val T) bool {
	_, ok := s.container[val]
	return ok
}

func (s *Set[T]) ToSlice() []T {
	var keys []T
	for key := range s.container {
		keys = append(keys, key)
	}

	return keys
}

func (s *Set[T]) Len() int {
	return len(s.container)
}

func (s *Set[T]) String() string {

	return fmt.Sprint(s.ToSlice())
}
