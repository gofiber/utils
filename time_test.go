package utils

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var timerTestMu sync.Mutex

// testEpoch is an arbitrary fixed instant. Nothing depends on its value, only
// on the differences between it and the times derived from it.
var testEpoch = time.Unix(1_700_000_000, 0)

// fakeClock stands in for the package clock so the updater can be stepped
// instead of waited on.
type fakeClock struct {
	mu   sync.Mutex
	t    time.Time
	tick chan time.Time

	creates atomic.Int32
	stops   atomic.Int32
}

// withFakeClock locks out the other timestamp tests, installs a fake clock
// reading start, and returns it with the teardown the caller must defer.
// Teardown order matters: the updater has to be stopped and the real clock put
// back before the next test takes the lock.
func withFakeClock(tb testing.TB, start time.Time) (*fakeClock, func()) {
	tb.Helper()
	timerTestMu.Lock()

	fc := &fakeClock{t: start, tick: make(chan time.Time)}
	realNow, realNewTicker := now, newTicker
	now, newTicker = fc.now, fc.newTicker

	return fc, func() {
		StopTimeStampUpdater()
		now, newTicker = realNow, realNewTicker
		timerTestMu.Unlock()
	}
}

func (fc *fakeClock) now() time.Time {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	return fc.t
}

func (fc *fakeClock) set(t time.Time) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.t = t
}

func (fc *fakeClock) newTicker(time.Duration) (<-chan time.Time, func()) {
	fc.creates.Add(1)
	return fc.tick, func() { fc.stops.Add(1) }
}

// advance moves the clock to t and ticks the updater, returning once the new
// value is published. The second send is the synchronization: the updater
// handles ticks one at a time, so it can only take that one after the store
// that followed the first. It stores the same value again, so it changes
// nothing.
func (fc *fakeClock) advance(t time.Time) {
	fc.set(t)
	fc.tick <- t
	fc.tick <- t
}

func checkTimeStamp(tb testing.TB, expectedCurrent, actualCurrent uint32) {
	tb.Helper()
	// Test with buffer of ±2 seconds for CI environment tolerance. The
	// failure message is built only on failure: boxing require.True's
	// variadic args allocated on every call, which made the fiber_asserted
	// benchmark's B/op metric nondeterministic.
	if actualCurrent < expectedCurrent-2 || actualCurrent > expectedCurrent+2 {
		tb.Fatalf("Expected timestamp %d (±2s), got %d (diff: %d)",
			expectedCurrent, actualCurrent, int64(actualCurrent)-int64(expectedCurrent))
	}
}

// Test_NewTicker exercises the production tick source. Every other test here
// replaces it, so without this one nothing proves that the real one delivers a
// live channel or that the stop it hands back works. A millisecond period keeps
// the only real-time wait in this file down to the cost of one tick; the
// generous timeout is there to fail loudly rather than hang a loaded CI box.
func Test_NewTicker(t *testing.T) {
	timerTestMu.Lock()
	defer timerTestMu.Unlock()

	tick, stop := newTicker(time.Millisecond)
	defer stop()

	select {
	case <-tick:
	case <-time.After(5 * time.Second):
		t.Fatal("the real tick source never fired")
	}
}

func Test_TimeStampUpdater(t *testing.T) {
	fc, done := withFakeClock(t, testEpoch)
	defer done()

	StartTimeStampUpdater()

	base := uint32(testEpoch.Unix())
	require.Equal(t, base, Timestamp(), "Start should publish the current time before returning")

	fc.advance(testEpoch.Add(time.Second))
	require.Equal(t, base+1, Timestamp(), "a tick should publish the current time")

	fc.advance(testEpoch.Add(90 * time.Second))
	require.Equal(t, base+90, Timestamp(), "each tick republishes, it does not increment")
}

func Test_StartTimeStampUpdater_Idempotent(t *testing.T) {
	fc, done := withFakeClock(t, testEpoch)
	defer done()

	StartTimeStampUpdater()
	StartTimeStampUpdater()
	require.Equal(t, int32(1), fc.creates.Load(), "a second Start must not launch a second updater")

	fc.advance(testEpoch.Add(time.Second))
	require.Equal(t, uint32(testEpoch.Unix())+1, Timestamp())
}

func Test_StopTimeStampUpdater(t *testing.T) {
	fc, done := withFakeClock(t, testEpoch)
	defer done()

	StartTimeStampUpdater()
	StopTimeStampUpdater()
	require.Equal(t, int32(1), fc.stops.Load(), "Stop should release the tick source before returning")

	StopTimeStampUpdater()
	require.Equal(t, int32(1), fc.stops.Load(), "a second Stop must be a no-op")

	// The updater is gone, so the clock can move without the cache following.
	// No tick is sent because there is nothing left to receive one.
	fc.set(testEpoch.Add(time.Hour))
	require.Equal(t, uint32(testEpoch.Unix()), Timestamp(), "the timestamp must not move after Stop")
}

func Test_StartStopRepeat(t *testing.T) {
	fc, done := withFakeClock(t, testEpoch)
	defer done()

	StartTimeStampUpdater()
	StopTimeStampUpdater()

	fc.set(testEpoch.Add(5 * time.Second))
	StartTimeStampUpdater()

	require.Equal(t, uint32(testEpoch.Unix())+5, Timestamp(), "a restart republishes the current time")
	require.Equal(t, int32(2), fc.creates.Load(), "the restart needs its own tick source")
	require.Equal(t, int32(1), fc.stops.Load(), "only the first run's tick source is released so far")

	fc.advance(testEpoch.Add(6 * time.Second))
	require.Equal(t, uint32(testEpoch.Unix())+6, Timestamp(), "the restarted updater still ticks")
}

// Benchmark_CalculateTimestamp runs on the real clock: it measures Timestamp
// against time.Now, and the ±2s tolerance below is the cache's own lag.
func Benchmark_CalculateTimestamp(b *testing.B) {
	timerTestMu.Lock()
	defer timerTestMu.Unlock()
	StartTimeStampUpdater()
	defer StopTimeStampUpdater()

	b.Run("fiber", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < bb.N; n++ {
			_ = Timestamp()
		}
	})

	b.Run("default", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < bb.N; n++ {
			_ = uint32(time.Now().Unix())
		}
	})

	// Simplified asserted benchmarks - measure the realistic cost including validation
	b.Run("fiber_asserted", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < bb.N; n++ {
			res := Timestamp()
			expected := uint32(time.Now().Unix())
			checkTimeStamp(bb, expected, res)
		}
	})

	b.Run("default_asserted", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < bb.N; n++ {
			expected := uint32(time.Now().Unix())
			res := uint32(time.Now().Unix())
			checkTimeStamp(bb, expected, res)
		}
	})
}
