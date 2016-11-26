package main

import (
	"math/rand"
	"time"

	appm "github.com/wgliang/goappmonitor"
)

func main() {
	go basic()  // 基础统计器
	go senior() // 高级统计器
	select {}
}

func basic() {
	for _ = range time.Tick(time.Second * time.Duration(10)) {
		// (常用) Meter,用于累加求和、计算变化率。使用场景如，统计首页访问次数、gvm的CG次数等。
		pv := int64(rand.Int() % 100)
		appm.Meter("test.meter", pv)
		appm.Meter("test.meter.2", pv-50)

		// (常用) Gauge,用于保存数值类型的瞬时记录值。使用场景如，统计队列长度、统计CPU使用率等
		queueSize := int64(rand.Int()%100 - 50)
		appm.Gauge("test.gauge", queueSize)

		cpuUtil := float64(rand.Int()%10000) / float64(100)
		appm.GaugeFloat64("test.gauge.float64", cpuUtil)
	}
}

func senior() {
	for _ = range time.Tick(time.Second) {
		// Histogram,使用指数衰减抽样的方式，计算被统计对象的概率分布情况。使用场景如，统计主页访问延时的概率分布
		delay := int64(rand.Int() % 100)
		appm.Histogram("test.histogram", delay)
	}
}
