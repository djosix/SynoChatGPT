package internal

import "sync"

type Queue[T any] struct {
	idx  int
	len  int
	data []T
	lock sync.Mutex
}

func NewQueue[T any](size int) Queue[T] {
	return Queue[T]{data: make([]T, size)}
}

func (q *Queue[T]) Push(value T) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.data[q.idx] = value
	q.idx = (q.idx + 1) % len(q.data)
	if q.len < len(q.data) {
		q.len++
	}
}

func (q *Queue[T]) Clear() {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.len = 0
}

func (q *Queue[T]) Data() []T {
	q.lock.Lock()
	defer q.lock.Unlock()
	startIdx := (q.idx - q.len + len(q.data)) % len(q.data)
	result := []T{}
	if startIdx >= q.idx {
		result = append(result, q.data[startIdx:len(q.data)]...)
		result = append(result, q.data[0:q.idx]...)
	} else {
		result = append(result, q.data[startIdx:q.idx]...)
	}
	return result
}
