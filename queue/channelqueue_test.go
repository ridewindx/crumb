package queue

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestChannelQueuePush(t *testing.T) {
	q := NewChannelQueue(10)

	q.Push(`test`)
	assert.Equal(t, 1, q.Size())

	result, err := q.Pop()
	assert.Nil(t, err)
	assert.Equal(t, `test`, result)
	assert.True(t, q.Empty())

	q.Push(`test2`)
	assert.Equal(t, 1, q.Size())

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
	assert.Equal(t, 0, q.Size())

	q.Push(`1`)
	q.Push(`2`)

	result, err = q.Pop()
	assert.Nil(t, err)

	assert.Equal(t, `1`, result)
	assert.Equal(t, 1, q.Size())

	result, err = q.Pop()
	assert.Nil(t, err)

	assert.Equal(t, `2`, result)
}

