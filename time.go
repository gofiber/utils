package utils

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	timestampTimer sync.Once
	timestamp      uint32
	stopChan       chan struct{}
)

// Timestamp returns the current time.
// Make sure to start the updater once using StartTimeStampUpdater() before calling.
func Timestamp() uint32 {
	return atomic.LoadUint32(&timestamp)
}

// StartTimeStampUpdater starts a concurrent function which stores the timestamp to an atomic value per second,
// which is much better for performance than determining it at runtime each time
func StartTimeStampUpdater() {
	timestampTimer.Do(func() {
		atomic.StoreUint32(&timestamp, uint32(time.Now().Unix()))

		c := make(chan struct{})
		stopChan = c

		go func(localChan chan struct{}, sleep time.Duration) {
			ticker := time.NewTicker(sleep)
			defer ticker.Stop()

			for {
				select {
				case t := <-ticker.C:
					atomic.StoreUint32(&timestamp, uint32(t.Unix()))
				case <-localChan:
					return
				}
			}
		}(c, 1*time.Second)
	})
}

// StopTimeStampUpdater stops the timestamp updater
// WARNING: Make sure to call this function before the program exits, otherwise it will leak goroutines
func StopTimeStampUpdater() {
	if stopChan != nil {
		close(stopChan)
		stopChan = nil
	}
}
