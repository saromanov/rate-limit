package limit

import (
	"sync"
	"time"
)

// Limiter defines limiting of the rate
type Limiter struct {
	limit    uint
	interval time.Duration
	m        sync.Mutex
}
