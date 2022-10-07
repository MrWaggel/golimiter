package golimiter

import (
	"github.com/akyoto/cache"
	"time"
)

// Limiter is an interface to the limiterManager structure.
// use New(int, time.Duration) to
// instantiate a new limiter. See limiterManager (scroll down)
// for more information.
type Limiter interface {
	// IsLimited checks if the given key is limited
	IsLimited(key interface{}) bool
	// Increment increments the counter
	Increment(key interface{})
	// Remove clears key from all limits
	Remove(key interface{})
	// Count returns the total amount of values
	// that are not expired
	Count(key interface{}) int
}

// limiterManager contains all the internals for the rate-limiter.
type limiterManager struct {
	c           *cache.Cache
	limit       int
	limitDouble int
	duration    time.Duration
	durationInt int64
}

// New returns a limiter manager, this is thread safe but does not
// lock keys like transactions would.
// 		limit: 			how many unexpired entries a key can have
// 		expiresAfter: 	how long before an entry expires for a key
// 		gcInterval: 	the interval when the underlying cache
//						should remove expired keys
func New(limit int, expiresAfter time.Duration) *limiterManager {
	o := new(limiterManager)
	o.limit = limit
	o.limitDouble = o.limit * 2
	o.duration = expiresAfter
	o.durationInt = int64(expiresAfter / time.Second)
	o.c = cache.New(expiresAfter + time.Second)
	return o
}

// IsLimited checks if given key is limited currently limited.
func (l *limiterManager) IsLimited(key interface{}) bool {
	val, exists := l.c.Get(key)
	if !exists {
		return false
	}
	return val.(limitEntry).isLimited(l)
}

// Increment adds a timestamp to the given key, thread safe but concurrent
// calls are not locked as transactions would, so multiple Increments at the
// same time for the same key may result in only one update.
func (l *limiterManager) Increment(key interface{}) {
	var ce limitEntry
	val, exists := l.c.Get(key)
	if exists {
		ce = val.(limitEntry)

		// prune expired entries, this is useful
		// to reduce memory on highly populated caches
		// with high limit
		if len(ce) >= l.limitDouble {
			ce = ce.prune(l)
		}
	} else {
		ce = make(limitEntry, 0)
	}

	ce.addTimeStamp()
	l.c.Set(key, ce, l.duration)
}

// Remove removes all entries for a given key.
func (l *limiterManager) Remove(key interface{}) {
	l.c.Delete(key)
}

// Count returns the total amount of unexpired entries
// of the given key. Returns zero if key is not found.
func (l *limiterManager) Count(key interface{}) int {
	val, exists := l.c.Get(key)
	if !exists {
		return 0
	}
	return val.(limitEntry).count(l)
}

// limitEntry is of type []int64 and holds timestamps that
// are associated with a key.
type limitEntry []int64

// addTimeStamp adds a time stamp to the slice
func (le *limitEntry) addTimeStamp() {
	*le = append(*le, time.Now().Unix())
}

// isLimited counts all the 64bit integers which are not expired yet
// and returns true if the count is >= the limit.
func (le limitEntry) isLimited(parent *limiterManager) bool {
	count := 0
	tb := time.Now().Unix() - parent.durationInt
	for _, v := range le {
		if v > tb {
			count++
		}
		if count >= parent.limit {
			return true
		}
	}
	return false
}

// count returns a sum of all unexpired entries in the limiterEntry
func (le limitEntry) count(parent *limiterManager) int {
	count := 0
	tb := time.Now().Unix() - parent.durationInt
	for _, v := range le {
		if v > tb {
			count++
		}
	}
	return count
}

// prune removes timestamps that are expired and leaves non expired
// timestamps in the array
func (le limitEntry) prune(parent *limiterManager) limitEntry {
	tb := time.Now().Unix() - parent.durationInt
	var cutAt int
	for i, v := range le {
		if v > tb {
			cutAt = i
			break
		}
	}
	if cutAt == 0 {
		return le
	}
	n := make(limitEntry, len(le)-cutAt)
	copy(n, le[cutAt:])
	return n
}
