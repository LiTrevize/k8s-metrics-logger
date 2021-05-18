package main

import (
	"os"
	"time"

	"github.com/LiTrevize/k8s-metrics-logger/util"
)

func main() {
	nc := util.NewNvmlClient()
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	nc.LogDeviceInfo(hostname)
	nc.LogMetrics(hostname)
	go func() {
		for range time.Tick(time.Hour) {
			nc.LogDeviceInfo(hostname)
		}
	}()
	for range time.Tick(time.Minute) {
		nc.LogMetrics(hostname)
	}
}
