package logic

import (
	"time"
	"fmt"
	"trace-analyser/pkg/info"
)

// ColdStartAnalyzer analyzes invocation data and determines cold start timestamps.
type ColdStartAnalyzer struct {
	KeepAlive time.Duration // Time (e.g., 60 seconds) an instance remains alive after invocation ends
}

type Instance struct {
	LastEndTime time.Time
	ExpiryTime time.Time
}

// AnalyzeColdStarts processes invocation timestamps and durations for multiple functions,
// returning a map where the key is the function identifier and the value is a list of cold start timestamps.
func (c *ColdStartAnalyzer) AnalyzeColdStarts(invocations []info.InvocationTimestamps) ([]time.Time, error) {
	// Results map: key is the HashFunction, value is a slice of cold start timestamps
	results := make([]time.Time, 0)

	for _, invocationData := range invocations {
		// hashFunction := invocationData.HashFunction
		timestamps := invocationData.Timestamps
		durations := invocationData.Duration

		// Check slices' length
		if len(timestamps) != len(durations) {
			return nil, fmt.Errorf("len(timestamps) is %d while len(durations) is %d", len(timestamps), len(durations))
		}

		// Track active instances and their expiry times for this specific function
		activeInstances := make([]Instance, 0) // Key: instance last end time, Value: expiry time
		var coldStartTimestamps []time.Time

		for i, start := range timestamps {
			duration := durations[i]
			instanceFound := false

			// Check if any instance is available
			// for i, instance := range activeInstances {
			for i := len(activeInstances); i >= 0; i-- {
				if start.After(activeInstances[i].LastEndTime) && start.Before(activeInstances[i].ExpiryTime) {
					// Use this instance and update its expiry time
					// Update instance info
					activeInstances[i] = Instance{
						LastEndTime: start.Add(duration),
						ExpiryTime: start.Add(duration).Add(c.KeepAlive),
					}
					instanceFound = true
					continue
				} else if start.After(activeInstances[i].ExpiryTime){
					// Remove instances expired
					activeInstances = append(activeInstances[:i], activeInstances[i+1:]...)
				}

			}

			// If no active instance is found, this is a cold start
			if !instanceFound {
				coldStartTimestamps = append(coldStartTimestamps, start)
				// Create a new instance and set its expiry time
				activeInstances = append(activeInstances, Instance{
					LastEndTime: start.Add(duration),
					ExpiryTime: start.Add(duration).Add(c.KeepAlive),
				})
			}
		}

		// Store cold start timestamps for this specific function
		results = append(results, coldStartTimestamps...)
	}

	return results, nil
}
