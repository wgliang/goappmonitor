package goappmonitor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// 全局配置文件
type GlobalConfig struct {
	Debug    bool        `json:"debug"`
	Hostname string      `json:"hostname"`
	Tags     string      `json:"tags"`
	Step     int64       `json:"step"`
	Bases    []string    `json:"bases"`
	Push     *PushConfig `json:"push"`
	Http     *HttpConfig `json:"http"`
}

// http配置
type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

// 推送数据配置
type PushConfig struct {
	Enabled bool   `json:"enabled"`
	Api     string `json:"api"`
}

var (
	// 默认配置
	configFn     = "./cfg.json"
	defaultTags  = ""
	defaultStep  = int64(60)
	defaultBases = []string{}
	defaultPush  = &PushConfig{Enabled: true, Api: "http://127.0.0.1:1988/v1/push"}
	defaultHttp  = &HttpConfig{Enabled: false, Listen: ""}
	// 全局变量
	cfg      *GlobalConfig
	cfgLock  = new(sync.RWMutex)
	step     int64
	api      string
	gdebug   bool
	endpoint string
	gtags    string
)

// 获取配置信息
func config() *GlobalConfig {
	cfgLock.RLock()
	defer cfgLock.RUnlock()
	return cfg
}

// 加载配置文件
func loadConfig() error {
	if !isFileExist(configFn) {
		return fmt.Errorf("config file not found: %s", configFn)
	}
	// 解析配置
	c, err := parseConfig(configFn)
	if err != nil {
		return err
	}
	// 更新配置
	updateConfig(c)
	return nil
}

// 默认配置
func setDefaultConfig() {
	dcfg := defaultConfig()
	updateConfig(dcfg)
}

// 默认配置
func defaultConfig() GlobalConfig {
	return GlobalConfig{
		Debug:    false,
		Hostname: defaultHostname(),
		Tags:     defaultTags,
		Step:     defaultStep,
		Bases:    defaultBases,
		Push:     defaultPush,
		Http:     defaultHttp,
	}
}

// 更新配置
func updateConfig(c GlobalConfig) {
	nc := formatConfig(c)
	cfgLock.Lock()
	defer cfgLock.Unlock()
	cfg = &nc
}

// 格式化配置
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
	if nc.Push.Enabled && nc.Push.Api == "" {
		nc.Push = defaultPush
	}
	if len(nc.Bases) < 1 {
		nc.Bases = defaultBases
	}

	return nc
}

// 解析配置文件
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

// 默认主机名称
func defaultHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

// 文件存在性检测
func isFileExist(fn string) bool {
	_, err := os.Stat(fn)
	return err == nil || os.IsExist(err)
}

// 读取配置文件
func readFileString(fn string) (string, error) {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
