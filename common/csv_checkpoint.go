package common

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

// Thread-safe CSV file writing lock
var csvFileLock sync.Mutex

// Load verification result file (only read CSV format failed records)
func loadVerificationResult(resultFileName string) (map[int]FailedLineInfo, int) {
	failedLines := make(map[int]FailedLineInfo)
	lastProcessedLine := 0

	// Convert to CSV filename
	csvFileName := strings.TrimSuffix(resultFileName, "_result") + "_failed.csv"

	file, err := os.Open(csvFileName)
	if err != nil {
		return failedLines, lastProcessedLine // Return empty if file doesn't exist
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || lineCount == 1 { // Skip empty lines and header
			continue
		}

		// Parse CSV format: digitalAsset,network,address,errorMessage
		parts := strings.Split(line, ",")
		if len(parts) >= 3 {
			// Since CSV doesn't have line number info, we use a virtual line number
			// This is mainly for marking these addresses for re-verification
			errorMsg := "Need re-verification"
			if len(parts) >= 4 {
				// Handle error messages that might be surrounded by quotes
				errorMsg = strings.TrimSpace(parts[3])
				if strings.HasPrefix(errorMsg, "\"") && strings.HasSuffix(errorMsg, "\"") {
					errorMsg = strings.TrimPrefix(strings.TrimSuffix(errorMsg, "\""), "\"")
					// Restore CSV escaped double quotes
					errorMsg = strings.ReplaceAll(errorMsg, "\"\"", "\"")
				}
			}

			failedLines[lineCount] = FailedLineInfo{
				LineNumber:   lineCount,
				Coin:         strings.TrimSpace(parts[1]), // network
				DigitalAsset: strings.TrimSpace(parts[0]), // digitalAsset
				Address:      strings.TrimSpace(parts[2]), // address
				ErrorMessage: errorMsg,
			}
		}
	}

	return failedLines, lastProcessedLine
}

// Save verification result
func saveVerificationResult(resultFileName string, failedLines map[int]FailedLineInfo, lastProcessedLine int) {
	// Only save CSV format failed records, no txt suffix
	csvFileName := strings.TrimSuffix(resultFileName, "_result") + "_failed.csv"

	if len(failedLines) > 0 {
		saveFailedAsCSV(csvFileName, failedLines)
	} else {
		fmt.Printf("No failed records in current batch, processed to line %d\n", lastProcessedLine)
	}
}

// Save failed records as CSV
func saveFailedAsCSV(csvFileName string, failedLines map[int]FailedLineInfo) {
	csvFileLock.Lock()
	defer csvFileLock.Unlock()

	// Check if file exists and read existing records to avoid duplicates
	existingRecords := make(map[string]bool)
	if file, err := os.Open(csvFileName); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				existingRecords[line] = true
			}
		}
		file.Close()
	}

	// Open file for appending (create if doesn't exist)
	file, err := os.OpenFile(csvFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Failed to create CSV file: %v\n", err)
		return
	}
	defer file.Close()

	// Write CSV header if file is empty
	if len(existingRecords) == 0 {
		fmt.Fprintf(file, "digitalAsset,network,address,errorMessage\n")
	}

	// Collect and sort failed line numbers
	var lineNumbers []int
	for lineNum := range failedLines {
		lineNumbers = append(lineNumbers, lineNum)
	}
	sort.Ints(lineNumbers)

	// Write failed records, including error information
	newRecordsCount := 0
	for _, lineNum := range lineNumbers {
		failedInfo := failedLines[lineNum]
		// Escape special characters in CSV (commas, newlines, double quotes)
		errorMsg := failedInfo.ErrorMessage
		if strings.Contains(errorMsg, ",") || strings.Contains(errorMsg, "\n") || strings.Contains(errorMsg, "\"") {
			errorMsg = "\"" + strings.ReplaceAll(errorMsg, "\"", "\"\"") + "\""
		}

		recordLine := fmt.Sprintf("%s,%s,%s,%s",
			failedInfo.DigitalAsset,
			failedInfo.Coin, // network is coin
			failedInfo.Address,
			errorMsg)

		// Only write if record doesn't already exist
		if !existingRecords[recordLine] {
			fmt.Fprintf(file, "%s\n", recordLine)
			newRecordsCount++
		}
	}

	fmt.Printf("Saved failed records CSV file: %s, new failed lines: %d, total: %d\n", csvFileName, newRecordsCount, len(failedLines))
}
