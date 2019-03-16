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

const (
	limitReached = iota + 1
	allowed
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

	errClosedInterval = errors.New("closed interval")
)

// Limiter defines limiting of the rate
type Limiter struct {
	limit      uint32
	interval   time.Duration
	m          sync.RWMutex
	counter    uint32
	metaData   map[string]bool
	message    Status
	afterLimit time.Duration
	reached    chan int
}

// New provides initialization of the Limiter
func New(c *Config) *Limiter {
	lim := &Limiter{
		limit:      c.Limit,
		interval:   c.Interval,
		metaData:   make(map[string]bool),
		message:    NoneMessage,
		afterLimit: c.AfterLimit,
		reached:    make(chan int),
	}

	go lim.events()
	return lim
}

// events consumes notifications about rate limit is reached
// and after this, its starts timer for checking allowed limit time
func (r *Limiter) events() {
	for data := range r.reached {
		switch data {
		case limitReached:
			go r.timer()
			continue
		case allowed:
			atomic.StoreUint32(&r.counter, 0)
		}
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
// it checks if limit is reached
func (r *Limiter) metaSelect() {
	r.m.RLock()
	defer r.m.RUnlock()
	value, ok := r.metaData[metaLimit]
	if ok && value {
		r.metaData[metaLimit] = false
	}
}

// apply method first, checks how many hits need to reach interval
// and then changed status of the Limiter to limitReached
func (r *Limiter) apply() (time.Duration, error) {
	r.m.Lock()
	defer r.m.Unlock()
	l := atomic.LoadUint32(&r.counter)
	if l < r.limit {
		atomic.AddUint32(&r.counter, 1)
		return 0, nil
	}
	r.message = LimitReached
	r.reached <- limitReached
	lim, ok := r.metaData[metaLimit]
	if ok && !lim {
		r.metaData[metaLimit] = true
		return 0, nil
	}
	r.metaData[metaLimit] = true
	return r.interval, errors.New("closed interval")
}

// timer starts when limit is reached
// and stops when AfterLimit interval is passed
func (r *Limiter) timer() {
	timer := time.NewTimer(r.afterLimit)
	<-timer.C
	go func() {
		r.reached <- allowed
	}()
}
