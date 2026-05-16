package status

import (
	"runtime"
	"sync/atomic"
	"time"
)

type Metrics struct {
	Healthy        bool    `json:"healthy"`
	ActiveRequests int64   `json:"active_requests"`
	AvgLatencyMs   float64 `json:"avg_latency_ms"`
	TotalRequests  int64   `json:"total_requests"`
	TotalErrors    int64   `json:"total_errors"`
	ErrorRate      float64 `json:"error_rate"`
	MemoryMB       uint64  `json:"memory_mb"`
}

var (
	activeRequests atomic.Int64
	totalRequests  atomic.Int64
	totalErrors    atomic.Int64
	totalLatencyNs atomic.Int64
)

func RequestStarted() {
	activeRequests.Add(1)
	totalRequests.Add(1)
}

func RequestFinished(latency time.Duration) {
	activeRequests.Add(-1)
	totalLatencyNs.Add(latency.Nanoseconds())
}

func RecordErrors() {
	totalErrors.Add(1)
}

func Status() Metrics {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	reqs := totalRequests.Load()
	errs := totalErrors.Load()

	var averageLatency float64
	if reqs > 0 {
		averageLatency = float64(totalLatencyNs.Load()) / float64(reqs) / 1_000_000
	}

	var errorRate float64
	if reqs > 0 {
		errorRate = float64(errs) / float64(reqs)
	}

	return Metrics{
		Healthy:        true,
		ActiveRequests: activeRequests.Load(),
		AvgLatencyMs:   averageLatency,
		TotalRequests:  reqs,
		TotalErrors:    errs,
		ErrorRate:      errorRate,
		MemoryMB:       mem.Alloc / 1024 / 1024,
	}
}

