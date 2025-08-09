package clock

import (
	"sync/atomic"
	"time"
)

// Now returns the current time. It can be overridden for tests/benchmarks.
var Now = time.Now

// UseMonotonicFake overrides Now to return a monotonically increasing time
// starting from the provided start time (or time.Unix(0,0) if start.IsZero())
// advancing by the given step on each call. It returns a restore function
// that resets Now back to time.Now.
func UseMonotonicFake(start time.Time, step time.Duration) (restore func()) {
	if start.IsZero() {
		start = time.Unix(0, 0)
	}
	var counter atomic.Int64
	Now = func() time.Time {
		n := counter.Add(1)
		return start.Add(time.Duration(n) * step)
	}
	return func() { Now = time.Now }
}
