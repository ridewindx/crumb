package queue

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	q := NewQueue(10)

	q.Put(`test`)
	assert.Equal(t, 1, q.Size())

	result, err := q.Get()
	assert.Nil(t, err)
	assert.Equal(t, `test`, result)
	assert.True(t, q.Empty())

	q.Put(`test2`)
	assert.Equal(t, 1, q.Size())

	result, err = q.Get()
	assert.Nil(t, err)

	assert.Equal(t, `test2`, result)
	assert.True(t, q.Empty())
}

func TestGet(t *testing.T) {
	q := NewQueue(10)

	q.Put(`test`)
	result, err := q.Get()
	assert.Nil(t, err)

	assert.Equal(t, `test`, result)
	assert.Equal(t, 0, q.Size())

	q.Put(`1`)
	q.Put(`2`)

	result, err = q.Get()
	assert.Nil(t, err)

	assert.Equal(t, `1`, result)
	assert.Equal(t, 1, q.Size())

	result, err = q.Get()
	assert.Nil(t, err)

	assert.Equal(t, `2`, result)
}

