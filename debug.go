package goappmonitor

import (
	"runtime/debug"
	"time"

	"github.com/wgliang/metrics"
)

var (
	debugMetrics DebugMetrics
	gcStats      = debug.GCStats{Pause: make([]time.Duration, 11)}
)

// Collect denug gc stats and it will run a goroutine.
func CollectDebugGCStats(r metrics.Collectry, d time.Duration) {
	collectDebugGCStats(r)
	go captureDebugGCStats(r, d)
}

// Capture debug status.
func captureDebugGCStats(r metrics.Collectry, d time.Duration) {
	for range time.Tick(d) {
		captureDebugGCStatsWorker(r)
	}
}

// Capture debug status unit.
func captureDebugGCStatsWorker(r metrics.Collectry) {
	lastGC := gcStats.LastGC
	t := time.Now()
	debug.ReadGCStats(&gcStats)
	debugMetrics.ReadGCStats.Update(int64(time.Since(t)))

	debugMetrics.GCStats.LastGC.Update(int64(gcStats.LastGC.UnixNano()))
	debugMetrics.GCStats.NumGC.Update(int64(gcStats.NumGC))
	if lastGC != gcStats.LastGC && 0 < len(gcStats.Pause) {
		debugMetrics.GCStats.Pause.Update(int64(gcStats.Pause[0]))
	}

	debugMetrics.GCStats.PauseTotal.Update(int64(gcStats.PauseTotal))
}

// Collect debug status.
func collectDebugGCStats(r metrics.Collectry) {
	debugMetrics.GCStats.LastGC = metrics.NewGauge()
	debugMetrics.GCStats.NumGC = metrics.NewGauge()
	debugMetrics.GCStats.Pause = metrics.NewHistogram(metrics.NewExpDecaySample(1028, 0.015))
	debugMetrics.GCStats.PauseTotal = metrics.NewGauge()
	debugMetrics.ReadGCStats = metrics.NewHistogram(metrics.NewExpDecaySample(1028, 0.015))

	r.Collector("debug.GCStats.LastGC", debugMetrics.GCStats.LastGC)
	r.Collector("debug.GCStats.NumGC", debugMetrics.GCStats.NumGC)
	r.Collector("debug.GCStats.Pause", debugMetrics.GCStats.Pause)
	r.Collector("debug.GCStats.PauseTotal", debugMetrics.GCStats.PauseTotal)
	r.Collector("debug.ReadGCStats", debugMetrics.ReadGCStats)
}
