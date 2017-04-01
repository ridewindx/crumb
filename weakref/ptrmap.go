package weakref

import (
	"unsafe"
	"reflect"
	"runtime"
	"sync"
	"fmt"
)

type WeakPtrMap struct {
	ptrMap map[uintptr]interface{}
	mutex *sync.RWMutex
}

var weakPtrMaps = make(map[uintptr]*WeakPtrMap)
var weakPtrMapsMutex = &sync.Mutex{}

func NewWeakPtrMap() *WeakPtrMap {
	wpm := &WeakPtrMap{
		ptrMap: make(map[uintptr]interface{}),
		mutex: &sync.RWMutex{},
	}

	weakPtrMapsMutex.Lock()
	defer weakPtrMapsMutex.Unlock()
	weakPtrMaps[uintptr(unsafe.Pointer(wpm))] = wpm

	runtime.SetFinalizer(wpm, func(wpm interface{}) {
		weakPtrMapsMutex.Lock()
		defer weakPtrMapsMutex.Unlock()
		delete(weakPtrMaps, uintptr(unsafe.Pointer(wpm.(*WeakPtrMap))))
	})

	return wpm
}

func (wpm *WeakPtrMap) Put(ptr interface{}, val interface{}) {
	wpm.mutex.Lock()
	wpm.ptrMap[reflect.ValueOf(ptr).Pointer()] = val
	wpm.mutex.Unlock()

	wpmPtr := uintptr(unsafe.Pointer(wpm))

	runtime.SetFinalizer(ptr, func(ptr interface{}) {
		fmt.Println("haha")
		weakPtrMapsMutex.Lock()
		defer weakPtrMapsMutex.Unlock()
		if wpm, ok := weakPtrMaps[wpmPtr]; ok {
			delete(wpm.ptrMap, reflect.ValueOf(ptr).Pointer())
		}
	})
}

func (wpm *WeakPtrMap) Get(ptr interface{}) (val interface{}, ok bool) {
	wpm.mutex.RLock()
	defer wpm.mutex.RUnlock()
	val, ok = wpm.ptrMap[reflect.ValueOf(ptr).Pointer()]
	return
}
