package goappmonitor

import (
	"testing"
	"time"
)

func TestPusher2Falcon(t *testing.T) {
	go push2Falcon()
	time.Sleep(10 * time.Second)
}

func TestPusher2InfluxDB(t *testing.T) {
	go push2InfluxDB()
	time.Sleep(10 * time.Second)
}
