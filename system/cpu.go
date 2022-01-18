package system

import (
	"context"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
)

// CPUArch get CPU architecture type.
func CPUArch() (arch string) {
	return runtime.GOARCH
}

// CPUNum get the number of logical CPUs usable by the current process.
func CPUNum() (coreNum int) {
	return runtime.NumCPU()
}

// CPUPhysicalCount get CPU physical count.
func CPUPhysicalCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	defer cancel()

	return cpu.CountsWithContext(ctx, false)
}

// CPULogicalCount get CPU logical count.
func CPULogicalCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	defer cancel()

	return cpu.CountsWithContext(ctx, true)
}

type CPU struct {
	ModeName string  `json:"modeName"`
	Mhz      float64 `json:"mhz"`
	CoreNum  int     `json:"coreNum"`
}

// CPUInfo get CPU info.
func CPUInfo() (map[string]CPU, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	defer cancel()

	info, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}

	var cpuMap = make(map[string]CPU, 64)
	for _, v := range info {
		cpu := cpuMap[v.ModelName]

		cpu.Mhz = v.Mhz
		cpu.ModeName = v.ModelName
		cpu.CoreNum++
		cpuMap[v.ModelName] = cpu
	}

	return cpuMap, nil
}

// CPULoad get CPU load.
func CPULoad() (*load.AvgStat, error) {
	return load.Avg()
}

// CPUUsageTotal get CPU usage.
func CPUUsageTotal() (float64, error) {
	info, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}

	return info[0], nil
}

// CPUUsagePer get CPU usage.
func CPUUsagePer() ([]float64, error) {
	return cpu.Percent(time.Second, true)
}

// CPUTimesTotal get CPU usage.
func CPUTimesTotal() (cpu.TimesStat, error) {
	var info cpu.TimesStat
	list, err := cpu.Times(false)
	if err != nil {
		return info, err
	}

	info = list[0]

	return info, nil
}

// CPUTimesTotal get CPU usage.
func CPUTimesPer() ([]cpu.TimesStat, error) {
	info, err := cpu.Times(true)
	if err != nil {
		return nil, err
	}

	return info, nil
}
