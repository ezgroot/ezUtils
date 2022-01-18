package system

import "github.com/shirou/gopsutil/mem"

// Memory info.
type Memory struct {
	Total float64 `json:"total"`
	Free  float64 `json:"free"`
	Usage float64 `json:"usage"`
	Used  float64 `json:"used"`
}

// MemInfo get memory info, per - GB.
func MemInfo() (Memory, error) {
	var memInfo Memory

	temp, err := mem.VirtualMemory()
	if err == nil {
		memInfo.Total = float64(temp.Total) / 1024 / 1024 / 1024
		memInfo.Free = float64(temp.Available) / 1024 / 1024 / 1024
		memInfo.Usage = float64(temp.Used) / 1024 / 1024 / 1024
		memInfo.Used = temp.UsedPercent
	}

	return memInfo, err
}
