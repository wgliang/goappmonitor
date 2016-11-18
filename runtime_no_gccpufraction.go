// +build !go1.5

package appmonitor

import "runtime"

func gcCPUFraction(memStats *runtime.MemStats) float64 {
	return 0
}
