package main

import (
	"fmt"
	"log"
	// "trace-analyser/pkg/logic"
	// "trace-analyser/pkg/plot"
	"trace-analyser/pkg/wrapper"
	// "time"
)

func main() {
	// Path to the input CSV file
	traceFile := "data/azure-2019/invocations_per_function_md.anon.d01.csv"
	duraFile := "data/azure-2019/function_durations_percentiles.anon.d01.csv"

	// Step 1: Process the trace file and get invocation data
	invocationTimestamps, err := wrapper.ParseAndConvert(traceFile, duraFile)
	if err != nil {
		log.Fatalf("Error processing trace file: %v", err)
	}

	// // Step 2: Analyze cold starts for each function
	// analyzer := logic.ColdStartAnalyzer{KeepAlive: time.Second * 60}
	// // Analyze cold starts
	// coldStartTimestamps, err := analyzer.AnalyzeColdStarts(invocationTimestamps)
	// if err != err {
	// 	log.Fatalf("Error calculating coldstarts", err)
	// }

	// // Step 3: Plot cold start statistics
	// startOfDay := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Fixed day for calculation
	// outputFile := "cold_starts_per_minute.png"

	// err = plot.PlotColdStarts(coldStartTimestamps, startOfDay, outputFile)
	// if err != nil {
	// 	log.Fatalf("Error creating plot: %v", err)
	// }

	// fmt.Printf("Cold start statistics plotted successfully: %s\n", outputFile)
	fmt.Println(invocationTimestamps)
}
