package queue

import "container/heap"

type PriorityQueue struct {
	*pqueue
}

// Function comparePriority reports whether the element a has higher priority than the element b.
func NewPriorityQueue(comparePriority func(a, b interface{}) bool, sizeHint ...int) *PriorityQueue {
	pq := &pqueue{
		comparePriority: comparePriority,
	}
	if len(sizeHint) > 0 {
		pq.items = make([]interface{}, 0, sizeHint[0])
	}
	return &PriorityQueue{
		pqueue: pq,
	}
}

func (pq *PriorityQueue) Push(item interface{}) {
	heap.Push(pq.pqueue, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	if len(pq.items) == 0 {
		return nil
	}
	return heap.Pop(pq.pqueue)
}


func (pq *PriorityQueue) Len() int {
	return len(pq.items)
}

func (pq *PriorityQueue) Empty() bool {
	return len(pq.items) == 0
}

type pqueue struct {
	items           []interface{}
	comparePriority func(a, b interface{}) bool
}

func (pq *pqueue) Len() int {
	return len(pq.items)
}

func (pq *pqueue) Less(i, j int) bool {
	return pq.comparePriority(pq.items[i], pq.items[j])
}

func (pq *pqueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

func (pq *pqueue) Push(item interface{}) {
	pq.items = append(pq.items, item)
}

func (pq *pqueue) Pop() interface{} {
	if len(pq.items) == 0 {
		return nil
	}
	item := pq.items[len(pq.items)-1]
	pq.items = pq.items[0:len(pq.items)-1]
	return item
}
