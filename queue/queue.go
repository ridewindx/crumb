package queue

import (
	"container/list"
)

type Queue struct {
	list *list.List
}

func NewQueue() *Queue {
	return &Queue{
		list: list.New(),
	}
}

func (q *Queue) Push(item interface{}) {
	q.list.PushBack(item)
}

func (q *Queue) Pop() interface{} {
	e := q.list.Front()
	if e == nil {
		return nil
	}

	return q.list.Remove(e)
}

func (q *Queue) Len() int {
	return q.list.Len()
}

func (q *Queue) Empty() bool {
	return q.list.Len() == 0
}
