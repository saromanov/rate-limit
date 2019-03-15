package limit

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const (
	metaLimit    = "limit"
	metaPreLimit = "prelimit"
)

// Status returns information from rate limiter
type Status string

var (
	// NoneMessage is defaul message from limiter
	NoneMessage Status
	// PreLimit returns in the case when limit is nealy reached
	PreLimit Status = "You around the limit"
	// LimitReached returns when limit is reached
	LimitReached Status = "limit is reached"
)

// Limiter defines limiting of the rate
type Limiter struct {
	limit    uint32
	interval time.Duration
	m        sync.RWMutex
	counter  uint32
	metaData map[string]bool
	message  Status
}

// New provides initialization of the Limiter
func New(c *Config) *Limiter {
	return &Limiter{
		limit:    c.Limit,
		interval: c.Interval,
		metaData: make(map[string]bool),
		message:  NoneMessage,
	}
}

// Status returns info from limiter
func (r *Limiter) Status() Status {
	return r.message
}

// Do provides trying to check limiter
func (r *Limiter) Do() {
	for {
		remaining, err := r.apply()
		if err == nil {
			break
		}
		time.Sleep(remaining)
		r.metaSelect()
	}
}

// metaSelect provides handling of meta data
func (r *Limiter) metaSelect() {
	r.m.RLock()
	defer r.m.RUnlock()
	value, ok := r.metaData[metaLimit]
	if ok && value {
		r.metaData[metaLimit] = false
	}
}

func (r *Limiter) apply() (time.Duration, error) {
	r.m.Lock()
	defer r.m.Unlock()
	l := atomic.LoadUint32(&r.counter)
	if l < r.limit {
		atomic.AddUint32(&r.counter, 1)
		return 0, nil
	}
	r.message = LimitReached
	lim, ok := r.metaData[metaLimit]
	if ok && !lim {
		r.metaData[metaLimit] = true
		return 0, nil
	}
	r.metaData[metaLimit] = true
	return r.interval, errors.New("closed interval")
}
