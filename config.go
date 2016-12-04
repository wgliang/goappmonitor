package goappmonitor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/wgliang/metrics"
)

// All can collect Type-of-Monitoring-Data. No matter what you add into
// goappmonitor,all data will one of the types. So you should choose which type
// is your best choice.
var (
	// float64-gauge
	appGaugeFloat64 = metrics.NewCollectry()
	// counter
	appCounter = metrics.NewCollectry()
	// meter
	appMeter = metrics.NewCollectry()
	// histogram
	appHistogram = metrics.NewCollectry()
	// debug status data
	appDebug = metrics.NewCollectry()
	// runtime statusc data
	appRuntime = metrics.NewCollectry()
	// self
	appSelf = metrics.NewCollectry()
	// all collect data
	values = make(map[string]metrics.Collectry)

	// default config
	configFn              = "./cfg.json" // 配置文件路径
	defaultTags           = ""           // 标签
	defaultStep           = int64(60)    // 默认采集频率60s一次
	defaultBases          = []string{}
	defaultPushOpenFalcon = OpenFalconx{Enabled: true, Api: "http://127.0.0.1:1988/v1/push"}
	defaultPushInfluxDB   = InfluxDBx{Enabled: true, Addr: "http://127.0.0.1:8086", Username: "root", Password: "root"}
	defaultHttp           = &HttpConfig{Enabled: false, Listen: ""}

	// global variables
	cfg              *GlobalConfig
	cfgLock          = new(sync.RWMutex)
	step             int64
	api              string
	gdebug           bool
	endpoint         string
	gtags            string
	influxDBAddr     string
	influxDBUsername string
	influxDBPassword string
)

// GlobalConfig of goappmonitor, and you can config it in cfg.json.
type GlobalConfig struct {
	Debug    bool        `json:"debug"`
	Hostname string      `json:"hostname"`
	Tags     string      `json:"tags"`
	Step     int64       `json:"step"`
	Bases    []string    `json:"bases"`
	Push     *PushConfig `json:"push"`
	Http     *HttpConfig `json:"http"`
}

// HttpConfig about whether open local server and server address.
type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type OpenFalconx struct {
	Enabled bool   `json:"enabled"`
	Api     string `json:"api"`
}

type InfluxDBx struct {
	Enabled  bool   `json:"enabled"`
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// PushConfig of pushing address and switcher.
type PushConfig struct {
	OpenFalcon OpenFalconx `json:"open-falcon"`
	InfluxDB   InfluxDBx   `json:"influxDB"`
}

// Initialize all your type.
func init() {
	values["gauge"] = appGaugeFloat64
	values["counter"] = appCounter
	values["meter"] = appMeter
	values["histogram"] = appHistogram
	values["debug"] = appDebug
	values["runtime"] = appRuntime
	values["self"] = appSelf
}

// Get config.
func config() *GlobalConfig {
	cfgLock.RLock()
	defer cfgLock.RUnlock()
	return cfg
}

// Load config form cfg.json.
func loadConfig() error {
	if !isFileExist(configFn) {
		return fmt.Errorf("config file not found: %s", configFn)
	}
	// parse config json file.
	c, err := parseConfig(configFn)
	if err != nil {
		return err
	}
	// update config
	updateConfig(c)
	return nil
}

// Set default config.
func setDefaultConfig() {
	dcfg := defaultConfig()
	updateConfig(dcfg)
}

// Get default config.
func defaultConfig() GlobalConfig {
	return GlobalConfig{
		Debug:    false,
		Hostname: defaultHostname(),
		Tags:     defaultTags,
		Step:     defaultStep,
		Bases:    defaultBases,
		Push:     &PushConfig{defaultPushOpenFalcon, defaultPushInfluxDB},
		Http:     defaultHttp,
	}
}

// Uodate config.
func updateConfig(c GlobalConfig) {
	nc := formatConfig(c)
	cfgLock.Lock()
	defer cfgLock.Unlock()
	cfg = &nc
}

// Format config.
func formatConfig(c GlobalConfig) GlobalConfig {
	nc := c
	if nc.Hostname == "" {
		nc.Hostname = defaultHostname()
	}
	if nc.Step < 1 {
		nc.Step = defaultStep
	}
	if nc.Tags != "" {
		tagsOk := true
		tagsSlice := strings.Split(nc.Tags, ",")
		for _, tag := range tagsSlice {
			kv := strings.Split(tag, "=")
			if len(kv) != 2 || kv[0] == "name" { // name是保留tag
				tagsOk = false
				break
			}
		}
		if !tagsOk {
			nc.Tags = defaultTags
		}
	}
	if nc.Push.OpenFalcon.Enabled && nc.Push.OpenFalcon.Api == "" {
		nc.Push.OpenFalcon = defaultPushOpenFalcon
	}

	if nc.Push.InfluxDB.Enabled && nc.Push.InfluxDB.Addr == "" {
		nc.Push.InfluxDB = defaultPushInfluxDB
	}

	if len(nc.Bases) < 1 {
		nc.Bases = defaultBases
	}

	return nc
}

// Parse config.
func parseConfig(cfg string) (GlobalConfig, error) {
	var c GlobalConfig

	if cfg == "" {
		return c, fmt.Errorf("config file not found")
	}

	configContent, err := readFileString(cfg)
	if err != nil {
		return c, fmt.Errorf("read config file %s error: %v", cfg, err.Error())
	}

	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return c, fmt.Errorf("parse config file %s error: %v", cfg, err.Error())
	}
	return c, nil
}

// Default host name.
func defaultHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

// Whether file is exist.
func isFileExist(fn string) bool {
	_, err := os.Stat(fn)
	return err == nil || os.IsExist(err)
}

// Read config file.
func readFileString(fn string) (string, error) {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

// Return raw data of a metric.
func rawMetric(types []string) map[string]interface{} {
	data := make(map[string]interface{})
	for _, mtype := range types {
		if v, ok := values[mtype]; ok {
			data[mtype] = v.Values()
		}
	}
	return data
}

// Return all-type metrics raw data.
func rawMetrics() map[string]interface{} {
	data := make(map[string]interface{})
	for key, v := range values {
		data[key] = v.Values()
	}
	return data
}

// Retuen all-type metrics data size.
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

// Collect all base or system data. And it contains debug and runtime status.
func collectBase(bases []string) {
	// collect data after 30s
	time.Sleep(time.Duration(30) * time.Second)
	// if open debug
	if contains(bases, "debug") {
		CollectDebugGCStats(appDebug, 5e9)
	}
	// if open runtime
	if contains(bases, "runtime") {
		CollectRuntimeMemStats(appRuntime, 5e9)
	}
}

// Check base status.
func contains(bases []string, name string) bool {
	for _, n := range bases {
		if n == name {
			return true
		}
	}
	return false
}
