package golimiter

import (
	"testing"
	"time"
)

func TestLimiterManager(t *testing.T) {
	target := 20
	key := 1
	lim := New(200, time.Minute)

	for i := 0; i < target; i++ {
		lim.Increment(key)
	}

	if lim.Count(key) != target {
		t.Error("did not match target count")
	}
}

func TestLimiterManagerDelayed(t *testing.T) {
	target := 20
	key := 1
	lim := New(200, time.Second*2)

	for i := 0; i < target; i++ {
		lim.Increment(key)
	}

	time.Sleep(time.Second * 3)
	if lim.Count(key) != 0 {
		t.Error("did not match target count")
	}
}

func TestLimiterManagerRemove(t *testing.T) {
	target := 20
	key := 1
	lim := New(200, time.Second*2)

	for i := 0; i < target; i++ {
		lim.Increment(key)
	}

	lim.Remove(key)

	if lim.Count(key) != 0 {
		t.Error("did not match target count")
	}
}
