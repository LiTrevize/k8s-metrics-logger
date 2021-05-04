package main

import (
	"time"
)

func main() {
	ml := NewMetricsLogger()

	for range time.Tick(time.Second * 5) {
		go func() {
			ml.LogMetrics()
		}()
	}
}
