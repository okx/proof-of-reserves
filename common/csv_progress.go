package common

import (
	"fmt"
	"time"
)

// Progress bar structure
type ProgressBar struct {
	total      int
	current    int
	startTime  time.Time
	lastUpdate time.Time
}

// Create progress bar
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		total:      total,
		current:    0,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
	}
}

// Update progress
func (pb *ProgressBar) Update(current int) {
	pb.current = current
	now := time.Now()

	// Update progress bar every 5 seconds to avoid too frequent output
	if now.Sub(pb.lastUpdate) < 5*time.Second {
		return
	}
	pb.lastUpdate = now

	// Calculate progress percentage
	percentage := float64(current) / float64(pb.total) * 100

	// Calculate elapsed time
	elapsed := now.Sub(pb.startTime)

	// Calculate estimated remaining time
	var eta time.Duration
	if current > 0 {
		avgTimePerItem := elapsed / time.Duration(current)
		remaining := pb.total - current
		eta = avgTimePerItem * time.Duration(remaining)
	}

	// Calculate processing speed (lines/second)
	speed := float64(current) / elapsed.Seconds()

	// Display progress bar
	fmt.Printf("\rProgress: %d/%d (%.1f%%) | Speed: %.1f lines/s | Elapsed: %v | ETA: %v",
		current, pb.total, percentage, speed, elapsed.Round(time.Second), eta.Round(time.Second))
}

// Finish progress bar
func (pb *ProgressBar) Finish() {
	elapsed := time.Since(pb.startTime)
	speed := float64(pb.current) / elapsed.Seconds()
	fmt.Printf("\nCompleted: %d lines in %v (%.1f lines/s)\n", pb.current, elapsed.Round(time.Second), speed)
}
