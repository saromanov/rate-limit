# rate-limit
[![GoDoc](https://godoc.org/github.com/saromanov/rate-limit?status.png)](https://godoc.org/github.com/saromanov/rate-limit)
[![Go Report Card](https://goreportcard.com/badge/github.com/saromanov/rate-limit)](https://goreportcard.com/report/github.com/saromanov/rate-limit)
[![Build Status](https://travis-ci.org/saromanov/rate-limit.svg?branch=master)](https://travis-ci.org/saromanov/rate-limit)
[![Coverage Status](https://coveralls.io/repos/github/saromanov/rate-limit/badge.svg?branch=master)](https://coveralls.io/github/saromanov/rate-limit?branch=master)


Implementation of simple rate limiter

### Example
```go
package main

import (
	"fmt"
	"time"

	"github.com/saromanov/rate-limit"
)

func main() {
	lim := limit.New(&limit.Config{
		Limit: 10, 
		Interval: 2*time.Second,
	})
	for index := 0; index < 20; index++ {
		fmt.Println(index)
		lim.Do()
	}
}
```

After 10-th iteration it'll slow down iteration and each will be after 2 seconds