package goappmonitor

import (
	"testing"
	"time"

	"github.com/wgliang/metrics"
)

func BenchmarkDebugGCStats(b *testing.B) {
	r := metrics.NewCollectry()
	collectDebugGCStats(r)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		captureDebugGCStatsWorker(r)
	}
}

func TestCollectDebugGCStats(t *testing.T) {
	r := metrics.NewCollectry()
	CollectDebugGCStats(r, time.Duration(1))
}
