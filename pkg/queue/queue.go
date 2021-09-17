package queue

import "container/list"

type Queue[T any] struct {
	list *list.List
}

func New[T any]() *Queue[T] {
	return &Queue[T]{
		list: list.New(),
	}
}

func (q *Queue[T]) Push(v T) {
	q.list.PushBack(v)
}

func (q *Queue[T]) Pop() T {
	return q.list.Remove(q.list.Front()).(T)
}

func (q *Queue[T]) IsEmpty() bool {
	return q.list.Len() == 0
}
