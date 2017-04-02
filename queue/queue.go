package queue

import (
	"time"
	"errors"
)

var (
	ErrTimeout = errors.New("queue timeout")
	ErrFull = errors.New("queue full")
	ErrEmpty = errors.New("queue empty")
)

type Queue chan interface{}

func NewQueue(capacity int) Queue {
	return make(chan interface{}, capacity)
}

func (q Queue) Put(item interface{}, timeout ...time.Duration) error {
	if len(timeout) == 0 {
		select {
		case q <- item:
			return nil
		}
	}

	select {
	case q <- item:
		return nil
	case <-time.After(timeout[0]):
		if timeout[0] <= 0 {
			return ErrFull
		} else {
			return ErrTimeout
		}
	}
}

func (q Queue) Get(timeout ...time.Duration) (interface{}, error) {
	if len(timeout) == 0 {
		select {
		case item := <-q:
			return item, nil
		}
	}

	select {
	case item := <-q:
		return item, nil
	case <-time.After(timeout[0]):
		if timeout[0] <= 0 {
			return nil, ErrEmpty
		} else {
			return nil, ErrTimeout
		}
	}
}

func (q Queue) Size() int {
	return len(q)
}

func (q Queue) Empty() bool {
	return len(q) == 0
}
