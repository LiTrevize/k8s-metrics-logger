package main

type MetricsLogger struct {
	kc *KubeletClient
	nc *NvmlClient
}

func NewMetricsLogger() *MetricsLogger {
	ml := &MetricsLogger{NewKubeletClient(), NewNvmlClient()}
	return ml
}

func (ml *MetricsLogger) LogMetrics() {
	ml.kc.LogMetrics()
	ml.nc.LogMetrics(ml.kc.Node)
}
