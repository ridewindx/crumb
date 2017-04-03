package queue

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestQueuePushPop(t *testing.T) {
	pq := NewQueue()

	assert.True(t, pq.Empty())
	assert.Zero(t, pq.Size())

	pq.Push(1)
	pq.Push(2)
	pq.Push(3)
	pq.Push(4)
	pq.Push(5)

	assert.Equal(t, 5, pq.Size())
	assert.EqualValues(t, 1, pq.Pop())

	pq.Push(6)

	assert.EqualValues(t, 2, pq.Pop())
	assert.EqualValues(t, 3, pq.Pop())
	assert.EqualValues(t, 4, pq.Pop())

	pq.Push(7)

	assert.EqualValues(t, 5, pq.Pop())
	assert.EqualValues(t, 6, pq.Pop())
	assert.EqualValues(t, 7, pq.Pop())

	assert.Zero(t, pq.Size())
	assert.Nil(t, pq.Pop())
}

func BenchmarkQueuePushIncremental(b *testing.B) {
	numItems := 1000000

	pqs := make([]*Queue, 0, b.N)

	for i := 0; i < b.N; i++ {
		pq := NewQueue()
		pqs = append(pqs, pq)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq := pqs[i]
		for j := 0; j < numItems; j++ {
			pq.Push(j)
		}
	}
}

func BenchmarkQueuePopIncremental(b *testing.B) {
	numItems := 1000000

	pqs := make([]*Queue, 0, b.N)

	for i := 0; i < b.N; i++ {
		pq := NewQueue()
		pqs = append(pqs, pq)
	}

	for i := 0; i < b.N; i++ {
		pq := pqs[i]
		for j := 0; j < numItems; j++ {
			pq.Push(j)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq := pqs[i]
		for j := 0; j < numItems; j++ {
			pq.Pop()
		}
	}
}
