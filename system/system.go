package system

import (
	"context"
	"runtime"

	"github.com/shirou/gopsutil/host"
)

// System get System info.
func System() (*host.InfoStat, error) {
	return host.InfoWithContext(context.Background())
}

// GetOS get os type.
func GetOS() (os string) {
	return runtime.GOOS
}
