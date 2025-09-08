package common

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// Thread-safe result collector
type ResultCollector struct {
	mu           sync.RWMutex
	failedLines  map[int]FailedLineInfo
	successCount int64
	failCount    int64
	skipCount    int64
	coinStats    map[string]map[string]int
}

// Create result collector
func NewResultCollector() *ResultCollector {
	return &ResultCollector{
		failedLines: make(map[int]FailedLineInfo),
		coinStats:   make(map[string]map[string]int),
	}
}

// Add verification result
func (rc *ResultCollector) AddResult(result VerifyResult) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// Initialize coin statistics
	if rc.coinStats[result.Coin] == nil {
		rc.coinStats[result.Coin] = make(map[string]int)
	}

	if result.Success {
		atomic.AddInt64(&rc.successCount, 1)
		rc.coinStats[result.Coin]["success"]++
	} else {
		atomic.AddInt64(&rc.failCount, 1)
		rc.coinStats[result.Coin]["fail"]++

		// Record failed line detailed information
		rc.failedLines[result.Line.LineNumber] = FailedLineInfo{
			LineNumber:   result.Line.LineNumber,
			Coin:         result.Coin,
			DigitalAsset: result.Line.DigitalAsset,
			Address:      result.Line.Address,
			ErrorMessage: result.Error,
		}
	}
}

// Add skip count
func (rc *ResultCollector) AddSkip() {
	atomic.AddInt64(&rc.skipCount, 1)
}

// Get statistics
func (rc *ResultCollector) GetStats() (int64, int64, int64, map[int]FailedLineInfo, map[string]map[string]int) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	// Deep copy failed lines
	failedLinesCopy := make(map[int]FailedLineInfo)
	for k, v := range rc.failedLines {
		failedLinesCopy[k] = v
	}

	// Deep copy coin statistics
	coinStatsCopy := make(map[string]map[string]int)
	for coin, stats := range rc.coinStats {
		coinStatsCopy[coin] = make(map[string]int)
		for status, count := range stats {
			coinStatsCopy[coin][status] = count
		}
	}

	return atomic.LoadInt64(&rc.successCount), atomic.LoadInt64(&rc.failCount),
		atomic.LoadInt64(&rc.skipCount), failedLinesCopy, coinStatsCopy
}

// Worker pool structure
type WorkerPool struct {
	workerCount int
	jobChan     chan CSVLine
	wg          sync.WaitGroup
	collector   *ResultCollector
	t           *testing.T
}

// Create worker pool
func NewWorkerPool(workerCount int, collector *ResultCollector, t *testing.T) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		jobChan:     make(chan CSVLine, workerCount*2), // Buffer size is twice the worker count
		collector:   collector,
		t:           t,
	}
}

// Start worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Worker function
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for line := range wp.jobChan {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Caught panic, convert to error handling
					errorMsg := fmt.Sprintf("Verification panic occurred: %v", r)
					result := VerifyResult{
						Line:    line,
						Success: false,
						Coin:    line.Network, // Use network as coin
						Error:   errorMsg,
					}
					wp.collector.AddResult(result)
				}
			}()

			// Verify this line (multithreaded version, skip StarkNet)
			success, coin, errorMsg := verifyCSVLineMultithread(line.Network, line.Address, line.Message,
				line.SignedMessage, line.SignedMessage2, line.PublicKey, line.Owner1, line.Owner2, line.DigitalAsset, line.LineNumber, wp.t)

			// Send result
			result := VerifyResult{
				Line:    line,
				Success: success,
				Coin:    coin,
				Error:   errorMsg,
			}

			wp.collector.AddResult(result)
		}()
	}
}

// Add job to worker pool
func (wp *WorkerPool) AddJob(line CSVLine) {
	wp.jobChan <- line
}

// Stop worker pool
func (wp *WorkerPool) Stop() {
	close(wp.jobChan)
	wp.wg.Wait()
}
