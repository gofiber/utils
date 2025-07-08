package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func checkTimeStamp(tb testing.TB, expectedCurrent, actualCurrent uint32) {
	tb.Helper()
	// Test with buffer of ±2 seconds for CI environment tolerance
	require.True(tb, actualCurrent >= expectedCurrent-2 && actualCurrent <= expectedCurrent+2,
		"Expected timestamp %d (±2s), got %d", expectedCurrent, actualCurrent)
}

func Test_TimeStampUpdater(t *testing.T) {
	StartTimeStampUpdater()

	now := uint32(time.Now().Unix())
	t.Logf("Test start wall time: %d, Timestamp(): %d", now, Timestamp())

	// Wait for the timestamp to catch up (max 100ms * 100 = 10s)
	for i := 0; i < 100; i++ {
		if Timestamp() >= now {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Logf("After wait: wall time: %d, Timestamp(): %d", uint32(time.Now().Unix()), Timestamp())

	checkTimeStamp(t, now, Timestamp())

	// one second later
	<-time.After(1 * time.Second)
	t.Logf("After 1s: wall time: %d, Timestamp(): %d", uint32(time.Now().Unix()), Timestamp())
	checkTimeStamp(t, now+1, Timestamp())

	// two seconds later
	<-time.After(1 * time.Second)
	t.Logf("After 2s: wall time: %d, Timestamp(): %d", uint32(time.Now().Unix()), Timestamp())
	checkTimeStamp(t, now+2, Timestamp())
}

func Test_StopTimeStampUpdater(t *testing.T) {
	// Start the timestamp updater
	StartTimeStampUpdater()

	// Stop the updater
	StopTimeStampUpdater()

	// Capture the timestamp after stopping
	stoppedTime := Timestamp()

	// Wait before checking the timestamp
	<-time.After(5 * time.Second)
	// It should not have changed since we've stopped the updater
	require.Equal(t, stoppedTime, Timestamp(), "timestamp should not change after stopping updater")
}

func Benchmark_CalculateTimestamp(b *testing.B) {
	StartTimeStampUpdater()

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
			// Only time the actual function call
			res := Timestamp()

			// Stop timer while validating
			bb.StopTimer()
			checkTimeStamp(bb, uint32(time.Now().Unix()), res)
			bb.StartTimer()
		}
	})

	b.Run("default_asserted", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < bb.N; n++ {
			// Capture expected time before starting timing
			bb.StopTimer()
			expected := uint32(time.Now().Unix())
			bb.StartTimer()

			// Only time the actual function call
			res := uint32(time.Now().Unix())

			// Stop timer while validating
			bb.StopTimer()
			checkTimeStamp(bb, expected, res)
			bb.StartTimer()
		}
	})
}
