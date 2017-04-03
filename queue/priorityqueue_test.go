package queue

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"time"
)

func TestPriorityQueuePushPop(t *testing.T) {
	pq := NewPriorityQueue(func(a, b interface{}) bool {
		return a.(int) > b.(int)
	})

	assert.True(t, pq.Empty())
	assert.Zero(t, pq.Size())

	pq.Push(4)
	pq.Push(2)
	pq.Push(5)
	pq.Push(1)
	pq.Push(3)

	assert.Equal(t, 5, pq.Size())
	assert.EqualValues(t, 5, pq.Pop())

	pq.Push(3)

	assert.EqualValues(t, 4, pq.Pop())
	assert.EqualValues(t, 3, pq.Pop())
	assert.EqualValues(t, 3, pq.Pop())

	pq.Push(5)

	assert.EqualValues(t, 5, pq.Pop())

	assert.EqualValues(t, 2, pq.Pop())
	assert.EqualValues(t, 1, pq.Pop())

	assert.Zero(t, pq.Size())
	assert.Nil(t, pq.Pop())
}

func BenchmarkPriorityQueuePushIncremental(b *testing.B) {
	numItems := 1000000

	pqs := make([]*PriorityQueue, 0, b.N)

	for i := 0; i < b.N; i++ {
		pq := NewPriorityQueue(func(a, b interface{}) bool {
			return a.(int) > b.(int)
		})
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

func BenchmarkPriorityQueuePopIncremental(b *testing.B) {
	numItems := 1000000

	pqs := make([]*PriorityQueue, 0, b.N)

	for i := 0; i < b.N; i++ {
		pq := NewPriorityQueue(func(a, b interface{}) bool {
			return a.(int) > b.(int)
		})
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

func BenchmarkPriorityQueuePushRandom(b *testing.B) {
	numItems := 1000000

	pqs := make([]*PriorityQueue, 0, b.N)

	for i := 0; i < b.N; i++ {
		pq := NewPriorityQueue(func(a, b interface{}) bool {
			return a.(int) > b.(int)
		})
		pqs = append(pqs, pq)
	}

	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq := pqs[i]
		for j := 0; j < numItems; j++ {
			pq.Push(r.Intn(numItems))
		}
	}
}

func BenchmarkPriorityQueuePopRandom(b *testing.B) {
	numItems := 1000000

	pqs := make([]*PriorityQueue, 0, b.N)

	for i := 0; i < b.N; i++ {
		pq := NewPriorityQueue(func(a, b interface{}) bool {
			return a.(int) > b.(int)
		})
		pqs = append(pqs, pq)
	}

	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	for i := 0; i < b.N; i++ {
		pq := pqs[i]
		for j := 0; j < numItems; j++ {
			pq.Push(r.Intn(numItems))
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
