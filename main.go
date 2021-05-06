package main

import (
	"time"
)

func main() {
	ml := NewMetricsLogger()

	ml.LogDeviceInfo()
	ml.LogMetrics()

	go func() {
		for range time.Tick(time.Second * 1) {
			ml.LogDeviceInfo()
		}
	}()
	for range time.Tick(time.Second * 5) {
		ml.LogMetrics()
	}

}
