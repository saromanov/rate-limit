package limit

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Limiter defines limiting of the rate
type Limiter struct {
	limit    uint32
	interval time.Duration
	m        sync.Mutex
	counter  uint32
}

// New provides initialization of the Limiter
func New(limit uint32, interval time.Duration) *Limiter {
	return &Limiter{
		limit:    limit,
		interval: interval,
	}
}

// Do provides trying to check limiter
func (r *Limiter) Do() {
	for {
		remaining, err := r.apply()
		if err == nil {
			break
		}
		time.Sleep(remaining)
	}
}

func (r *Limiter) apply() (time.Duration, error) {
	r.m.Lock()
	defer r.m.Unlock()
	if l := atomic.LoadUint32(&r.counter); l < r.limit {
		atomic.AddUint32(&r.counter, 1)
		return 0, nil
	}
	return 0, errors.New("closed interval")
}
