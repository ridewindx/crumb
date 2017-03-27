package concurrency

import (
	"sync"
)

func All(fns ...func()) (done <-chan struct{}) {
	ch := make(chan struct{}, 1)
	var wg sync.WaitGroup
	wg.Add(len(fns))

	for _, fn := range fns {
		go func(f func()) {
			f()
			wg.Done()
		}(fn)
	}

	go func() {
		wg.Wait()
		ch <- struct {}{}
		close(ch)
	}()

	done = ch
	return
}
