package concurrency

import (
	"testing"
	"time"
)

func TestAll(test *testing.T) {
	start := time.Now()
	var val1, val2, val3 bool
	done := All(
		func() {
			val1 = true
			time.Sleep(100 * time.Millisecond)
		},
		func() {
			val2 = true
			time.Sleep(100 * time.Millisecond)
		},
		func() {
			val3 = true
			time.Sleep(100 * time.Millisecond)
		},
	)
	<-done
	diff := time.Now().Sub(start)
	if diff > 105*time.Millisecond {
		test.Errorf("All takes too long to complete")
	}
	if !(val1 && val2 && val3) {
		test.Errorf("Expected all to run, but at least one didn't")
	}
}
