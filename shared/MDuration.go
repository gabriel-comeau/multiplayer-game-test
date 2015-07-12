package shared

import (
	"time"
)

// Wrap a time.Duration by embedding it for a more convenient GetMillis() call
// (since you can't add methods to out-of-package types)
type MDuration struct {
	time.Duration
}

func (m MDuration) Milliseconds() int64 {
	return m.Nanoseconds() / NANO_TO_MILLI
}
