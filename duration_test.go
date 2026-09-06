package utils

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func Test_AppendDuration(t *testing.T) {
	t.Parallel()
	check := func(d time.Duration) {
		t.Helper()
		want := d.String()
		if got := string(AppendDuration(nil, d)); got != want {
			t.Fatalf("AppendDuration(%d) = %q, want %q", int64(d), got, want)
		}
		if got := string(AppendDuration([]byte("x"), d)); got != "x"+want {
			t.Fatalf("AppendDuration(%d) with prefix = %q, want %q", int64(d), got, "x"+want)
		}
	}

	// Every unit boundary, the values just around it, and the extremes.
	for _, base := range []time.Duration{
		0, 1, 999, 1000, 1001, 999999, time.Millisecond, time.Millisecond + 1,
		999999999, time.Second, time.Second + 1, 1500 * time.Millisecond,
		59 * time.Second, time.Minute, time.Minute + time.Second, 59*time.Minute + 59*time.Second,
		time.Hour, time.Hour + time.Nanosecond, 100 * time.Hour, 2540400 * time.Hour,
		math.MaxInt64, math.MinInt64, math.MinInt64 + 1,
	} {
		check(base)
		check(-base)
		check(base + 1)
		check(base - 1)
	}
	// Every count of trailing zeros in the fraction, in every unit.
	for _, unit := range []time.Duration{time.Microsecond, time.Millisecond, time.Second} {
		for f := time.Duration(1); f < unit; f *= 10 {
			check(unit + f)
			check(unit + 7*f)
			check(unit - f)
		}
	}
	// Dense sweeps through the small units.
	for d := range 3000 {
		check(time.Duration(d))
	}
	for d := time.Duration(0); d < 3*time.Second; d += 1234567 {
		check(d)
		check(-d)
	}
	// A deterministic pseudo-random spread across every magnitude.
	rng := rand.New(rand.NewSource(1)) //nolint:gosec // deterministic test data
	for range 200000 {
		x := rng.Uint64()
		check(time.Duration(x >> (x % 64)))
		check(-time.Duration(x >> (x % 64)))
	}
}

func Benchmark_AppendDuration(b *testing.B) {
	inputs := []struct {
		name  string
		value time.Duration
	}{
		{"micros", 823 * time.Microsecond},
		{"millis-frac", 12345678 * time.Nanosecond},
		{"hours", 3*time.Hour + 4*time.Minute + 5*time.Second},
	}
	dst := make([]byte, 0, 32)
	for _, input := range inputs {
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				AppendDuration(dst, input.value)
			}
		})
		b.Run(input.name+"/default", func(b *testing.B) {
			b.ReportAllocs()
			var s string
			for b.Loop() {
				s = input.value.String()
			}
			_ = s
		})
		b.Run(input.name+"/fmt", func(b *testing.B) {
			b.ReportAllocs()
			var out []byte
			for b.Loop() {
				out = fmt.Appendf(dst, "%v", input.value)
			}
			_ = out
		})
	}
}
