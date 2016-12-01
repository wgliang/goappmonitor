package main

import (
	"math/rand"
	"time"

	appm "github.com/wgliang/goappmonitor"
)

// Base or system performance data,such as memeory,gc,network and so on.
func baseOrsystem() {
	for _ = range time.Tick(time.Second * time.Duration(10)) {
		// (commonly used) Meter, used to sum and calculate the rate of change. Use scenarios
		// such as the number of home visits statistics, CG etc..
		pv := int64(rand.Int31n(100))
		appm.Meter("appm.meter", pv)
		appm.Meter("appm.meter.2", pv-50)

		// (commonly used) Gauge, used to preserve the value of the instantaneous value of the
		// type of record. Use scenarios such as statistical queue length, statistics CPU usage,
		// and so on.
		queueSize := int64(rand.Int31n(100) - 50)
		appm.Gauge("appm.gauge", queueSize)

		cpuUtil := float64(rand.Int31n(10000)) / float64(100)
		appm.GaugeFloat64("appm.gauge.float64", cpuUtil)
	}
}

// Custom or business performance data,such as qps,num of function be called, task queue and so on.
func customOrbusiness() {
	for _ = range time.Tick(time.Second) {
		// Histogram, using the exponential decay sampling method, the probability distribution of
		// the statistical object is calculated. Using scenarios such as the probability distribution
		// of the statistics home page to access the delay
		delay := int64(rand.Int31n(100))
		appm.Histogram("appm.histogram", delay)
	}
}

func main() {
	var ch chan int
	go baseOrsystem()
	go customOrbusiness()
	<-ch
}
