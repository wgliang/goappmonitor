package goappmonitor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // pprof collector
	"strings"
)

// Start Http serever.
func startHttp(addr string, debug bool) {
	configCommonRoutes()
	configProcRoutes()
	if len(addr) >= 9 {
		s := &http.Server{
			Addr:           addr,
			MaxHeaderBytes: 1 << 30,
		}
		go func() {
			if debug {
				log.Println("[goappmonitor] http server start, listening on", addr)
			}
			s.ListenAndServe()
			if debug {
				log.Println("[goappmonitor] http server stop,", addr)
			}
		}()
	}
}

// Routers config.
func configProcRoutes() {
	http.HandleFunc("/pfc/proc/metrics/json", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		RenderJson(w, rawMetrics())
	})
	http.HandleFunc("/pfc/proc/metrics/falcon", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		RenderJson(w, falconMetrics())
	})
	http.HandleFunc("/pfc/proc/metrics/influxdb", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		InfluxDBJson(w, influxDBMetrics())
	})
	// url=/pfc/proc/metric/{json,falcon}
	http.HandleFunc("/pfc/proc/metrics/", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		urlParam := r.URL.Path[len("/pfc/proc/metrics/"):]
		args := strings.Split(urlParam, "/")
		argsLen := len(args)
		if argsLen != 2 {
			RenderJson(w, "")
			return
		}

		types := []string{}
		typeslice := strings.Split(args[0], ",")
		for _, t := range typeslice {
			nt := strings.TrimSpace(t)
			if nt != "" {
				types = append(types, nt)
			}
		}

		if args[1] == "json" {
			RenderJson(w, rawMetric(types))
			return
		}
		if args[1] == "falcon" {
			RenderJson(w, falconMetric(types))
			return
		}

		if args[1] == "influxdb" {
			InfluxDBJson(w, influxDBMetric(types))
			return
		}

	})

	http.HandleFunc("/pfc/proc/metrics/size", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		RenderJson(w, rawSizes())
	})

}

// Common router config.
func configCommonRoutes() {
	http.HandleFunc("/pfc/health", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/pfc/version", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		w.Write([]byte(fmt.Sprintf("%s\n", VERSION)))
	})

	http.HandleFunc("/pfc/config", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		RenderJson(w, config())
	})

	http.HandleFunc("/pfc/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		loadConfig()
		RenderJson(w, "ok")
	})
}

func isLocalReq(raddr string) bool {
	return strings.HasPrefix(raddr, "127.0.0.1")
}

// RenderJson json
func RenderJson(w http.ResponseWriter, data interface{}) {
	renderJson(w, Response{Msg: "success", Data: data})
}

// RenderString string
func RenderString(w http.ResponseWriter, msg string) {
	renderJson(w, map[string]string{"msg": msg})
}

func renderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

// Response http struct
type Response struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Render json
func InfluxDBJson(w http.ResponseWriter, data interface{}) {
	influxDBJson(w, InfluxDBResponse{
		"success",
		endpoint,
		[]string{
			"metric",
			"counterType",
			"tags",
			"value",
			"step",
			"timestamp",
		},
		data,
	})
}

// Render string
func InfluxDBString(w http.ResponseWriter, msg string) {
	influxDBJson(w, map[string]string{"msg": msg})
}

func influxDBJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

type InfluxDBResponse struct {
	Msg     string      `json:"msg"`
	Name    string      `json:"name"`
	Columns []string    `json:"columns"`
	Points  interface{} `json:"points"`
}
