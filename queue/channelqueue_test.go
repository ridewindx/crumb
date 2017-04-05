package queue

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"sync"
)

func TestChannelQueuePush(t *testing.T) {
	q := NewChannelQueue(10)

	q.Push(`test`)
	assert.Equal(t, 1, q.Len())

	result, err := q.Pop()
	assert.Nil(t, err)
	assert.Equal(t, `test`, result)
	assert.True(t, q.Empty())

	q.Push(`test2`)
	assert.Equal(t, 1, q.Len())

	result, err = q.Pop()
	assert.Nil(t, err)

	assert.Equal(t, `test2`, result)
	assert.True(t, q.Empty())
}

func TestChannelQueuePop(t *testing.T) {
	q := NewChannelQueue(10)

	q.Push(`test`)
	result, err := q.Pop()
	assert.Nil(t, err)

	assert.Equal(t, `test`, result)
	assert.Equal(t, 0, q.Len())

	q.Push(`1`)
	q.Push(`2`)

	result, err = q.Pop()
	assert.Nil(t, err)

	assert.Equal(t, `1`, result)
	assert.Equal(t, 1, q.Len())

	result, err = q.Pop()
	assert.Nil(t, err)

	assert.Equal(t, `2`, result)
}

func BenchmarkChannelQueuePushPop(b *testing.B) {
	rq := NewChannelQueue(64)

	var count int
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			_, err := rq.Pop()
			assert.Nil(b, err)

			count++
			if count == b.N {
				return
			}
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rq.Push(i)
	}

	wg.Wait()
}

func BenchmarkChannelQueuePush(b *testing.B) {
	rq := NewChannelQueue(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rq.Push(i)
	}
}

func BenchmarkChannelQueuePop(b *testing.B) {
	rq := NewChannelQueue(b.N)

	for i := 0; i < b.N; i++ {
		err := rq.Push(i)
		assert.Nil(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rq.Pop()
	}
}

func BenchmarkChannelQueueParallelPushPop(b *testing.B) {
	rq := NewChannelQueue(b.N)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			rq.Push(i)
			rq.Pop()
			i++
		}
	})
}

func BenchmarkChannelQueueParallelPush(b *testing.B) {
	rq := NewChannelQueue(b.N)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			rq.Push(i)
			i++
		}
	})
}

func BenchmarkChannelQueueParallelPop(b *testing.B) {
	rq := NewChannelQueue(b.N)

	for i := 0; i < b.N; i++ {
		err := rq.Push(i)
		assert.Nil(b, err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rq.Pop()
		}
	})
}
