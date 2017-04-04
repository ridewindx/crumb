package queue

import (
	"runtime"
	"sync/atomic"
	"fmt"
	"time"
)

// RingQueue is a bounded MPMC ring buffer queue that achieves concurrency
// with CAS operations only.
type RingQueue struct {
	padding0       [8]uint64
	tail           uint64
	padding1       [8]uint64
	head           uint64
	padding2       [8]uint64
	mask uint64
	padding3       [8]uint64
	nodes          []*ringnode
	spin bool
	notFull        chan struct{}
	notEmpty       chan struct{}
}

type ringnode struct {
	position uint64
	data     interface{}
}

// NewRingQueue will allocate a RingQueue with the specified capacity.
// The `spin` specifies the waiting strategy when operation will be blocked:
// if true, spin on CAS; otherwise, wait on channel for notification.
func NewRingQueue(capacity uint64, spin bool) *RingQueue {
	if capacity == 0 {
		panic("RingQueue capacity must be greater than 0")
	}
	capacity = roundUp(capacity)

	rb := &RingQueue{
		nodes: make([]*ringnode, capacity),
		mask: capacity-1,
		spin: spin,
		notFull: make(chan struct{}, 1),
		notEmpty: make(chan struct{}, 1),
	}
	for i := uint64(0); i < capacity; i++ {
		rb.nodes[i] = &ringnode{position: i}
	}
	return rb
}

// Push adds the item to the queue. If the queue is full, will block
// until an item is added to the queue. If a nonzero timeout is specified,
// block no more than the timeout duration and return ErrTimeout. If timeout
// is zero, immediately return ErrFull.
func (rb *RingQueue) Push(item interface{}, timeout ...time.Duration) error {
	var n *ringnode
	var pos uint64
	i := 0
	for {
		pos = atomic.LoadUint64(&rb.tail)
		n = rb.nodes[pos&rb.mask]
		seq := atomic.LoadUint64(&n.position)
		switch dif := int64(seq) - int64(pos); {
		case dif == 0:
			if atomic.CompareAndSwapUint64(&rb.tail, pos, pos+1) {
				n.data = item
				atomic.StoreUint64(&n.position, pos+1)
				rb.print()
				if !rb.spin {
					select {
					case rb.notEmpty <- struct{}{}:
					default:
					}
				}
				return nil
			}
		case dif < 0: // queue is full
			if len(timeout) == 0 {
				if !rb.spin {
					<-rb.notFull // wait for a pop
				}
			} else if timeout[0] > 0 {
				if !rb.spin {
					select { // wait for a pop, until timeout
					case <-rb.notFull:
					case <-time.After(timeout[0]):
						return ErrTimeout
					}
				}
			} else {
				return ErrFull
			}
		default: // another push occurred
		}

		if i == 10000 {
			runtime.Gosched() // free up the cpu before the next iteration
			i = 0
		} else {
			i++
		}
	}
}

// Push will return the next item in the queue. If the queue is empty,
// block until an item can be returned. If a nonzero timeout is specified,
// block no more than the timeout duration and return ErrTimeout. If timeout
// is zero, immediately return ErrEmpty.
func (rb *RingQueue) Pop(timeout ...time.Duration) (interface{}, error) {
	var n *ringnode
	pos := atomic.LoadUint64(&rb.head)
	i := 0
	for {
		pos = atomic.LoadUint64(&rb.head)
		n = rb.nodes[pos&rb.mask]
		seq := atomic.LoadUint64(&n.position)
		switch dif := int64(seq) - (int64(pos) + 1); {
		case dif == 0:
			if atomic.CompareAndSwapUint64(&rb.head, pos, pos+1) {
				data := n.data
				n.data = nil
				atomic.StoreUint64(&n.position, pos+rb.mask+1)
				rb.print()
				select {
				case rb.notFull <- struct{}{}:
				default:
				}
				return data, nil
			}
		case dif < 0: // queue is empty
			if len(timeout) == 0 {
				if !rb.spin {
					<-rb.notEmpty // wait for a push
				}
			} else if timeout[0] > 0 {
				if !rb.spin {
					select { // wait for a push, until timeout
					case <-rb.notEmpty:
					case <-time.After(timeout[0]):
						return nil, ErrTimeout
					}
				}
			} else {
				return nil, ErrEmpty
			}
		default: // another pop occurred
		}

		if i == 10000 {
			runtime.Gosched() // free up the cpu before the next iteration
			i = 0
		} else {
			i++
		}
	}
}

// Len returns the number of items in the queue.
func (rb *RingQueue) Len() uint64 {
	return atomic.LoadUint64(&rb.tail) - atomic.LoadUint64(&rb.head)
}

// Empty returns whether the queue is empty.
func (rq *RingQueue) Empty() bool {
	return atomic.LoadUint64(&rq.tail) == atomic.LoadUint64(&rq.head)
}

// roundUp rounds the uint64 v (v > 0) up to the next
// power of 2.
func roundUp(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}

func (rq *RingQueue) print() {
	fmt.Printf("queue %d, dequeue %d\n", rq.tail, rq.head)
	fmt.Print("nodes: ")
	for i := uint64(0); i <= rq.mask; i++ {
		fmt.Printf("%d(%v) ", rq.nodes[i].position, rq.nodes[i].data)
	}
	fmt.Println()
}
