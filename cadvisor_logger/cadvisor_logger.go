package main

import (
	"time"

	"github.com/LiTrevize/k8s-metrics-logger/util"
)

func main() {
	kc := util.NewKubeletClient()
	for range time.Tick(time.Minute) {
		kc.LogMetrics()
	}
}
