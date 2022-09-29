package generic

import (
	"math"
	"sync"
)

type Number interface {
	int | int64 | int32 | int16 | int8 | uint | uint64 | uint32 | uint16 | uint8 | float64 | float32
}

type NumericCircularBuffer[T Number] struct {
	buffer []T
	size int
	head int
	mutex *sync.Mutex
	count float64
	sum float64
	sqsum float64
	min float64
	max float64
}

func NewNumericCircularBuffer[T Number](size int) *NumericCircularBuffer[T] {
	return &NumericCircularBuffer[T]{
		buffer: make([]T, size),
		size: size,
		head: 0,
		mutex: &sync.Mutex{},
		min: math.Inf(1),
		max: math.Inf(-1),
	}
}

func (cb *NumericCircularBuffer[T]) Append(val T) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	old := cb.head - cb.size
	idx := cb.head % cb.size
	if old >= 0 {
		oldVal := float64(cb.buffer[idx])
		if !math.IsNaN(oldVal) {
			cb.count -= 1
			cb.sum -= oldVal
			cb.sqsum -= oldVal * oldVal
			if math.IsInf(oldVal, 1) || oldVal >= cb.max {
				// need to recalc max
				max := math.Inf(-1)
				for i := cb.head + 1 - cb.size; i < cb.head; i++ {
					fval := float64(cb.buffer[i % cb.size])
					if math.IsNaN(fval) {
						continue
					} else if math.IsInf(fval, 1) {
						max = fval
						break
					} else if fval > max {
						max = fval
					}
				}
				cb.max = max
			} else if math.IsInf(oldVal, -1) || oldVal <= cb.min {
				// need to recalc min
				min := math.Inf(1)
				for i := cb.head + 1 - cb.size; i < cb.head; i++ {
					fval := float64(cb.buffer[i % cb.size])
					if math.IsNaN(fval) {
						continue
					} else if math.IsInf(fval, -1) {
						min = fval
						break
					} else if fval < min {
						min = fval
					}
				}
				cb.min = min
			}
		}
	}
	cb.buffer[idx] = val
	fval := float64(val)
	if !math.IsNaN(fval) {
		cb.count += 1
		cb.sum += fval
		cb.sqsum += fval * fval
		if fval > cb.max {
			cb.max = fval
		}
		if fval < cb.min {
			cb.min = fval
		}
	}
	cb.head += 1
}

func (cb *NumericCircularBuffer[T]) Get(idx int) (val T, err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	if idx < 0 || idx >= cb.head || idx < cb.head - cb.size {
		return val, ErrIndexOutOfRange
	}
	return cb.buffer[idx % cb.size], nil
}

func (cb *NumericCircularBuffer[T]) Slice() []T {
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

func (cb *NumericCircularBuffer[T]) Min() float64 {
	return cb.min
}

func (cb *NumericCircularBuffer[T]) Max() float64 {
	return cb.max
}

func (cb *NumericCircularBuffer[T]) Sum() float64 {
	return cb.sum
}

func (cb *NumericCircularBuffer[T]) Mean() float64 {
	return cb.sum / cb.count
}

func (cb *NumericCircularBuffer[T]) RMS() float64 {
	return math.Sqrt(cb.sqsum) / cb.count
}
