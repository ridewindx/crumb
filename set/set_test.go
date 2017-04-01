package set

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	s := NewSet()

	s.Add(1)
	s.Add(2)
	assert.EqualValues(t, map[interface{}]struct{}{
		1: struct{}{},
		2: struct{}{},
	}, s)

	s.Remove(2)
	assert.EqualValues(t, map[interface{}]struct{}{
		1: struct{}{},
	}, s)

	assert.True(t, s.Contains(1))
	assert.False(t, s.Contains(2))
	assert.Equal(t, s.Size(), 1)

	s.Clear()
	assert.Empty(t, s)
	assert.Equal(t, s.Size(), 0)

	s = NewSet()

	s.Add(1)
	s.Add(2)
	ss := NewSet()
	for _, i := range s.ToSlice() {
		ss.Add(i)
	}
	assert.EqualValues(t, s, ss)
	assert.True(t, s.Equal(ss))
	assert.True(t, s.IsSubset(ss))
	assert.True(t, s.IsSuperset(ss))
	assert.EqualValues(t, s, s.Union(ss))
	assert.EqualValues(t, s, s.Intersect(ss))
	assert.Empty(t, s.Difference(ss))
	assert.Empty(t, s.SymmetricDifference(ss))

	cloned := s.Clone()
	assert.EqualValues(t, s, cloned)
	s.Add(3)
	assert.False(t, cloned.Contains(3))

	assert.False(t, s.Equal(ss))
	assert.False(t, s.IsSubset(ss))
	assert.True(t, s.IsSuperset(ss))
	assert.EqualValues(t, s, s.Union(ss))
	assert.EqualValues(t, ss, s.Intersect(ss))
	assert.EqualValues(t, map[interface{}]struct{}{
		3: struct{}{},
	}, s.Difference(ss))
	assert.EqualValues(t, map[interface{}]struct{}{
		3: struct{}{},
	}, s.SymmetricDifference(ss))
}
