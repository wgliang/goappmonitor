package goappmonitor

import (
	"github.com/wgliang/metrics"
)

// RuntimeMetrics ,Monitoring index of application in production environment.
type RuntimeMetrics struct {
	MemStats struct {
		Alloc         metrics.Gauge        `json:"alloc"`
		BuckHashSys   metrics.Gauge        `json:"buckHashSys"`
		DebugGC       metrics.Gauge        `json:"debugGC"`
		EnableGC      metrics.Gauge        `json:"enableGC"`
		Frees         metrics.Gauge        `json:"frees"`
		HeapAlloc     metrics.Gauge        `json:"heapAlloc"`
		HeapIdle      metrics.Gauge        `json:"heapIdle"`
		HeapInuse     metrics.Gauge        `json:"heapInuse"`
		HeapObjects   metrics.Gauge        `json:"heapObjects"`
		HeapReleased  metrics.Gauge        `json:"heapReleased"`
		HeapSys       metrics.Gauge        `json:"heapSys"`
		LastGC        metrics.Gauge        `json:"lastGC"`
		Lookups       metrics.Gauge        `json:"lookups"`
		Mallocs       metrics.Gauge        `json:"mallocs"`
		MCacheInuse   metrics.Gauge        `json:"mCacheInuse"`
		MCacheSys     metrics.Gauge        `json:"mCacheSys"`
		MSpanInuse    metrics.Gauge        `json:"mSpanInuse"`
		MSpanSys      metrics.Gauge        `json:"mSpanSys"`
		NextGC        metrics.Gauge        `json:"nextGC"`
		NumGC         metrics.Gauge        `json:"numGC"`
		GCCPUFraction metrics.GaugeFloat64 `json:"gCCPUFraction"`
		PauseNs       metrics.Histogram    `json:"pauseNs"`
		PauseTotalNs  metrics.Gauge        `json:"pauseTotalNs"`
		StackInuse    metrics.Gauge        `json:"stackInuse"`
		StackSys      metrics.Gauge        `json:"stackSys"`
		Sys           metrics.Gauge        `json:"sys"`
		TotalAlloc    metrics.Gauge        `json:"totalAlloc"`
	}
	NumCgoCall   metrics.Gauge     `json:"numCgoCall"`
	NumGoroutine metrics.Gauge     `json:"numGoroutine"`
	ReadMemStats metrics.Histogram `json:"readMemStats"`
}

// DebugMetrics ,Monitoring index of application in debug environment.
type DebugMetrics struct {
	GCStats struct {
		LastGC     metrics.Gauge     `json:"lastGC"`
		NumGC      metrics.Gauge     `json:"numGC"`
		Pause      metrics.Histogram `json:"pause"`
		PauseTotal metrics.Gauge     `json:"pauseTotal"`
	}
	ReadGCStats metrics.Histogram `json:"readGCStats"`
}
