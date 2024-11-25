package plot

import (
	"fmt"
	"os"
	"time"

	"github.com/wcharczuk/go-chart/v2"
)

// PlotColdStarts takes the cold start timestamps and generates a chart showing the count per minute.
func PlotColdStarts(coldStartTimestamps []time.Time, startOfDay time.Time, outputFilePath string) error {
	// Count cold starts per minute
	countsPerMinute := make(map[time.Time]int)
	for _, timestamp := range coldStartTimestamps {
		// Truncate timestamp to the nearest minute
		minute := timestamp.Truncate(time.Minute)
		countsPerMinute[minute]++
	}

	// Generate x-axis labels and y-axis counts
	var xLabels []time.Time
	var yCounts []float64

	// Assume a full day (1440 minutes)
	for i := 0; i < 1440; i++ {
		minute := startOfDay.Add(time.Duration(i) * time.Minute)
		xLabels = append(xLabels, minute)
		yCounts = append(yCounts, float64(countsPerMinute[minute]))
	}

	// Create the chart
	c := chart.Chart{
		Width:  1280,
		Height: 720,
		XAxis: chart.XAxis{
			Name:           "Time (Minutes)",
			ValueFormatter: chart.TimeValueFormatter,
		},
		YAxis: chart.YAxis{
			Name: "Cold Starts",
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: xLabels,
				YValues: yCounts,
			},
		},
	}

	// Write the chart to an output file
	file, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Render the chart as PNG
	err = c.Render(chart.PNG, file)
	if err != nil {
		return fmt.Errorf("failed to render chart: %w", err)
	}

	fmt.Printf("Chart saved to %s\n", outputFilePath)
	return nil
}
