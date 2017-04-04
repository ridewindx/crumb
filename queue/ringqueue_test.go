package queue

import (
	"testing"
	//"github.com/stretchr/testify/assert"
	"time"
	"fmt"
)

func TestRingQueue(t *testing.T) {
	rq := NewRingQueue(8, false)

	rq.Push(1)
	rq.Push(1)
	rq.Push(1)
	rq.Push(1)
	rq.Push(1)
	rq.Push(1)
	rq.Push(1)
	rq.Push(1)

	t.Run("haha", func(t *testing.T) {
		t.Parallel()
		fmt.Println("get sleep")
		time.Sleep(time.Second)
		fmt.Println("get awake")
		rq.Pop()
		fmt.Println("get awake")
	})

	t.Run("hahaput", func(t *testing.T) {
		t.Parallel()
		fmt.Println("put begin")
		rq.Push(1)
		fmt.Println("put end")
	})
	/*v, err := rq.Get()
	rq.Put(1)
	v, err = rq.Get()
	assert.Nil(t, err)
	assert.Equal(t, 1, v)
	rq.Get()*/
}
