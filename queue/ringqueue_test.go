package queue

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"sync"
)

func TestRingQueuePushPop(t *testing.T) {
	rq := NewRingQueue(8, false)

	assert.True(t, rq.Empty())
	assert.Zero(t, rq.Len())

	err := rq.Push(1)
	assert.Nil(t, err)

	assert.False(t, rq.Empty())
	assert.Equal(t, 1, rq.Len())

	val, err := rq.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 1, val)

	rq.Push(1)
	rq.Push(2)

	val, err = rq.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 1, val)

	rq.Push(3)
	rq.Push(4)

	val, err = rq.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 2, val)

	val, err = rq.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 3, val)

	assert.False(t, rq.Empty())
	assert.Equal(t, 1, rq.Len())
}

func TestRingQueueSpinNonblocking(t *testing.T) {
	rq := NewRingQueue(2, true)

	err := rq.Push(1, 0)
	assert.Nil(t, err)
	err = rq.Push(2, 0)
	assert.Nil(t, err)
	err = rq.Push(3, 0)
	assert.Equal(t, ErrFull, err)

	val, err := rq.Pop(0)
	assert.Nil(t, err)
	assert.Equal(t, 1, val)
	val, err = rq.Pop(0)
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
	val, err = rq.Pop(0)
	assert.Equal(t, ErrEmpty, err)
	assert.Nil(t, val)

	err = rq.Push(3, 0)
	assert.Nil(t, err)

	val, err = rq.Pop(0)
	assert.Nil(t, err)
	assert.Equal(t, 3, val)
}

func TestRingQueueSpinTimeout(t *testing.T) {
	rq := NewRingQueue(2, true)

	err := rq.Push(1)
	assert.Nil(t, err)
	err = rq.Push(2)
	assert.Nil(t, err)
	err = rq.Push(3, time.Microsecond)
	assert.Equal(t, ErrTimeout, err)

	val, err := rq.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 1, val)
	val, err = rq.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
	val, err = rq.Pop(time.Microsecond)
	assert.Equal(t, ErrTimeout, err)
	assert.Nil(t, val)

	go func() {
		time.Sleep(2*time.Millisecond)
		err = rq.Push(3)
		assert.Nil(t, err)
	}()

	val, err = rq.Pop(time.Millisecond)
	assert.Equal(t, ErrTimeout, err)
	assert.Nil(t, val)

	val, err = rq.Pop(5*time.Millisecond)
	assert.Nil(t, err)
	assert.Equal(t, 3, val)
}

func TestRingQueueSpinBlocking(t *testing.T) {
	rq := NewRingQueue(2, true)

	err := rq.Push(1)
	assert.Nil(t, err)
	err = rq.Push(2)
	assert.Nil(t, err)

	go func() {
		time.Sleep(2*time.Millisecond)
		val, err := rq.Pop()
		assert.Nil(t, err)
		assert.Equal(t, 1, val)
	}()

	err = rq.Push(3)
	assert.Nil(t, err)

	val, err := rq.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 2, val)

	go func() {
		time.Sleep(2*time.Millisecond)
		err = rq.Push(3)
		assert.Nil(t, err)
	}()

	val, err = rq.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 3, val)
}

func TestRingQueueChannelNonblocking(t *testing.T) {
	rq := NewRingQueue(2, false)

	err := rq.Push(1, 0)
	assert.Nil(t, err)
	err = rq.Push(2, 0)
	assert.Nil(t, err)
	err = rq.Push(3, 0)
	assert.Equal(t, ErrFull, err)

	val, err := rq.Pop(0)
	assert.Nil(t, err)
	assert.Equal(t, 1, val)
	val, err = rq.Pop(0)
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
	val, err = rq.Pop(0)
	assert.Equal(t, ErrEmpty, err)
	assert.Nil(t, val)

	err = rq.Push(3, 0)
	assert.Nil(t, err)

	val, err = rq.Pop(0)
	assert.Nil(t, err)
	assert.Equal(t, 3, val)
}

func TestRingQueueChannelTimeout(t *testing.T) {
	rq := NewRingQueue(2, false)

	err := rq.Push(1)
	assert.Nil(t, err)
	err = rq.Push(2)
	assert.Nil(t, err)
	err = rq.Push(3, time.Microsecond)
	assert.Equal(t, ErrTimeout, err)

	val, err := rq.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 1, val)
	val, err = rq.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
	val, err = rq.Pop(time.Microsecond)
	assert.Equal(t, ErrTimeout, err)
	assert.Nil(t, val)

	go func() {
		time.Sleep(2*time.Millisecond)
		err = rq.Push(3)
		assert.Nil(t, err)
	}()

	val, err = rq.Pop(time.Millisecond)
	assert.Equal(t, ErrTimeout, err)
	assert.Nil(t, val)

	val, err = rq.Pop(5*time.Millisecond)
	assert.Nil(t, err)
	assert.Equal(t, 3, val)
}

func TestRingQueueChannelBlocking(t *testing.T) {
	rq := NewRingQueue(2, false)

	err := rq.Push(1)
	assert.Nil(t, err)
	err = rq.Push(2)
	assert.Nil(t, err)

	go func() {
		time.Sleep(2*time.Millisecond)
		val, err := rq.Pop()
		assert.Nil(t, err)
		assert.Equal(t, 1, val)
	}()

	err = rq.Push(3)
	assert.Nil(t, err)

	val, err := rq.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 2, val)

	go func() {
		time.Sleep(2*time.Millisecond)
		err = rq.Push(3)
		assert.Nil(t, err)
	}()

	val, err = rq.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 3, val)
}

func BenchmarkRingQueueSpinPushPop(b *testing.B) {
	rq := NewRingQueue(64, true)

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

func BenchmarkRingQueueSpinPush(b *testing.B) {
	rq := NewRingQueue(b.N, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rq.Push(i)
	}
}

func BenchmarkRingQueueSpinPop(b *testing.B) {
	rq := NewRingQueue(b.N, true)

	for i := 0; i < b.N; i++ {
		err := rq.Push(i)
		assert.Nil(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rq.Pop()
	}
}

func BenchmarkRingQueueChannelPushPop(b *testing.B) {
	rq := NewRingQueue(64, false)

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

func BenchmarkRingQueueChannelPush(b *testing.B) {
	rq := NewRingQueue(b.N, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rq.Push(i)
	}
}

func BenchmarkRingQueueChannelPop(b *testing.B) {
	rq := NewRingQueue(b.N, false)

	for i := 0; i < b.N; i++ {
		err := rq.Push(i)
		assert.Nil(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rq.Pop()
	}
}

func BenchmarkRingQueueSpinParallelPushPop(b *testing.B) {
	rq := NewRingQueue(b.N, true)

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

func BenchmarkRingQueueSpinParallelPush(b *testing.B) {
	rq := NewRingQueue(b.N, true)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			rq.Push(i)
			i++
		}
	})
}

func BenchmarkRingQueueSpinParallelPop(b *testing.B) {
	rq := NewRingQueue(b.N, true)

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
