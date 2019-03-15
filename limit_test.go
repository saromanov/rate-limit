package limit

import (
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	lim := New(&Config{
		Limit:    2,
		Interval: 1 * time.Second,
	})
	for index := 0; index < 10; index++ {
		lim.Do()
		st := lim.Status()
		if index > 1 {
			if st != LimitReached {
				t.Fatal("limit should be reached")
			}
		} else {
			if st != NoneMessage {
				t.Fatal("limit shouldn't be reached")
			}
		}
	}
}
