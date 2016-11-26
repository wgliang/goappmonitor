package goappmonitor

import (
	"time"

	"github.com/wgliang/metrics"
)

// All can collect Type-of-Monitoring-Data. No matter what you add into
// goappmonitor,all data will one of the types. So you should choose which type
// is your best choice.
var (
	// float64-gauge
	appGaugeFloat64 = metrics.NewCollectry()
	// counter
	appCounter = metrics.NewCollectry()
	// meter
	appMeter = metrics.NewCollectry()
	// histogram
	appHistogram = metrics.NewCollectry()
	// debug status data
	appDebug = metrics.NewCollectry()
	// runtime statusc data
	appRuntime = metrics.NewCollectry()
	// self
	appSelf = metrics.NewCollectry()
	// all collect data
	values = make(map[string]metrics.Collectry)
)

// Initialize all your type.
func init() {
	values["gauge"] = appGaugeFloat64
	values["counter"] = appCounter
	values["meter"] = appMeter
	values["histogram"] = appHistogram
	values["debug"] = appDebug
	values["runtime"] = appRuntime
	values["self"] = appSelf
}

// Return raw data of a metric.
func rawMetric(types []string) map[string]interface{} {
	data := make(map[string]interface{})
	for _, mtype := range types {
		if v, ok := values[mtype]; ok {
			data[mtype] = v.Values()
		}
	}
	return data
}

// Return all-type metrics raw data.
func rawMetrics() map[string]interface{} {
	data := make(map[string]interface{})
	for key, v := range values {
		data[key] = v.Values()
	}
	return data
}

// Retuen all-type metrics data size.
func rawSizes() map[string]int64 {
	data := map[string]int64{}
	all := int64(0)
	for key, v := range values {
		kv := v.Size()
		all += kv
		data[key] = kv
	}
	data["all"] = all
	return data
}

// Collect all base or system data. And it contains debug and runtime status.
func collectBase(bases []string) {
	// collect data after 30s
	time.Sleep(time.Duration(30) * time.Second)
	// if open debug
	if contains(bases, "debug") {
		CollectDebugGCStats(appDebug, 5e9)
	}
	// if open runtime
	if contains(bases, "runtime") {
		CollectRuntimeMemStats(appRuntime, 5e9)
	}
}

// Check base status.
func contains(bases []string, name string) bool {
	for _, n := range bases {
		if n == name {
			return true
		}
	}
	return false
}
