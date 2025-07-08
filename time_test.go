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
		"Expected timestamp %d (±2s), got %d", expectedCurrent, actualCurrent)
}

func Test_TimeStampUpdater(t *testing.T) {
	timerTestMu.Lock()
	defer timerTestMu.Unlock()
	StartTimeStampUpdater()

	now := uint32(time.Now().Unix())
	t.Logf("Test start wall time: %d, Timestamp(): %d", now, Timestamp())

	// Wait up to 10s for the timestamp to catch up
	timeout := time.After(10 * time.Second)
	tick := time.Tick(100 * time.Millisecond)
waitLoop:
	for {
		select {
		case <-timeout:
			t.Fatalf("Timeout waiting for timestamp to catch up: now=%d, Timestamp()=%d", now, Timestamp())
		case <-tick:
			if Timestamp() >= now {
				break waitLoop
			}
		}
		if Timestamp() >= now {
			break waitLoop
		}
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
	timerTestMu.Lock()
	defer timerTestMu.Unlock()
	StartTimeStampUpdater()

	StopTimeStampUpdater()

	// Wait a short moment to ensure the updater goroutine has exited
	time.Sleep(50 * time.Millisecond)

	stoppedTime := Timestamp()

	<-time.After(5 * time.Second)
	require.Equal(t, stoppedTime, Timestamp(), "timestamp should not change after stopping updater")
}

func Benchmark_CalculateTimestamp(b *testing.B) {
	timerTestMu.Lock()
	defer timerTestMu.Unlock()
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
