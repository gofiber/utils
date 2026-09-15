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

// now and newTicker are this package's clock: every production read of the wall
// clock goes through them, so replacing the pair is enough to drive the updater
// from a test without waiting on real time. Neither is reachable from outside
// the package, and Timestamp itself touches neither.
//
// newTicker hands back a channel and a stop function instead of a *time.Ticker
// because a Ticker the runtime did not build cannot stand in for one: its Stop
// panics.
var (
	now       = time.Now
	newTicker = func(d time.Duration) (<-chan time.Time, func()) {
		t := time.NewTicker(d)
		return t.C, t.Stop
	}
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

	timestamp.Store(uint32(now().Unix()))
	stopUpdater = make(chan struct{})
	updaterDone = make(chan struct{})

	// The tick source is created here, not in the goroutine, so it exists by
	// the time this returns and a caller cannot race the goroutine's first
	// statement.
	tick, stopTick := newTicker(time.Second)

	go func(stop, done chan struct{}) {
		// stopTick before close(done), so that once StopTimeStampUpdater's
		// wait returns the tick source is released rather than about to be.
		defer close(done)
		defer stopTick()

		for {
			select {
			case <-tick:
				timestamp.Store(uint32(now().Unix()))
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
