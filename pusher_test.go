package goappmonitor

import (
	"testing"
	"time"
)

func TestPusher(t *testing.T) {
	go push2Falcon()
	time.Sleep(10 * time.Second)
}
