API
====

Several types of statistics are provided by goappmonitor, Counter、Gauge、GaugeFloat64、Meter、Histogram、Timer、Health。The meaning of statistics, see Synonyms at [java-metrics](http://metrics.dropwizard.io/3.1.0/getting-started/)。


Counter
----

A counter is just a gauge for an AtomicLong instance. You can increment or decrement its value. 

##### Counter
+ Interface: Counter(name string, value int64)
+ Parameter: name - counter name; value - initial value
+ example:

```go
Counter("counterSize", int64(0))
```

##### DecCounter
+ Interface: DecCounter(name string, value int64)
+ Parameter: name - counter name; value - initial value
+ example:

```go
DecCounter("counterSize", int64(1))
```

##### IncCounter
+ Interface: IncCounter(name string, value int64)
+ Parameter: name - counter name; value - initial value
+ example:

```go
IncCounter("counterSize", int64(2))
```

##### GetCounter
+ Interface: GetCounter(name string) int64
+ Parameter: name - counter name
+ example:

```go
counterSize := GetCounter("counterSize")
```

Gauge
----

A gauge is an instantaneous measurement of a value(int64). 

##### Gauge
+ Interface: Gauge(name string, value int64)
+ Parameter: name - gauge name; value - initial value
+ example:

```go
Gauge("queueSize", int64(100))
```

##### SetGauge
+ Interface: SetGauge(name string, value float64)
+ Parameter: name - gauge name; value - initial value
+ example:

```go
SetGauge("queueSize", float64(18))
```

##### GetGauge
+ Interface: GetGauge(name string) float64
+ Parameter: name - gauge name
+ example:

```go
queueSize := GetGauge("queueSize")
```

GaugeFloat64
----

A gauge is an instantaneous measurement of a value(float64). 

##### GaugeFloat64
+ Interface: Gauge(name string, value float64)
+ Parameter: name - gauge name; value - initial value
+ example:

```go
GaugeFloat64("queueSize", float64(100.00))
```

##### SetGauge
+ Interface: SetGauge(name string, value float64)
+ Parameter: name - gauge name; value - initial value
+ example:

```go
SetGauge("queueSize", float64(18.0))
```

##### GetGauge
+ Interface: GetGauge(name string) float64
+ Parameter: name - gauge name
+ example:

```go
queueSize := GetGauge("queueSize")
```

Meter
----

A meter measures the rate of events over time (e.g., “requests per second”). In addition to the mean rate, meters also track 1-, 5-, and 15-minute moving averages.

##### New 
+ Interface: Meter(name string, value int64)
+ Parameter: value - the number of events that occurred 	
+ example:

```go
Meter("pageView", int64(1))
```

##### SetMeter 
+ Interface: SetMeter(name string, value int64)
+ Parameter: name - meter name;value - the number of events that occurred 	
+ example:

```go
SetMeter("pageView", int64(1))
```

##### GetMeter
+ Interface: GetMeter(name string) int64
+ Parameter: name - meter name;value - the number of events that occurred 	
+ example:

```go
pvSum := GetMeter("pageView")
```

##### GetMeterRateStep
+ Interface: GetMeterRateStep(name string) float64
+ Parameter: name - meter name
+ example:

```go
pageRateStep := GetMeterRateStep("pageView")
```

##### GetMeterRateMean
+ Interface: GetMeterRateMean(name string) float64
+ Parameter: name - meter name
+ example:

```go
pageRateMean := GetMeterRateMean("pageView")
```

##### GetMeterRate1
+ Interface: GetMeterRate1(name string) float64
+ Parameter: name - meter name
+ example:

```go
pageRate1 := GetMeterRate1("pageView")
```

##### GetMeterRate5
+ Interface: GetMeterRate5(name string) float64
+ Parameter: name - meter name
+ example:

```go
pageRate5 := GetMeterRate5("pageView")
```

##### GetMeterRate15
+ Interface: GetMeterRate15(name string) float64
+ Parameter: name - meter name
+ example:

```go
pageRate15 := GetMeterRate15("pageView")
```

Histogram
----

A histogram measures the statistical distribution of values in a stream of data. In addition to minimum, maximum, mean, etc., it also measures median, 75th, 90th, 95th, 98th, 99th, and 99.9th percentiles.

##### Histogram 
+ Interface: Histogram(name string, value int64)
+ Parameter: name - histogram name;value - the value of the current sampling point is recorded.
+ example:

```go
Histogram("processNum", int64(1))
```

##### Set 
+ Interface: SetHistogram(name string, value int64)
+ Parameter: name - histogram name;value - the value of the current sampling point is recorded.
+ example:

```go
SetHistogram("processNum", int64(1))
```

##### GetHistogram
+ Interface: GetHistogram(name string) int64
+ Parameter: name - histogram name
+ example:

```go
processNum := GetHistogram("processNum")
```

##### GetHistogramMax
+ Interface: GetHistogramMax(name string) float64
+ Parameter: name - histogram name
+ example:

```go
processNumMax := GetHistogramMax("processNum")
```

##### GetHistogramMin
+ Interface: GetHistogramMin(name string) float64
+ Parameter: name - histogram name
+ example:

```go
processNumMin := GetHistogramMin("processNum")
```

##### GetHistogramSum
+ Interface: GetHistogramSum(name string) float64
+ Parameter: name - histogram name
+ example:

```go
processNumSum := GetHistogramSum("processNum")
```

##### GetHistogramMean
+ Interface: GetHistogramMean(name string) float64
+ Parameter: name - histogram name
+ example:

```go
processNumMean := GetHistogramMean("processNum")
```

##### GetHistogramStdDev
+ Interface: GetHistogramStdDev(name string) float64
+ Parameter: name - histogram name
+ example:

```go
processNumStdDev := GetHistogramStdDev("processNum")
```

##### GetHistogram50th
+ Interface: GetHistogram50th(name string) float64
+ Parameter: name - histogram name
+ example:

```go
processNum50th := GetHistogram50th("processNum")
```

##### GetHistogram75th
+ Interface: GetHistogram75th(name string) float64
+ Parameter: name - histogram name
+ example:

```go
processNum75th := GetHistogram75th("processNum")
```

##### GetHistogram95th
+ Interface: GetHistogram95th(name string) float64
+ Parameter: name - histogram name
+ example:

```go
processNum95th := GetHistogram95th("processNum")
```

##### GetHistogram99th
+ Interface: GetHistogram99th(name string) float64
+ Parameter: name - histogram name
+ example:

```go
processNum99th := GetHistogram99th("processNum")
```

##### GetHistogram999th
+ Interface: GetHistogram999th(name string) float64
+ Parameter: name - histogram name
+ example:

```go
processNum999th := GetHistogram999th("processNum")
```
