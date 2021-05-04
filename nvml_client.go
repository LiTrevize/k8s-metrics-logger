package main

import (
	"fmt"
	"log"
	"time"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
)

type NvmlClient struct {
	DeviceCount   uint
	DriverVersion string
	Devices       []*nvml.Device
}

func NewNvmlClient() *NvmlClient {
	nc := &NvmlClient{}
	err := nc.Init()
	if err != nil {
		return nil
	}
	return nc
}

func (nc *NvmlClient) Init() error {
	nvml.Init()

	count, err := nvml.GetDeviceCount()
	if err != nil {
		return err
	}
	nc.DeviceCount = count

	driverVersion, err := nvml.GetDriverVersion()
	if err != nil {
		return err
	}
	nc.DriverVersion = driverVersion

	for i := uint(0); i < count; i++ {
		device, err := nvml.NewDevice(i)
		if err != nil {
			return err
		}
		nc.Devices = append(nc.Devices, device)
	}

	return nil
}

func (nc *NvmlClient) LogMetrics(node string) {
	// device info
	for i, device := range nc.Devices {
		tag := map[string]string{
			"node":  node,
			"GPU":   fmt.Sprint(i),
			"UUID":  device.UUID,
			"model": *device.Model}
		field := map[string]interface{}{
			"power":        fmt.Sprintf("%d W", *device.Power),
			"memory":       fmt.Sprintf("%d MiB", *device.Memory),
			"bandwidth":    fmt.Sprintf("%d MB/s", *device.PCI.Bandwidth),
			"clock_cores":  fmt.Sprintf("%d MHz", *device.Clocks.Cores),
			"clock_memory": fmt.Sprintf("%d MHz", *device.Clocks.Memory),
		}
		(&MetricsLog{"gpu_device_info", tag, nil, time.Now(), field}).Log()
	}
}

func (nc *NvmlClient) GetDeviceStatus() {
	for i, device := range nc.Devices {
		st, err := device.Status()
		if err != nil {
			log.Panicf("Error getting device %d status: %v\n", i, err)
		}
		fmt.Printf("%5d %5d %5d %5d %5d %5d %5d %5d %5d\n",
			i, *st.Power, *st.Temperature, *st.Utilization.GPU, *st.Utilization.Memory,
			*st.Utilization.Encoder, *st.Utilization.Decoder, *st.Clocks.Memory, *st.Clocks.Cores)
	}
}
