package limit

import (
	"sync"
	"sync/atomic"
	"time"
)

// Limiter defines limiting of the rate
type Limiter struct {
	limit    uint
	interval time.Duration
	m        sync.Mutex
	counter  uint32
}

// New provides initialization of the Limiter
func New(limit uint, interval time.Duration) *Limiter {
	return &Limiter{
		limit:    limit,
		interval: interval,
	}
}

// Do provides trying to check limiter
func (r *Limiter) Do() {
	for {
		ok, remaining := r.Try()
		if ok {
			break
		}
		time.Sleep(remaining)
	}
}

func (r *Limiter) apply() (ok bool, remaining time.Duration) {
	r.m.Lock()
	defer r.m.Unlock()
	now := time.Now()
	var cr uint32
	c := atomic.LoadUint32(&cr)
	if l := c; l < r.limit {
		atomic.AddUint32(&r.couner, 1)
		return true, 0
	}
	frnt := r.times.Front()
	if diff := now.Sub(frnt.Value.(time.Time)); diff < r.interval {
		return false, r.interval - diff
	}
	frnt.Value = now
	r.times.MoveToBack(frnt)
	return true, 0
}
