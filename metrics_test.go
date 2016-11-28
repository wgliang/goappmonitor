package goappmonitor

import (
	"testing"
)

func TestGauge(t *testing.T) {
	Gauge("testint64", 101)
	if v := GetGauge("testint64"); v-101 > 0.00001 || v-101 < -0.00001 {
		t.Errorf("expect value:0.0 but response:%f", v)
	}
	Gauge("testfloat64", 101.0)
	if v := GetGauge("testfloat64"); v-101 > 0.00001 || v-101 < -0.00001 {
		t.Errorf("expect value:101 but response:%f", v)
	}
}

func TestMeter(t *testing.T) {
	Meter("testmeter", 101)
	if v := GetMeter("testmeter"); v != 101 {
		t.Errorf("expect value:101 but response:%d", v)
	}

	if v := GetMeterRateStep("testmeter"); v-0 < 0.00001 && v-0 > -0.00001 {
		t.Errorf("expect value:0.0 but response:%f", v)
	}

	if v := GetMeterRateMean("testmeter"); v-0.0 < 0.00001 && v-0.0 > -0.00001 {
		t.Errorf("expect value:0.0 but response:%f", v)
	}

	if v := GetMeterRate1("testmeter"); v-0.0 > 0.00001 || v-0.0 < -0.00001 {
		t.Errorf("expect value:0.0 but response:%f", v)
	}

	if v := GetMeterRate5("testmeter"); v-0.0 > 0.00001 || v-0.0 < -0.00001 {
		t.Errorf("expect value:0.0 but response:%f", v)
	}

	if v := GetMeterRate15("testmeter"); v-0.0 > 0.00001 || v-0.0 < -0.00001 {
		t.Errorf("expect value:0.0 but response:%f", v)
	}
}

func TestHistogram(t *testing.T) {
	Histogram("testhistogram", 100)
	Histogram("testhistogram", 101)
	Histogram("testhistogram", 99)
	if v := GetHistogram("testhistogram"); v != 3 {
		t.Errorf("expect value:3 but response:%d", v)
	}

	if v := GetHistogramMax("testhistogram"); v != 101 {
		t.Errorf("expect value:101 but response:%d", v)
	}

	if v := GetHistogramMin("testhistogram"); v != 99 {
		t.Errorf("expect value:99 but response:%d", v)
	}

	if v := GetHistogramSum("testhistogram"); v != 300 {
		t.Errorf("expect value:300 but response:%d", v)
	}

	if v := GetHistogramMean("testhistogram"); v-0 < 0.00001 && v-0 > -0.00001 {
		t.Errorf("expect value:0.0 but response:%f", v)
	}

	if v := GetHistogramStdDev("testhistogram"); v-0.0 < 0.00001 && v-0.0 > -0.00001 {
		t.Errorf("expect value:0.0 but response:%f", v)
	}

	if v := GetHistogram50th("testhistogram"); v-100 > 0.00001 || v-100 < -0.00001 {
		t.Errorf("expect value:100.0 but response:%f", v)
	}

	if v := GetHistogram75th("testhistogram"); v-101 > 0.00001 || v-101 < -0.00001 {
		t.Errorf("expect value:101.0 but response:%f", v)
	}

	if v := GetHistogram99th("testhistogram"); v-101 > 0.00001 || v-101 < -0.00001 {
		t.Errorf("expect value:0.0 but response:%f", v)
	}

	if v := GetHistogram999th("testhistogram"); v-101 > 0.00001 || v-101 < -0.00001 {
		t.Errorf("expect value:101 but response:%f", v)
	}
}

func TestCounter(t *testing.T) {
	Counter("testhistogram", 100)
	Counter("testhistogram", 101)
	if v := GetCounter("testhistogram"); v != 201 {
		t.Errorf("expect value:201 but response:%d", v)
	}
}

func TestSelf(t *testing.T) {
	Gauge("tesself1", 100)
	selfGauge("tesself1", 100)
	if v := GetGauge("tesself1"); v-0 > 0.00001 && v-0 < -0.00001 {
		t.Errorf("expect value:0.0 but response:%f", v)
	}

	Meter("tesself1", 101)
	selfMeter("tesself2", 101)
	if v := GetMeter("tesself2"); v != 0 {
		t.Errorf("expect value:0 but response:%d", v)
	}
}
