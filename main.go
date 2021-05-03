package main

import (
	"time"
)

func main() {
	kc := NewKubeletClient()

	for range time.Tick(time.Second * 60) {
		go func() {
			kc.LogMetrics()
		}()
	}
}
