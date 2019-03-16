package limit

import "time"

// Config defines configuration for the app
type Config struct {
	Limit      uint32
	Interval   time.Duration
	StopLimit  uint32
	AfterLimit time.Duration
}
