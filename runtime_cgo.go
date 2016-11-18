// +build cgo
// +build !appengine

package appmonitor

import "runtime"

func numCgoCall() int64 {
	return runtime.NumCgoCall()
}
