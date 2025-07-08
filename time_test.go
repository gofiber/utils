package utils

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var timerTestMu sync.Mutex

func checkTimeStamp(tb testing.TB, expectedCurrent, actualCurrent uint32) {
	tb.Helper()
	// Test with buffer of ±2 seconds for CI environment tolerance
	require.True(tb, actualCurrent >= expectedCurrent-2 && actualCurrent <= expectedCurrent+2,
		"Expected timestamp %d (±2s), got %d (diff: %d)", expectedCurrent, actualCurrent, int64(actualCurrent)-int64(expectedCurrent))
}

func Test_TimeStampUpdater(t *testing.T) {
	timerTestMu.Lock()
	defer timerTestMu.Unlock()

	StartTimeStampUpdater()
	defer StopTimeStampUpdater()

	now := uint32(time.Now().Unix())

	// Give it a moment to initialize
	time.Sleep(100 * time.Millisecond)

	ts := Timestamp()
	require.True(t, ts >= now-1 && ts <= now+1,
		"Initial timestamp should be within ±1 second of current time. Expected: %d±1, got: %d (diff: %d)",
		now, ts, int64(ts)-int64(now))

	// Wait for next update
	time.Sleep(1100 * time.Millisecond)

	ts2 := Timestamp()
	require.Greater(t, ts2, ts,
		"Timestamp should have updated after 1+ seconds. Initial: %d, after 1s: %d",
		ts, ts2)

	currentTime := uint32(time.Now().Unix())
	require.True(t, ts2 >= now && ts2 <= currentTime+1,
		"Updated timestamp should be between test start (%d) and current time (%d), got: %d",
		now, currentTime, ts2)
}

func Test_StopTimeStampUpdater(t *testing.T) {
	timerTestMu.Lock()
	defer timerTestMu.Unlock()

	StartTimeStampUpdater()

	// Get initial timestamp
	time.Sleep(100 * time.Millisecond)
	initial := Timestamp()

	StopTimeStampUpdater()

	// Verify it stops updating
	time.Sleep(2 * time.Second)
	final := Timestamp()
	currentTime := uint32(time.Now().Unix())
	require.Equal(t, initial, final,
		"Timestamp should not change after stopping updater. Stopped at: %d, checked at: %d (wall time: %d)",
		initial, final, currentTime)
}

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
			bb.StopTimer()
			checkTimeStamp(bb, uint32(time.Now().Unix()), res)
			bb.StartTimer()
		}
	})

	b.Run("default_asserted", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < bb.N; n++ {
			bb.StopTimer()
			expected := uint32(time.Now().Unix())
			bb.StartTimer()
			res := uint32(time.Now().Unix())
			bb.StopTimer()
			checkTimeStamp(bb, expected, res)
			bb.StartTimer()
		}
	})
}
