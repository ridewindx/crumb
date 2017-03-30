package concurrency

import (
	"time"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	w := NewWorker()
	ch := make(chan int)

	for i := 0; i < 3; i++ {
		assert.True(t, w.Stopped())

		w.Start(func(sentry *Sentry) {
			<-ch
			assert.True(t, sentry.Paused())
			ch <- 1
			sentry.Sleep()
			ch <- 2

			<-ch
			ch <- 3
			sentry.Sleep()
			ch <- 4

			<-ch
			time.Sleep(5 * time.Microsecond)
			sentry.Sleep()
			if sentry.Stopped() {
				return
			}
		})

		assert.False(t, w.Stopped())
		assert.False(t, w.Paused())
		assert.Panics(t, func() {
			w.Start(func(*Sentry){})
		})

		w.Pause()
		assert.True(t, w.Paused())
		ch <- 0
		assert.Equal(t, <-ch, 1)

		time.Sleep(5 * time.Microsecond)
		w.Resume()
		assert.False(t, w.Paused())
		assert.Equal(t, <-ch, 2)

		w.Pause()
		assert.True(t, w.Paused())
		ch <- 0
		assert.Equal(t, <-ch, 3)
		w.Resume()
		assert.False(t, w.Paused())
		assert.Equal(t, <-ch, 4)

		w.Pause()
		ch <- 0
		w.Stop()
		assert.True(t, w.Stopped())

		assert.NotPanics(t, w.Stop)
	}
}
