package system

import (
	"github.com/shirou/gopsutil/disk"
)

type Partition struct {
	Mountpoint string
	Total      uint64  // per - GB
	Free       uint64  // per - GB
	Usage      float64 // 100%
}

// Disk info.
type Disk struct {
	Partitions []Partition
}

// DiskInfo get disk info.
func DiskInfo() (Disk, error) {
	var diskInfo Disk
	list, err := disk.Partitions(false)
	if err != nil {
		return diskInfo, err
	}

	for _, v := range list {
		var partition Partition
		partition.Mountpoint = v.Mountpoint

		usageInfo, err := disk.Usage(v.Mountpoint)
		if err == nil {
			partition.Total = usageInfo.Total / 1024 / 1024 / 1024
			partition.Free = usageInfo.Free / 1024 / 1024 / 1024
			partition.Usage = usageInfo.UsedPercent
		}

		diskInfo.Partitions = append(diskInfo.Partitions, partition)
	}

	return diskInfo, err
}
