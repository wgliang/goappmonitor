package goappmonitor

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var Urls []string

func TestRenderJson(t *testing.T) {
	Urls = []string{"/pfc/proc/metrics/json",
		"/pfc/proc/metrics/falcon",
		"/pfc/proc/metrics/influxdb",
		"/pfc/proc/metrics/",
		"/pfc/proc/metrics/size",
		"/pfc/health",
		"/pfc/version",
		"/pfc/config",
		"/pfc/config/reload"}
	go func() {
		for _, v := range Urls {
			resp, err := http.Get("http://127.0.0.1:2015" + v)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
		}
	}()
	time.Sleep(time.Second * 10)
}
