package wrapper

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
	"regexp"
	"strings"
	"trace-analyser/pkg/info"
)


// ParseAndConvert reads the input CSV file, processes invocation counts, and converts them into timestamps.
func ParseAndConvert(invocationFilePath string, durationFilePath string) ([]info.InvocationTimestamps, error) {
	// Open the CSV file
	invoFile, err := os.Open(invocationFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer invoFile.Close()
	duraFile, err := os.Open(durationFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer duraFile.Close()

	// Parse the CSV file
	invoReader := csv.NewReader(invoFile)
	invoRows, err := invoReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}
	duraReader := csv.NewReader(duraFile)
	duraRows, err := duraReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	var results []info.InvocationTimestamps
	blocks := strings.Split(invocationFilePath, ".")
	re := regexp.MustCompile(`d(\d+)`)
	matches := re.FindStringSubmatch(blocks[len(blocks)-2])
	day, _ := strconv.Atoi(matches[1])
	startOfDay := time.Date(2019, 7, day, 0, 0, 0, 0, time.UTC) // Set a fixed date for calculation

	// Process invocation
	for i, row := range invoRows {
		fmt.Println(i)
		// Skip header row
		if i == 0 {
			continue
		}
		// hashApp := row[1]
		hashFunction := row[1] + row[2]

		// Search in durationFile to find the corresponding duration
		var durations []time.Duration
		var duration time.Duration
		for j, duraRow := range duraRows {
			if duraRow[1]+duraRow[2] == hashFunction {
				d, err := strconv.Atoi(duraRow[3])
				if err != nil {
					return nil, fmt.Errorf("invalid duration at line %d, column %d: %w", j+1, 3, err)
				}
				duration = time.Duration(d) * time.Millisecond // Convert to time.Duration
				break
			}
		}

		// Convert invocation counts to timestamps
		var timestamps []time.Time
		for minute, countStr := range row[4:] {
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return nil, fmt.Errorf("invalid invocation count at line %d, column %d: %w", i+1, minute+5, err)
			}

			if count > 0 {
				// Calculate timestamps for the current minute
				minuteStart := startOfDay.Add(time.Duration(minute) * time.Minute)
				interval := time.Second * 60 / time.Duration(count) // Interval in time.Duration

				for j := 0; j < count; j++ {
					invocationTime := minuteStart.Add(time.Duration(j) * interval)
					timestamps = append(timestamps, invocationTime)
					durations = append(durations, duration)
				}
			}
		}


		results = append(results, info.InvocationTimestamps{
			// HashApp:      hashApp,
			HashFunction: hashFunction,
			Timestamps:   timestamps,
			Duration: 	  durations,
		})
	}

	return results, nil
}
