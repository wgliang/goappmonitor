package goappmonitor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// Global config of goappmonitor, and you can config it in cfg.json.
type GlobalConfig struct {
	Debug    bool        `json:"debug"`
	Hostname string      `json:"hostname"`
	Tags     string      `json:"tags"`
	Step     int64       `json:"step"`
	Bases    []string    `json:"bases"`
	Push     *PushConfig `json:"push"`
	Http     *HttpConfig `json:"http"`
}

// Http config about whether open local server and server address.
type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

// Push config of pushing address and switcher.
type PushConfig struct {
	Enabled bool   `json:"enabled"`
	Api     string `json:"api"`
}

var (
	// default config
	configFn     = "./cfg.json" // 配置文件路径
	defaultTags  = ""           // 标签
	defaultStep  = int64(60)    // 默认采集频率60s一次
	defaultBases = []string{}
	defaultPush  = &PushConfig{Enabled: true, Api: "http://127.0.0.1:1988/v1/push"}
	defaultHttp  = &HttpConfig{Enabled: false, Listen: ""}

	// global variables
	cfg      *GlobalConfig
	cfgLock  = new(sync.RWMutex)
	step     int64
	api      string
	gdebug   bool
	endpoint string
	gtags    string
)

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
		Push:     defaultPush,
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
	if nc.Push.Enabled && nc.Push.Api == "" {
		nc.Push = defaultPush
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
