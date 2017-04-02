package weakref

import (
	"reflect"
	"runtime"
	"sync"
)

type WeakPtrMap struct {
	ptrMap map[uintptr]interface{}
	mutex  *sync.RWMutex
	closed bool
}

var wpmRegistry = make(map[*WeakPtrMap]struct{})
var wpmMutex = &sync.Mutex{}

func NewWeakPtrMap() *WeakPtrMap {
	wpm := &WeakPtrMap{
		ptrMap: make(map[uintptr]interface{}),
		mutex:  &sync.RWMutex{},
	}

	wpmMutex.Lock()
	defer wpmMutex.Unlock()
	wpmRegistry[wpm] = struct{}{}

	return wpm
}

func (wpm *WeakPtrMap) Close() {
	wpmMutex.Lock()
	defer wpmMutex.Unlock()
	delete(wpmRegistry, wpm)

	wpm.mutex.RLock()
	defer wpm.mutex.RUnlock()
	wpm.closed = true
	wpm.ptrMap = nil
}

func (wpm *WeakPtrMap) Put(ptr interface{}, val interface{}) {
	wpm.mutex.Lock()
	defer wpm.mutex.Unlock()
	if wpm.closed {
		panic("WeakPtrMap has been closed")
	}
	wpm.ptrMap[reflect.ValueOf(ptr).Pointer()] = val

	runtime.SetFinalizer(ptr, func(ptr interface{}) {
		wpmMutex.Lock()
		defer wpmMutex.Unlock()
		if _, ok := wpmRegistry[wpm]; ok {
			delete(wpm.ptrMap, reflect.ValueOf(ptr).Pointer())
		}
	})
}

func (wpm *WeakPtrMap) Get(ptr interface{}) (val interface{}, ok bool) {
	wpm.mutex.RLock()
	defer wpm.mutex.RUnlock()
	if wpm.closed {
		panic("WeakPtrMap has been closed")
	}
	val, ok = wpm.ptrMap[reflect.ValueOf(ptr).Pointer()]
	return
}
