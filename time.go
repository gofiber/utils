package utils

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	timestamp   uint32
	updaterOnce sync.Once
	stopUpdater chan struct{}
	updaterDone chan struct{}
)

// Timestamp returns the current cached Unix timestamp (seconds).
// Call StartTimeStampUpdater() once at app startup for best performance.
func Timestamp() uint32 {
	return atomic.LoadUint32(&timestamp)
}

// StartTimeStampUpdater launches a background goroutine that updates the cached timestamp every second.
// It is safe to call multiple times; only the first call will start the updater.
func StartTimeStampUpdater() {
	updaterOnce.Do(func() {
		atomic.StoreUint32(&timestamp, uint32(time.Now().Unix()))
		stopUpdater = make(chan struct{})
		updaterDone = make(chan struct{})

		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			defer close(updaterDone)

			for {
				select {
				case <-ticker.C:
					atomic.StoreUint32(&timestamp, uint32(time.Now().Unix()))
				case <-stopUpdater:
					return
				}
			}
		}()
	})
}

// StopTimeStampUpdater stops the background updater goroutine.
// Call this on app shutdown to avoid leaking goroutines.
func StopTimeStampUpdater() {
	if stopUpdater != nil {
		close(stopUpdater)
		<-updaterDone
		// Reset the sync.Once so StartTimeStampUpdater can be called again
		updaterOnce = sync.Once{}
		stopUpdater = nil
		updaterDone = nil
	}
}
