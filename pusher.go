package appmonitor

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	bhttp "github.com/niean/gotools/http/httpclient/beego"
	"github.com/wgliang/metrics"
)

const (
	GAUGE = "GAUGE"
)

func init() {
	// 读取配置文件
	cfg = config()
	// 数据采集频率
	step = cfg.Step
	// 推送地址
	api = cfg.Push.Api
	// 采集数据状态，调试还是生产环境
	gdebug = cfg.Debug
	// 主机名
	endpoint = cfg.Hostname
	// 标签
	gtags = cfg.Tags
}

// 推送到open－falcon的agent上
func push2Falcon() {

	// 准备下然后开始工作
	alignPushStartTs(step)
	// 定一个闹钟
	ti := time.Tick(time.Duration(step) * time.Second)
	for {
		select {
		case <-ti:
			// 采集事件计数
			selfMeter("pfc.push.cnt", 1) // statistics
			// 当前采集的所有数据指标
			fms := falconMetrics()
			// 获取当前时间
			start := time.Now()
			// 推送出去
			fmt.Println(fms)
			err := push(fms, api, gdebug)
			// 推送耗时
			selfGauge("pfc.push.ms", int64(time.Since(start)/time.Millisecond)) // statistics

			if err != nil {
				if gdebug {
					log.Printf("[perfcounter] send to %s error: %v", api, err)
				}
				// 失败情况下，推送数据大小为0
				selfGauge("pfc.push.size", int64(0)) // statistics
			} else {
				// 推送数据大小
				selfGauge("pfc.push.size", int64(len(fms))) // statistics
			}
		}
	}
}

// open-falcon类型metric
func falconMetric(types []string) (fd []*MetricValue) {
	for _, ty := range types {
		if r, ok := values[ty]; ok && r != nil {
			data := _falconMetric(r)
			fd = append(fd, data...)
		}
	}
	return fd
}

// open-falcon类型metrics
func falconMetrics() []*MetricValue {
	data := make([]*MetricValue, 0)
	for _, r := range values {
		nd := _falconMetric(r)
		data = append(data, nd...)
	}
	return data
}

// open-falcon 转换internal
func _falconMetric(r metrics.Collectry) []*MetricValue {
	ts := time.Now().Unix()
	data := make([]*MetricValue, 0)
	r.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Gauge:
			m := gaugeMetricValue(metric, name, endpoint, gtags, step, ts)
			data = append(data, m...)
		case metrics.GaugeFloat64:
			m := gaugeFloat64MetricValue(metric, name, endpoint, gtags, step, ts)
			data = append(data, m...)
		case metrics.Counter:
			m := counterMetricValue(metric, name, endpoint, gtags, step, ts)
			data = append(data, m...)
		case metrics.Meter:
			// m := metric.Snapshot()
			ms := meterMetricValue(metric, name, endpoint, gtags, step, ts)
			data = append(data, ms...)
		case metrics.Histogram:
			// h := metric.Snapshot()
			ms := histogramMetricValue(metric, name, endpoint, gtags, step, ts)
			data = append(data, ms...)
		}
	})

	return data
}

// gauge类型数据转换
func gaugeMetricValue(metric metrics.Gauge, metricName, endpoint, oldtags string, step, ts int64) []*MetricValue {
	tags := getTags(metricName, oldtags)
	c := newMetricValue(endpoint, "value", metric.Value(), step, GAUGE, tags, ts)
	return []*MetricValue{c}
}

// gauge64类型数据转换
func gaugeFloat64MetricValue(metric metrics.GaugeFloat64, metricName, endpoint, oldtags string, step, ts int64) []*MetricValue {
	tags := getTags(metricName, oldtags)
	c := newMetricValue(endpoint, "value", metric.Value(), step, GAUGE, tags, ts)
	return []*MetricValue{c}
}

// counter类型数据转换
func counterMetricValue(metric metrics.Counter, metricName, endpoint, oldtags string, step, ts int64) []*MetricValue {
	tags := getTags(metricName, oldtags)
	c1 := newMetricValue(endpoint, "count", metric.Count(), step, GAUGE, tags, ts)
	return []*MetricValue{c1}
}

// meter类型数据转换
func meterMetricValue(metric metrics.Meter, metricName, endpoint, oldtags string, step, ts int64) []*MetricValue {
	data := make([]*MetricValue, 0)
	tags := getTags(metricName, oldtags)

	c1 := newMetricValue(endpoint, "rate", metric.RateMean(), step, GAUGE, tags, ts)
	c2 := newMetricValue(endpoint, "sum", metric.Count(), step, GAUGE, tags, ts)
	data = append(data, c1, c2)

	return data
}

// histogram数据类型转换
func histogramMetricValue(metric metrics.Histogram, metricName, endpoint, oldtags string, step, ts int64) []*MetricValue {
	data := make([]*MetricValue, 0)
	tags := getTags(metricName, oldtags)

	values := make(map[string]interface{})
	ps := metric.Percentiles([]float64{0.75, 0.95, 0.99})
	values["min"] = metric.Min()
	values["max"] = metric.Max()
	values["mean"] = metric.Mean()
	values["75th"] = ps[0]
	values["95th"] = ps[1]
	values["99th"] = ps[2]
	for key, val := range values {
		c := newMetricValue(endpoint, key, val, step, GAUGE, tags, ts)
		data = append(data, c)
	}

	return data
}

// 创建出metric类型数据
func newMetricValue(endpoint, metric string, value interface{}, step int64, t, tags string, ts int64) *MetricValue {
	return &MetricValue{
		Endpoint:  endpoint,
		Metric:    metric,
		Value:     value,
		Step:      step,
		Type:      t,
		Tags:      tags,
		Timestamp: ts,
	}
}

// 获取标签
func getTags(name string, tags string) string {
	if tags == "" {
		return fmt.Sprintf("name=%s", name)
	}
	return fmt.Sprintf("%s,name=%s", tags, name)
}

// 推送到代理地址
func push(data []*MetricValue, url string, debug bool) error {
	dlen := len(data)
	pkg := 200 //send pkg items once
	sent := 0
	for {
		if sent >= dlen {
			break
		}

		end := sent + pkg
		if end > dlen {
			end = dlen
		}

		pkgData := data[sent:end]
		jr, err := json.Marshal(pkgData)
		if err != nil {
			return err
		}

		response, err := bhttp.Post(url).Body(jr).String()
		if err != nil {
			return err
		}
		sent = end

		if debug {
			log.Printf("[perfcounter] push result: %v, data: %v\n", response, pkgData)
		}
	}
	return nil
}

// 凑整
func alignPushStartTs(stepSec int64) {
	nw := time.Duration(time.Now().UnixNano())
	step := time.Duration(stepSec) * time.Second
	sleepNano := step - nw%step
	if sleepNano > 0 {
		time.Sleep(sleepNano)
	}
}

// 采集数据
type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Step      int64       `json:"step"`
	Type      string      `json:"counterType"`
	Tags      string      `json:"tags"`
	Timestamp int64       `json:"timestamp"`
}

// 转化成格式化字符串
func (this *MetricValue) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Metric:%s, Tags:%s, Type:%s, Step:%d, Timestamp:%d, Value:%v>",
		this.Endpoint,
		this.Metric,
		this.Tags,
		this.Type,
		this.Step,
		this.Timestamp,
		this.Value,
	)
}
