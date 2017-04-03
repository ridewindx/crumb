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

type ChannelQueue chan interface{}

func NewChannelQueue(capacity int) ChannelQueue {
	return make(chan interface{}, capacity)
}

func (cq ChannelQueue) Push(item interface{}, timeout ...time.Duration) error {
	if len(timeout) == 0 {
		select {
		case cq <- item:
			return nil
		}
	}

	select {
	case cq <- item:
		return nil
	case <-time.After(timeout[0]):
		if timeout[0] <= 0 {
			return ErrFull
		} else {
			return ErrTimeout
		}
	}
}

func (cq ChannelQueue) Pop(timeout ...time.Duration) (interface{}, error) {
	if len(timeout) == 0 {
		select {
		case item := <-cq:
			return item, nil
		}
	}

	select {
	case item := <-cq:
		return item, nil
	case <-time.After(timeout[0]):
		if timeout[0] <= 0 {
			return nil, ErrEmpty
		} else {
			return nil, ErrTimeout
		}
	}
}

func (cq ChannelQueue) Len() int {
	return len(cq)
}

func (cq ChannelQueue) Empty() bool {
	return len(cq) == 0
}
