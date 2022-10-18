package utils

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	timestampTimer sync.Once
	// Timestamp please start the timer function before you use this value
	Timestamp atomic.Uint32
)

// StartTimeStampUpdater starts a concurrent function which stores the timestamp to an atomic value per second,
// which is much better for performance than determining it at runtime each time
func StartTimeStampUpdater() {
	timestampTimer.Do(func() {
		// set initial value
		Timestamp.Store(uint32(time.Now().Unix()))
		go func(sleep time.Duration) {
			ticker := time.NewTicker(sleep)
			defer ticker.Stop()
			for {
				select {
				case t := <-ticker.C:
					// update timestamp
					Timestamp.Store(uint32(t.Unix()))
				}
			}
		}(1 * time.Second) // duration
	})
}
