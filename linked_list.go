package generic

import (
	"errors"
	"sync"
)

var ErrNoData = errors.New("no data")
var ErrNotStarted = errors.New("iterator not yet started (call Next() first)")

type LinkedListElem[T any] struct {
	prev *LinkedListElem[T]
	next *LinkedListElem[T]
	data T
}

func (elem *LinkedListElem[T]) Next() *LinkedListElem[T] {
	return elem.next
}

func (elem *LinkedListElem[T]) Prev() *LinkedListElem[T] {
	return elem.prev
}

func (elem *LinkedListElem[T]) Data() T {
	return elem.data
}

type LinkedListIter[T any] struct {
	elem *LinkedListElem[T]
	started bool
	zero T
}

func (iter *LinkedListIter[T]) Next() bool {
	if iter.elem == nil {
		return false
	}
	if iter.started {
		iter.elem = iter.elem.Next()
	} else {
		iter.started = true
	}
	return iter.elem != nil
}

func (iter *LinkedListIter[T]) Get() (T, error) {
	if !iter.started {
		return iter.zero, ErrNotStarted
	}
	if iter.elem == nil {
		return iter.zero, ErrNoData
	}
	return iter.elem.Data(), nil
}

type LinkedList[T any] struct {
	head *LinkedListElem[T]
	tail *LinkedListElem[T]
	zero T
	size int
	mutex *sync.Mutex
}

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{mutex: &sync.Mutex{}}
}

func (ll *LinkedList[T]) Push(item T) {
	ll.mutex.Lock()
	defer ll.mutex.Unlock()
	if ll.tail == nil {
		ll.head = &LinkedListElem[T]{data: item}
		ll.tail = ll.head
	} else {
		elem := &LinkedListElem[T]{prev: ll.tail, data: item}
		ll.tail.next = elem
		ll.tail = elem
	}
	ll.size += 1
}

func (ll *LinkedList[T]) Unshift(item T) {
	ll.mutex.Lock()
	defer ll.mutex.Unlock()
	if ll.head == nil {
		ll.head = &LinkedListElem[T]{data: item}
		ll.tail = ll.head
	} else {
		elem := &LinkedListElem[T]{next: ll.head, data: item}
		ll.head.prev = elem
		ll.head = elem
	}
	ll.size += 1
}

func (ll *LinkedList[T]) Pop() T {
	ll.mutex.Lock()
	defer ll.mutex.Unlock()
	return ll.pop()
}

func (ll *LinkedList[T]) pop() T {
	if ll.tail == nil {
		return ll.zero
	}
	elem := ll.tail
	elem.prev.next = nil
	ll.tail = elem.prev
	ll.size -= 1
	return elem.data
}

func (ll *LinkedList[T]) PopIf(test func(T) bool) (T, bool) {
	ll.mutex.Lock()
	defer ll.mutex.Unlock()
	if ll.tail == nil {
		return ll.zero, false
	}
	if test(ll.tail.data) {
		return ll.pop(), true
	}
	return ll.tail.data, false
}

func (ll *LinkedList[T]) Shift() T {
	ll.mutex.Lock()
	defer ll.mutex.Unlock()
	return ll.shift()
}

func (ll *LinkedList[T]) shift() T {
	if ll.head == nil {
		return ll.zero
	}
	elem := ll.head
	elem.next.prev = nil
	ll.head = elem.next
	ll.size -= 1
	return elem.data
}

func (ll *LinkedList[T]) ShiftIf(test func(T) bool) (T, bool) {
	ll.mutex.Lock()
	defer ll.mutex.Unlock()
	if ll.head == nil {
		return ll.zero, false
	}
	if test(ll.head.data) {
		return ll.shift(), true
	}
	return ll.head.data, false
}

func (ll *LinkedList[T]) First() T {
	elem := ll.head
	if elem == nil {
		return ll.zero
	}
	return elem.data
}

func (ll *LinkedList[T]) Last() T {
	elem := ll.tail
	if elem == nil {
		return ll.zero
	}
	return elem.data
}

func (ll *LinkedList[T]) Len() int {
	return ll.size
}

func (ll *LinkedList[T]) Iter() *LinkedListIter[T] {
	return &LinkedListIter[T]{elem: ll.head}
}

func (ll *LinkedList[T]) Slice() []T {
	ll.mutex.Lock()
	defer ll.mutex.Unlock()
	if ll.head == nil {
		return nil
	}
	s := make([]T, ll.size)
	iter := ll.Iter()
	i := 0
	for iter.Next() {
		if i >= ll.size {
			break
		}
		val, err := iter.Get()
		if err != nil {
			break
		}
		s[i] = val
		i += 1
	}
	if i < ll.size {
		return s[:i]
	}
	return s
}
