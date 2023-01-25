# golimiter
[![Go Reference](https://pkg.go.dev/static/frontend/badge/badge.svg)](https://pkg.go.dev/github.com/mrwaggel/golimiter)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrwaggel/golimiter)](https://goreportcard.com/report/github.com/mrwaggel/golimiter)

This packages provides a key based limiter. [Read more here.](https://mrwaggel.be/page/golimiter)

## Usage
### go get
```
go get github.com/mrwaggel/golimiter
```

### example
```go
package main

import (
	"github.com/mrwaggel/golimiter"
	"time"
)

func main() {
	l := golimiter.New(4, time.Second*5)
	key := "a"

	l.Increment(key)
	l.Increment(key)
	l.Increment(key)

	l.Count(key)     // 3
	l.IsLimited(key) // false

	l.Increment(key)

	l.Count(key)     // 4
	l.IsLimited(key) // true

	time.Sleep(time.Second * 6)
	l.Count(key)     // 0
	l.IsLimited(key) // false
}
```