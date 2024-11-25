package info

import (
	"time"
)

// InvocationTimestamps holds the timestamps of all invocations for a specific function.
type InvocationTimestamps struct {
	// HashApp      string    // Application identifier
	HashFunction string    // Function identifier
	Timestamps   []time.Time   // List of invocation timestamps (in seconds since epoch)
	Duration     []time.Duration   // List of invocation durations (in seconds since epoch)
}