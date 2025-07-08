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
	updaterMu   sync.Mutex
)

// Timestamp returns the current cached Unix timestamp (seconds).
// Call StartTimeStampUpdater() once at app startup for best performance.
func Timestamp() uint32 {
	return atomic.LoadUint32(&timestamp)
}

// StartTimeStampUpdater launches a background goroutine that updates the cached timestamp every second.
// It is safe to call multiple times; only the first call will start the updater.
func StartTimeStampUpdater() {
	updaterMu.Lock()
	defer updaterMu.Unlock()
	updaterOnce.Do(func() {
		atomic.StoreUint32(&timestamp, uint32(time.Now().Unix()))
		stop := make(chan struct{})
		done := make(chan struct{})
		stopUpdater = stop
		updaterDone = done
		go func(stopCh, doneCh chan struct{}) {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			defer close(doneCh)
			for {
				select {
				case <-ticker.C:
					atomic.StoreUint32(&timestamp, uint32(time.Now().Unix()))
				case <-stopCh:
					return
				}
			}
		}(stop, done)
	})
}

// StopTimeStampUpdater stops the background updater goroutine.
// Call this on app shutdown to avoid leaking goroutines.
func StopTimeStampUpdater() {
	updaterMu.Lock()
	defer updaterMu.Unlock()
	if stopUpdater != nil {
		close(stopUpdater)
		stopUpdater = nil
	}
	if updaterDone != nil {
		<-updaterDone
		updaterDone = nil
	}
	// Reset the sync.Once so StartTimeStampUpdater can be called again
	updaterOnce = sync.Once{}
}
