package utils

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var timerTestMu sync.Mutex

var testEpoch = time.Unix(1_700_000_000, 0)

type fakeClock struct {
	mu   sync.Mutex
	t    time.Time
	tick chan time.Time

	creates atomic.Int32
	stops   atomic.Int32
}

// withFakeClock locks out other timestamp tests while it swaps the package clock; the deferred
// teardown stops the updater and restores the clock before unlocking.
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

// advance ticks twice: the updater takes the second tick only after storing the first, so the
// value is published on return.
func (fc *fakeClock) advance(t time.Time) {
	fc.set(t)
	fc.tick <- t
	fc.tick <- t
}

func checkTimeStamp(tb testing.TB, expectedCurrent, actualCurrent uint32) {
	tb.Helper()
	// Not require.True: its boxed variadic args made fiber_asserted's B/op nondeterministic.
	if actualCurrent < expectedCurrent-2 || actualCurrent > expectedCurrent+2 {
		tb.Fatalf("Expected timestamp %d (±2s), got %d (diff: %d)",
			expectedCurrent, actualCurrent, int64(actualCurrent)-int64(expectedCurrent))
	}
}

// Test_NewTicker covers the production tick source, which every other test replaces.
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

	// Set, not advance: nothing is left to receive a tick.
	fc.set(testEpoch.Add(time.Hour))
	require.Equal(t, uint32(testEpoch.Unix())+3600, Timestamp(), "after Stop the timestamp must follow the clock")
}

// Test_Timestamp_WithoutUpdater: without an updater (before Start, after Stop) Timestamp reads the clock.
func Test_Timestamp_WithoutUpdater(t *testing.T) {
	fc, done := withFakeClock(t, testEpoch)
	defer done()

	timestamp.Store(0)
	defer timestamp.Store(0)

	base := uint32(testEpoch.Unix())
	fc.set(testEpoch.Add(time.Hour))
	require.Equal(t, base+3600, Timestamp(), "without an updater Timestamp reads the clock")

	StartTimeStampUpdater()
	require.Equal(t, base+3600, Timestamp(), "Start publishes before it returns")
	fc.set(testEpoch.Add(2 * time.Hour))
	require.Equal(t, base+3600, Timestamp(), "a running updater serves what it published")

	StopTimeStampUpdater()
	require.Equal(t, base+7200, Timestamp(), "after Stop Timestamp reads the clock again")
}

// Test_Timestamp_Wraps2106: Timestamp wraps to 0 on 2106-02-07, whether the cache or the clock answers.
func Test_Timestamp_Wraps2106(t *testing.T) {
	wrap := time.Unix(1<<32, 0)
	fc, done := withFakeClock(t, wrap.Add(-time.Second))
	defer done()
	timestamp.Store(0)
	defer timestamp.Store(0)

	StartTimeStampUpdater()
	require.Equal(t, uint32(1<<32-1), Timestamp(), "the last second before the wrap")

	fc.advance(wrap)
	require.Zero(t, Timestamp(), "the wrapped second, whether the cache or the clock answers")

	fc.advance(wrap.Add(time.Second))
	require.Equal(t, uint32(1), timestamp.Load(), "the cache is nonzero again")
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

// Benchmark_CalculateTimestamp runs on the real clock; the ±2s tolerance below is the cache's own lag.
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
