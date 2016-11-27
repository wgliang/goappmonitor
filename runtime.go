package goappmonitor

import (
	"runtime"
	"time"

	"github.com/wgliang/metrics"
)

var (
	memStats       runtime.MemStats
	runtimeMetrics RuntimeMetrics
	frees          uint64
	lookups        uint64
	mallocs        uint64
	numGC          uint32
	numCgoCalls    int64
)

// Collect runtime memeory-data:CollectRuntimeMemStats.
func CollectRuntimeMemStats(r metrics.Collectry, d time.Duration) {
	collectRuntimeMemStats(r)
	go captureRuntimeMemStats(r, d)
}

// Capture runtime memeory-data.
func captureRuntimeMemStats(r metrics.Collectry, d time.Duration) {
	for _ = range time.Tick(d) {
		captureRuntimeMemStatsWorker(r)
	}
}

// Capture runtime memeory-data worker.
func captureRuntimeMemStatsWorker(r metrics.Collectry) {
	t := time.Now()
	runtime.ReadMemStats(&memStats)
	runtimeMetrics.ReadMemStats.Update(int64(time.Since(t)))

	runtimeMetrics.MemStats.Alloc.Update(int64(memStats.Alloc))
	runtimeMetrics.MemStats.BuckHashSys.Update(int64(memStats.BuckHashSys))
	if memStats.DebugGC {
		runtimeMetrics.MemStats.DebugGC.Update(1)
	} else {
		runtimeMetrics.MemStats.DebugGC.Update(0)
	}
	if memStats.EnableGC {
		runtimeMetrics.MemStats.EnableGC.Update(1)
	} else {
		runtimeMetrics.MemStats.EnableGC.Update(0)
	}

	runtimeMetrics.MemStats.Frees.Update(int64(memStats.Frees - frees))
	runtimeMetrics.MemStats.HeapAlloc.Update(int64(memStats.HeapAlloc))
	runtimeMetrics.MemStats.HeapIdle.Update(int64(memStats.HeapIdle))
	runtimeMetrics.MemStats.HeapInuse.Update(int64(memStats.HeapInuse))
	runtimeMetrics.MemStats.HeapObjects.Update(int64(memStats.HeapObjects))
	runtimeMetrics.MemStats.HeapReleased.Update(int64(memStats.HeapReleased))
	runtimeMetrics.MemStats.HeapSys.Update(int64(memStats.HeapSys))
	runtimeMetrics.MemStats.LastGC.Update(int64(memStats.LastGC))
	runtimeMetrics.MemStats.Lookups.Update(int64(memStats.Lookups - lookups))
	runtimeMetrics.MemStats.Mallocs.Update(int64(memStats.Mallocs - mallocs))
	runtimeMetrics.MemStats.MCacheInuse.Update(int64(memStats.MCacheInuse))
	runtimeMetrics.MemStats.MCacheSys.Update(int64(memStats.MCacheSys))
	runtimeMetrics.MemStats.MSpanInuse.Update(int64(memStats.MSpanInuse))
	runtimeMetrics.MemStats.MSpanSys.Update(int64(memStats.MSpanSys))
	runtimeMetrics.MemStats.NextGC.Update(int64(memStats.NextGC))
	runtimeMetrics.MemStats.NumGC.Update(int64(memStats.NumGC - numGC))
	runtimeMetrics.MemStats.GCCPUFraction.Update(gcCPUFraction(&memStats))

	i := numGC % uint32(len(memStats.PauseNs))
	ii := memStats.NumGC % uint32(len(memStats.PauseNs))
	if memStats.NumGC-numGC >= uint32(len(memStats.PauseNs)) {
		for i = 0; i < uint32(len(memStats.PauseNs)); i++ {
			runtimeMetrics.MemStats.PauseNs.Update(int64(memStats.PauseNs[i]))
		}
	} else {
		if i > ii {
			for ; i < uint32(len(memStats.PauseNs)); i++ {
				runtimeMetrics.MemStats.PauseNs.Update(int64(memStats.PauseNs[i]))
			}
			i = 0
		}
		for ; i < ii; i++ {
			runtimeMetrics.MemStats.PauseNs.Update(int64(memStats.PauseNs[i]))
		}
	}
	frees = memStats.Frees
	lookups = memStats.Lookups
	mallocs = memStats.Mallocs
	numGC = memStats.NumGC

	runtimeMetrics.MemStats.PauseTotalNs.Update(int64(memStats.PauseTotalNs))
	runtimeMetrics.MemStats.StackInuse.Update(int64(memStats.StackInuse))
	runtimeMetrics.MemStats.StackSys.Update(int64(memStats.StackSys))
	runtimeMetrics.MemStats.Sys.Update(int64(memStats.Sys))
	runtimeMetrics.MemStats.TotalAlloc.Update(int64(memStats.TotalAlloc))

	currentNumCgoCalls := numCgoCall()
	runtimeMetrics.NumCgoCall.Update(currentNumCgoCalls - numCgoCalls)
	numCgoCalls = currentNumCgoCalls

	runtimeMetrics.NumGoroutine.Update(int64(runtime.NumGoroutine()))
}

// Collect runtime memory stats.
func collectRuntimeMemStats(r metrics.Collectry) {
	runtimeMetrics.MemStats.Alloc = metrics.NewGauge()
	runtimeMetrics.MemStats.BuckHashSys = metrics.NewGauge()
	runtimeMetrics.MemStats.DebugGC = metrics.NewGauge()
	runtimeMetrics.MemStats.EnableGC = metrics.NewGauge()
	runtimeMetrics.MemStats.Frees = metrics.NewGauge()
	runtimeMetrics.MemStats.HeapAlloc = metrics.NewGauge()
	runtimeMetrics.MemStats.HeapIdle = metrics.NewGauge()
	runtimeMetrics.MemStats.HeapInuse = metrics.NewGauge()
	runtimeMetrics.MemStats.HeapObjects = metrics.NewGauge()
	runtimeMetrics.MemStats.HeapReleased = metrics.NewGauge()
	runtimeMetrics.MemStats.HeapSys = metrics.NewGauge()
	runtimeMetrics.MemStats.LastGC = metrics.NewGauge()
	runtimeMetrics.MemStats.Lookups = metrics.NewGauge()
	runtimeMetrics.MemStats.Mallocs = metrics.NewGauge()
	runtimeMetrics.MemStats.MCacheInuse = metrics.NewGauge()
	runtimeMetrics.MemStats.MCacheSys = metrics.NewGauge()
	runtimeMetrics.MemStats.MSpanInuse = metrics.NewGauge()
	runtimeMetrics.MemStats.MSpanSys = metrics.NewGauge()
	runtimeMetrics.MemStats.NextGC = metrics.NewGauge()
	runtimeMetrics.MemStats.NumGC = metrics.NewGauge()
	runtimeMetrics.MemStats.GCCPUFraction = metrics.NewGaugeFloat64()
	runtimeMetrics.MemStats.PauseNs = metrics.NewHistogram(metrics.NewExpDecaySample(1028, 0.015))
	runtimeMetrics.MemStats.PauseTotalNs = metrics.NewGauge()
	runtimeMetrics.MemStats.StackInuse = metrics.NewGauge()
	runtimeMetrics.MemStats.StackSys = metrics.NewGauge()
	runtimeMetrics.MemStats.Sys = metrics.NewGauge()
	runtimeMetrics.MemStats.TotalAlloc = metrics.NewGauge()
	runtimeMetrics.NumCgoCall = metrics.NewGauge()
	runtimeMetrics.NumGoroutine = metrics.NewGauge()
	runtimeMetrics.ReadMemStats = metrics.NewHistogram(metrics.NewExpDecaySample(1028, 0.015))

	r.Collector("runtime.MemStats.Alloc", runtimeMetrics.MemStats.Alloc)
	r.Collector("runtime.MemStats.BuckHashSys", runtimeMetrics.MemStats.BuckHashSys)
	r.Collector("runtime.MemStats.DebugGC", runtimeMetrics.MemStats.DebugGC)
	r.Collector("runtime.MemStats.EnableGC", runtimeMetrics.MemStats.EnableGC)
	r.Collector("runtime.MemStats.Frees", runtimeMetrics.MemStats.Frees)
	r.Collector("runtime.MemStats.HeapAlloc", runtimeMetrics.MemStats.HeapAlloc)
	r.Collector("runtime.MemStats.HeapIdle", runtimeMetrics.MemStats.HeapIdle)
	r.Collector("runtime.MemStats.HeapInuse", runtimeMetrics.MemStats.HeapInuse)
	r.Collector("runtime.MemStats.HeapObjects", runtimeMetrics.MemStats.HeapObjects)
	r.Collector("runtime.MemStats.HeapReleased", runtimeMetrics.MemStats.HeapReleased)
	r.Collector("runtime.MemStats.HeapSys", runtimeMetrics.MemStats.HeapSys)
	r.Collector("runtime.MemStats.LastGC", runtimeMetrics.MemStats.LastGC)
	r.Collector("runtime.MemStats.Lookups", runtimeMetrics.MemStats.Lookups)
	r.Collector("runtime.MemStats.Mallocs", runtimeMetrics.MemStats.Mallocs)
	r.Collector("runtime.MemStats.MCacheInuse", runtimeMetrics.MemStats.MCacheInuse)
	r.Collector("runtime.MemStats.MCacheSys", runtimeMetrics.MemStats.MCacheSys)
	r.Collector("runtime.MemStats.MSpanInuse", runtimeMetrics.MemStats.MSpanInuse)
	r.Collector("runtime.MemStats.MSpanSys", runtimeMetrics.MemStats.MSpanSys)
	r.Collector("runtime.MemStats.NextGC", runtimeMetrics.MemStats.NextGC)
	r.Collector("runtime.MemStats.NumGC", runtimeMetrics.MemStats.NumGC)
	r.Collector("runtime.MemStats.GCCPUFraction", runtimeMetrics.MemStats.GCCPUFraction)
	r.Collector("runtime.MemStats.PauseNs", runtimeMetrics.MemStats.PauseNs)
	r.Collector("runtime.MemStats.PauseTotalNs", runtimeMetrics.MemStats.PauseTotalNs)
	r.Collector("runtime.MemStats.StackInuse", runtimeMetrics.MemStats.StackInuse)
	r.Collector("runtime.MemStats.StackSys", runtimeMetrics.MemStats.StackSys)
	r.Collector("runtime.MemStats.Sys", runtimeMetrics.MemStats.Sys)
	r.Collector("runtime.MemStats.TotalAlloc", runtimeMetrics.MemStats.TotalAlloc)
	r.Collector("runtime.NumCgoCall", runtimeMetrics.NumCgoCall)
	r.Collector("runtime.NumGoroutine", runtimeMetrics.NumGoroutine)
	r.Collector("runtime.ReadMemStats", runtimeMetrics.ReadMemStats)
}

// Cgo call
func numCgoCall() int64 {
	return 0
}

// gcCPUFraction call
func gcCPUFraction(memStats *runtime.MemStats) float64 {
	return memStats.GCCPUFraction
}
