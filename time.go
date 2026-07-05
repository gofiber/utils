package utils

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	timestamp   atomic.Uint32
	updaterMu   sync.Mutex
	stopUpdater chan struct{}
	updaterDone chan struct{}
)

// Timestamp returns the current cached Unix timestamp (seconds).
// Call StartTimeStampUpdater() once at app startup for best performance.
func Timestamp() uint32 {
	return timestamp.Load()
}

// StartTimeStampUpdater launches a background goroutine that updates the cached timestamp every second.
// It is safe to call multiple times and from multiple goroutines; only one updater runs at a time.
func StartTimeStampUpdater() {
	updaterMu.Lock()
	defer updaterMu.Unlock()
	if stopUpdater != nil {
		return
	}

	timestamp.Store(uint32(time.Now().Unix()))
	stopUpdater = make(chan struct{})
	updaterDone = make(chan struct{})

	go func(stop, done chan struct{}) {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		defer close(done)

		for {
			select {
			case <-ticker.C:
				timestamp.Store(uint32(time.Now().Unix()))
			case <-stop:
				return
			}
		}
	}(stopUpdater, updaterDone)
}

// StopTimeStampUpdater stops the background updater goroutine.
// Call this on app shutdown to avoid leaking goroutines.
// It is safe to call multiple times and from multiple goroutines.
func StopTimeStampUpdater() {
	updaterMu.Lock()
	defer updaterMu.Unlock()
	if stopUpdater == nil {
		return
	}

	close(stopUpdater)
	<-updaterDone
	stopUpdater = nil
	updaterDone = nil
}
