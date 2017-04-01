package weakref

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"runtime"
	"fmt"
	"time"
)

func TestWeakPtrMap(t *testing.T) {
	type Obj struct {}

	wpm := NewWeakPtrMap()

	done := make(chan bool, 1)
	obj := &Obj{}
	runtime.SetFinalizer(obj, func(obj *Obj) {
		fmt.Println("hehe")
		done <- true
	})

	wpm.Put(obj, 1)
	val, ok := wpm.Get(obj)
	assert.True(t, ok)
	assert.Equal(t, val.(int), 1)
	obj = nil
	runtime.GC()
	select {
	case <-done:
	case <-time.After(time.Second * 4):
		t.Error("finalizer for type WeakPtrMap didn't run")
	}
	assert.Empty(t, wpm.ptrMap)
}
