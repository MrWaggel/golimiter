package golimiter

import (
	"testing"
	"time"
)

func BenchmarkLimiterManager_IsLimited(b *testing.B) {
	key := 1
	lim := New(200, time.Minute)
	for i := 0; i < 5000; i++ {
		lim.IsLimited(key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lim.IsLimited(key)
	}
}

func BenchmarkLimiterManager_Increment(b *testing.B) {
	key := 1
	lim := New(200, time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lim.Increment(key)
	}
}
