package common

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
)

// Multithreaded CSV file verification
func TestVerifyCSVFileMultithread(t *testing.T) {
	csvFileName := "../example/test.csv"

	// Temporary marking function switch - can be deleted after testing
	enableMarking := true // Set to false to disable marking function
	batchSize := 100000   // Save every 100k lines

	// Get CPU core count, set worker thread count
	workerCount := runtime.NumCPU()

	fmt.Printf("Using %d worker threads for verification\n", workerCount)

	var existingFailedLines map[int]FailedLineInfo
	var lastProcessedLine int
	var resultFileName string
	if enableMarking {
		resultFileName = csvFileName + "_result"
		existingFailedLines, lastProcessedLine = loadVerificationResult(resultFileName)
		fmt.Printf("Loaded verification results: processed to line %d, %d failed lines need re-verification\n", lastProcessedLine, len(existingFailedLines))
	}

	fmt.Printf("Starting multithreaded CSV file verification: %s\n", csvFileName)

	// First calculate total lines (for progress bar)
	fmt.Printf("Calculating total file lines...\n")
	totalLines := 0
	file, err := os.Open(csvFileName)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		totalLines++
	}
	file.Close()
	fmt.Printf("Total lines: %d\n", totalLines)

	// Create progress bar
	var progressBar *ProgressBar
	if totalLines > 1000 { // Only show progress bar for large files
		progressBar = NewProgressBar(totalLines - 1) // Exclude header line
	}

	// Create result collector and worker pool
	collector := NewResultCollector()
	workerPool := NewWorkerPool(workerCount, collector, t)

	// Start worker pool
	workerPool.Start()

	// Open file for processing
	file, err = os.Open(csvFileName)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	lineNumber := 0
	processedCount := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and header
		if line == "" || lineNumber == 1 {
			continue
		}

		// Skip already processed lines (unless they are failed lines)
		if enableMarking && lineNumber <= lastProcessedLine {
			if _, isFailed := existingFailedLines[lineNumber]; !isFailed {
				collector.AddSkip()
				continue
			}
		}

		// Parse CSV line
		fields := strings.Split(line, ",")
		if len(fields) < 7 {
			// Format error, add directly to failed results
			result := VerifyResult{
				Line: CSVLine{
					LineNumber:   lineNumber,
					DigitalAsset: "UNKNOWN",
					Network:      "UNKNOWN",
					Address:      "UNKNOWN",
					RawLine:      line,
				},
				Success: false,
				Coin:    "UNKNOWN",
				Error:   fmt.Sprintf("Format error, insufficient fields: %s", line),
			}
			collector.AddResult(result)
			continue
		}

		// Create CSV line data
		csvLine := CSVLine{
			LineNumber:     lineNumber,
			DigitalAsset:   strings.TrimSpace(fields[0]),
			Network:        strings.TrimSpace(fields[1]),
			Address:        strings.TrimSpace(fields[2]),
			SignedMessage:  strings.TrimSpace(fields[3]),
			SignedMessage2: strings.TrimSpace(fields[4]),
			Message:        strings.TrimSpace(fields[5]),
			PublicKey:      strings.TrimSpace(fields[6]),
			RawLine:        line,
		}

		// Parse owner1 and owner2 fields (if exist)
		if len(fields) > 7 {
			csvLine.Owner1 = strings.TrimSpace(fields[7])
		}
		if len(fields) > 8 {
			csvLine.Owner2 = strings.TrimSpace(fields[8])
		}

		// Send job to worker pool
		workerPool.AddJob(csvLine)
		processedCount++

		// Update progress bar
		if progressBar != nil {
			progressBar.Update(lineNumber - 1) // Subtract header line
		}

		// Save results every batch
		if enableMarking && processedCount%batchSize == 0 {
			successCount, failCount, skipCount, failedLines, _ := collector.GetStats()
			saveVerificationResult(resultFileName, failedLines, lineNumber)
			fmt.Printf("\nSaved verification results at line %d (Success: %d, Failed: %d, Skipped: %d)\n",
				lineNumber, successCount, failCount, skipCount)
		}
	}

	// Stop worker pool and wait for completion
	workerPool.Stop()

	// Finish progress bar
	if progressBar != nil {
		progressBar.Finish()
	}

	// Get final statistics
	successCount, failCount, skipCount, failedLines, coinStats := collector.GetStats()

	// Final save
	if enableMarking {
		saveVerificationResult(resultFileName, failedLines, lineNumber)
	}

	// Display results
	fmt.Printf("\n=== Verification Results ===\n")
	fmt.Printf("Total processed: %d lines\n", lineNumber-1)
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failCount)
	fmt.Printf("Skipped: %d\n", skipCount)

	// Display coin statistics
	fmt.Printf("\n=== Coin Statistics ===\n")
	for coin, stats := range coinStats {
		total := stats["success"] + stats["fail"]
		if total > 0 {
			fmt.Printf("%s: Total %d, Success %d, Failed %d\n", coin, total, stats["success"], stats["fail"])
		}
	}

	if failCount == 0 {
		fmt.Println("\nAll addresses verified successfully!")
		fmt.Printf("CSV file preserved: %s\n", csvFileName)
		if enableMarking {
			// Only clean up result files, keep the original CSV
			os.Remove(resultFileName)
			fmt.Printf("Verification result file deleted: %s\n", resultFileName)
		}
	} else {
		fmt.Printf("\n%d addresses failed verification, CSV file preserved for debugging\n", failCount)
	}
}

// Single-threaded CSV file verification (StarkNet only)
func TestVerifyCSVFileStarknetOnly(t *testing.T) {
	csvFileName := "../example/test.csv"

	// Temporary marking function switch - can be deleted after testing
	enableMarking := true // Set to false to disable marking function
	batchSize := 100000   // Save every 100k lines

	var failedLines map[int]FailedLineInfo
	var lastProcessedLine int
	var resultFileName string
	if enableMarking {
		resultFileName = csvFileName + "_result"
		failedLines, lastProcessedLine = loadVerificationResult(resultFileName)
		fmt.Printf("Loaded verification results: processed to line %d, %d failed lines need re-verification\n", lastProcessedLine, len(failedLines))
	}

	fmt.Printf("Starting CSV file verification (StarkNet only): %s\n", csvFileName)

	// First calculate total lines (for progress bar)
	fmt.Printf("Calculating total file lines...\n")
	totalLines := 0
	file, err := os.Open(csvFileName)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		totalLines++
	}
	file.Close()
	fmt.Printf("Total lines: %d\n", totalLines)

	// Create progress bar
	var progressBar *ProgressBar
	if totalLines > 1000 { // Only show progress bar for large files
		progressBar = NewProgressBar(totalLines - 1) // Exclude header line
	}

	// Open file again for processing
	file, err = os.Open(csvFileName)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	lineNumber := 0
	successCount := 0
	failCount := 0
	skipCount := 0

	// Coin statistics
	coinStats := make(map[string]map[string]int)

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and header
		if line == "" || lineNumber == 1 {
			continue
		}

		// Skip already processed lines (unless they are failed lines)
		if enableMarking && lineNumber <= lastProcessedLine {
			if _, isFailed := failedLines[lineNumber]; !isFailed {
				skipCount++
				continue
			}
		}

		// Parse CSV line
		fields := strings.Split(line, ",")
		if len(fields) < 7 {
			errorMsg := fmt.Sprintf("Format error, insufficient fields: %s", line)
			t.Logf("Line %d: %s", lineNumber, errorMsg)
			failCount++
			if enableMarking {
				// Record format error details
				failedInfo := FailedLineInfo{
					LineNumber:   lineNumber,
					Coin:         "UNKNOWN",
					DigitalAsset: "UNKNOWN",
					Address:      "UNKNOWN",
					ErrorMessage: errorMsg,
				}
				failedLines[lineNumber] = failedInfo
			}
			continue
		}

		digitalAsset := strings.TrimSpace(fields[0]) // Keep digitalAsset field for error messages
		network := strings.TrimSpace(fields[1])      // Use network field as coin
		address := strings.TrimSpace(fields[2])
		signedMessage := strings.TrimSpace(fields[3])
		signedMessage2 := strings.TrimSpace(fields[4])
		message := strings.TrimSpace(fields[5])
		publicKey := strings.TrimSpace(fields[6])

		// Verify this line, add panic protection
		var success bool
		var coin string
		var errorMsg string

		func() {
			defer func() {
				if r := recover(); r != nil {
					// Caught panic, convert to error handling
					success = false
					coin = network
					errorMsg = fmt.Sprintf("Verification panic occurred: %v", r)
					t.Logf("Line %d panic occurred during verification: %v (digitalAsset:%s, network:%s, addr:%s)", lineNumber, r, digitalAsset, network, address)
				}
			}()

			success, coin, errorMsg = verifyCSVLineStarknetOnly(network, address, message, signedMessage, signedMessage2, publicKey, digitalAsset, lineNumber, t)
		}()

		// Update statistics
		if coinStats[coin] == nil {
			coinStats[coin] = make(map[string]int)
		}
		if success {
			successCount++
			coinStats[coin]["success"]++
		} else {
			failCount++
			coinStats[coin]["fail"]++
			if enableMarking {
				// Record failed details
				failedInfo := FailedLineInfo{
					LineNumber:   lineNumber,
					Coin:         coin,
					DigitalAsset: digitalAsset,
					Address:      address,
					ErrorMessage: errorMsg,
				}
				failedLines[lineNumber] = failedInfo
			}
		}

		// Update progress bar
		if progressBar != nil {
			progressBar.Update(lineNumber - 1) // Subtract header line
		}

		// Save results every 100k lines
		if enableMarking && lineNumber%batchSize == 0 {
			saveVerificationResult(resultFileName, failedLines, lineNumber)
			fmt.Printf("\nSaved verification results at line %d\n", lineNumber)
		}
	}

	// Finish progress bar
	if progressBar != nil {
		progressBar.Finish()
	}

	// Final save
	if enableMarking {
		saveVerificationResult(resultFileName, failedLines, lineNumber)
	}

	// Display results
	fmt.Printf("\n=== Verification Results ===\n")
	fmt.Printf("Total processed: %d lines\n", lineNumber-1)
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failCount)
	fmt.Printf("Skipped: %d\n", skipCount)

	// Display coin statistics
	fmt.Printf("\n=== Coin Statistics ===\n")
	for coin, stats := range coinStats {
		total := stats["success"] + stats["fail"]
		if total > 0 {
			fmt.Printf("%s: Total %d, Success %d, Failed %d\n", coin, total, stats["success"], stats["fail"])
		}
	}

	if failCount == 0 {
		fmt.Println("\nAll addresses verified successfully!")
		fmt.Printf("CSV file preserved: %s\n", csvFileName)
		if enableMarking {
			// Only clean up result files, keep the original CSV
			os.Remove(resultFileName)
			fmt.Printf("Verification result file deleted: %s\n", resultFileName)
		}
	} else {
		fmt.Printf("\n%d addresses failed verification, CSV file preserved for debugging\n", failCount)
	}
}
