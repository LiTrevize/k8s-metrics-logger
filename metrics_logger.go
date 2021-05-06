package main

type MetricsLogger struct {
	kc   *KubeletClient
	nc   *NvmlClient
	Node string
}

func NewMetricsLogger() *MetricsLogger {
	ml := &MetricsLogger{kc: NewKubeletClient(), nc: NewNvmlClient()}
	ml.Node = ml.kc.Node
	return ml
}

func (ml *MetricsLogger) LogMetrics() {
	ml.kc.LogMetrics()
	ml.nc.LogMetrics(ml.Node)
}

func (ml *MetricsLogger) LogDeviceInfo() {
	ml.nc.LogDeviceInfo(ml.Node)
}
