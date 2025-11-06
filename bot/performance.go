package bot

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

type PerformanceMonitor struct {
	searchTimes    []time.Duration
	embeddingTimes []time.Duration
	mutex          sync.RWMutex
	maxSamples     int
}

func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		searchTimes:    make([]time.Duration, 0),
		embeddingTimes: make([]time.Duration, 0),
		maxSamples:     100, // Keep last 100 samples
	}
}

func (pm *PerformanceMonitor) RecordSearchTime(duration time.Duration) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.searchTimes = append(pm.searchTimes, duration)
	if len(pm.searchTimes) > pm.maxSamples {
		pm.searchTimes = pm.searchTimes[1:] // Remove oldest
	}
}

func (pm *PerformanceMonitor) RecordEmbeddingTime(duration time.Duration) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.embeddingTimes = append(pm.embeddingTimes, duration)
	if len(pm.embeddingTimes) > pm.maxSamples {
		pm.embeddingTimes = pm.embeddingTimes[1:] // Remove oldest
	}
}

func (pm *PerformanceMonitor) GetStats() (searchAvg, embeddingAvg time.Duration, memUsage string) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	// Calculate search time average
	if len(pm.searchTimes) > 0 {
		var total time.Duration
		for _, t := range pm.searchTimes {
			total += t
		}
		searchAvg = total / time.Duration(len(pm.searchTimes))
	}

	// Calculate embedding time average
	if len(pm.embeddingTimes) > 0 {
		var total time.Duration
		for _, t := range pm.embeddingTimes {
			total += t
		}
		embeddingAvg = total / time.Duration(len(pm.embeddingTimes))
	}

	// Get memory usage
	var m runtime.MemStats
	runtime.GC() // Force garbage collection for accurate reading
	runtime.ReadMemStats(&m)
	memUsage = formatBytes(m.Alloc)

	return
}

func (pm *PerformanceMonitor) LogPerformanceStats() {
	searchAvg, embeddingAvg, memUsage := pm.GetStats()

	log.Printf("📊 Performance Stats:")
	log.Printf("   Search avg: %v", searchAvg)
	log.Printf("   Embedding avg: %v", embeddingAvg)
	log.Printf("   Memory usage: %s", memUsage)
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return "< 1 KB"
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Start performance monitoring goroutine
func (pm *PerformanceMonitor) StartMonitoring(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			pm.LogPerformanceStats()
		}
	}()
}
