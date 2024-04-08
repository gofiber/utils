package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func checkTimeStamp(tb testing.TB, expectedCurrent, actualCurrent uint32) {
	tb.Helper()
	// test with some buffer in front and back of the expectedCurrent time -> because of the timing on the work machine
	require.True(tb, true, actualCurrent >= expectedCurrent-1 || actualCurrent <= expectedCurrent+1)
}

func Test_TimeStampUpdater(t *testing.T) {
	t.Parallel()

	StartTimeStampUpdater()

	now := uint32(time.Now().Unix())
	checkTimeStamp(t, now, Timestamp())

	// one second later
	time.Sleep(1 * time.Second)
	checkTimeStamp(t, now+1, Timestamp())

	// two seconds later
	time.Sleep(1 * time.Second)
	checkTimeStamp(t, now+2, Timestamp())
}

func Benchmark_CalculateTimestamp(b *testing.B) {
	var res uint32
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
			res = Timestamp()
			checkTimeStamp(bb, uint32(time.Now().Unix()), res)
		}
	})
	b.Run("default_asserted", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < bb.N; n++ {
			res = uint32(time.Now().Unix())
			checkTimeStamp(bb, uint32(time.Now().Unix()), res)
		}
	})
}
