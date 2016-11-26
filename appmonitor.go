package goappmonitor

import (
	"time"

	"github.com/wgliang/metrics"
)

// 所有的应用监控指标集合
var (
	// float64类型的测量器
	appGaugeFloat64 = metrics.NewCollectry()
	// 计数器
	appCounter = metrics.NewCollectry()
	// 仪表测量器
	appMeter = metrics.NewCollectry()
	// 柱状图
	appHistogram = metrics.NewCollectry()
	// 调试状态信息
	appDebug = metrics.NewCollectry()
	// 运行状态信息
	appRuntime = metrics.NewCollectry()
	//
	appSelf = metrics.NewCollectry()
	// 只读的采集信息
	values = make(map[string]metrics.Collectry)
)

// 初始化所有的采集器
func init() {
	values["gauge"] = appGaugeFloat64
	values["counter"] = appCounter
	values["meter"] = appMeter
	values["histogram"] = appHistogram
	values["debug"] = appDebug
	values["runtime"] = appRuntime
	values["self"] = appSelf
}

// 指定类型元数据
func rawMetric(types []string) map[string]interface{} {
	data := make(map[string]interface{})
	for _, mtype := range types {
		if v, ok := values[mtype]; ok {
			data[mtype] = v.Values()
		}
	}
	return data
}

// 所有类型元数据
func rawMetrics() map[string]interface{} {
	data := make(map[string]interface{})
	for key, v := range values {
		data[key] = v.Values()
	}
	return data
}

// 数据量大小，单个类型数据量和all数据量总和
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

// 采集应用基础监控数据
func collectBase(bases []string) {
	// 30秒后开始采集数据
	time.Sleep(time.Duration(30) * time.Second)
	// 包含debug信息
	if contains(bases, "debug") {
		CollectDebugGCStats(appDebug, 5e9)
	}
	// 包含运行时信息
	if contains(bases, "runtime") {
		CollectRuntimeMemStats(appRuntime, 5e9)
	}
}

// 基础数据类型检查
func contains(bases []string, name string) bool {
	for _, n := range bases {
		if n == name {
			return true
		}
	}
	return false
}
