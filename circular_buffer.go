package generic

import (
	"errors"
	"sync"
)

var ErrIndexOutOfRange = errors.New("index out of range")

type CircularBuffer[T any] struct {
	buffer []T
	size int
	head int
	mutex *sync.Mutex
}

func NewCircularBuffer[T any](size int) *CircularBuffer[T] {
	return &CircularBuffer[T]{
		buffer: make([]T, size),
		size: size,
		head: 0,
		mutex: &sync.Mutex{},
	}
}

func (cb *CircularBuffer[T]) Append(val T) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	idx := cb.head % cb.size
	cb.buffer[idx] = val
	cb.head += 1
}

func (cb *CircularBuffer[T]) Get(idx int) (val T, err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	if idx < 0 || idx >= cb.head || idx < cb.head - cb.size {
		return val, ErrIndexOutOfRange
	}
	return cb.buffer[idx % cb.size], nil
}

func (cb *CircularBuffer[T]) Slice() []T {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	if cb.head < cb.size {
		slice := make([]T, cb.head)
		copy(slice, cb.buffer)
		return slice
	}
	idx := cb.head % cb.size
	slice := make([]T, cb.size)
	copy(slice, cb.buffer[:idx])
	copy(slice[idx:], cb.buffer)
	return slice
}
