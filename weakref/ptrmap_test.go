package weakref

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"runtime"
	"testing"
	"time"
	"unsafe"
)

func TestWeakPtrMapTiny(t *testing.T) {
	type Obj struct{}

	// Allocate struct with pointer to avoid hitting tinyalloc.
	// Otherwise we can't be sure when the allocation will be freed.
	type T struct {
		obj Obj
		p   unsafe.Pointer
	}

	wpm := NewWeakPtrMap()

	obj := &new(T).obj
	ptr := reflect.ValueOf(obj).Pointer()

	wpm.Put(obj, 1)
	val, ok := wpm.Get(obj)
	assert.True(t, ok)
	assert.Equal(t, val.(int), 1)

	obj = nil
	runtime.GC()
	time.Sleep(time.Millisecond)

	assert.Empty(t, wpm.ptrMap)
	val, ok = wpm.Get((*Obj)(unsafe.Pointer(ptr)))
	assert.False(t, ok)
	assert.Nil(t, val)

	wpm.Close()
	assert.Empty(t, wpmRegistry)
}

func TestWeakPtrMapBig(t *testing.T) {
	type Obj struct {
		fill uint64
		it   bool
		up   string
	}

	wpm := NewWeakPtrMap()

	objs := make(map[uintptr]*Obj)
	vals := make(map[uintptr]int)
	for i := 0; i < 100; i++ {
		obj := &Obj{0xDEADBEEFDEADBEEF, true, "It matters not how strait the gate"}
		ptr := reflect.ValueOf(obj).Pointer()
		objs[ptr] = obj
		vals[ptr] = i

		wpm.Put(obj, i)
		val, ok := wpm.Get(obj)
		assert.True(t, ok)
		assert.Equal(t, val.(int), i)
	}

	time.Sleep(time.Millisecond)
	for ptr, obj := range objs {
		val, ok := wpm.Get(obj)
		assert.True(t, ok)
		assert.Equal(t, val.(int), vals[ptr])
	}
	for ptr := range vals {
		delete(objs, ptr)
	}

	runtime.GC()
	time.Sleep(time.Millisecond)

	assert.Empty(t, wpm.ptrMap)
	for ptr := range vals {
		val, ok := wpm.Get((*Obj)(unsafe.Pointer(ptr)))
		assert.False(t, ok)
		assert.Nil(t, val)
	}

	wpm.Close()
	assert.Empty(t, wpmRegistry)
}
