package main

import (
	"os"
	"time"

	"github.com/LiTrevize/k8s-metrics-logger/util"
)

func main() {
	nc := util.NewNvmlClient()
	nodename := os.Getenv("NODE_NAME")
	nc.LogDeviceInfo(nodename)
	nc.LogMetrics(nodename)
	go func() {
		for range time.Tick(time.Hour) {
			nc.LogDeviceInfo(nodename)
		}
	}()
	for range time.Tick(time.Minute) {
		nc.LogMetrics(nodename)
	}
}
