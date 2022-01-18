package system_test

import (
	"testing"

	"github.com/ezgroot/ezUtils/system"
)

func TestCPUArch(t *testing.T) {
	t.Logf("CPUArch    = %s\n", system.CPUArch())
}

func TestCPUNum(t *testing.T) {
	t.Logf("CPUNum    = %d\n", system.CPUNum())
}

func TestCPUPhysicalCount(t *testing.T) {
	count, err := system.CPUPhysicalCount()
	if err != nil {
		t.Fatalf("err = %s\n", err)
	} else {
		t.Logf("CPUPhysicalCount    = %d\n", count)
	}
}

func TestCPULogicalCount(t *testing.T) {
	count, err := system.CPULogicalCount()
	if err != nil {
		t.Fatalf("err = %s\n", err)
	} else {
		t.Logf("CPULogicalCount    = %d\n", count)
	}
}

func TestCPUInfo(t *testing.T) {
	info, err := system.CPUInfo()
	if err != nil {
		t.Fatalf("err = %s\n", err)
	} else {
		t.Logf("CPUInfo    = %+v\n", info)
	}
}

func TestCPULoad(t *testing.T) {
	load, err := system.CPULoad()
	if err != nil {
		t.Fatalf("err = %s\n", err)
	} else {
		t.Logf("CPULoad    = %+v\n", load)
	}
}

func TestCPUUsageTotal(t *testing.T) {
	usage, err := system.CPUUsageTotal()
	if err != nil {
		t.Fatalf("err = %s\n", err)
	} else {
		t.Logf("CPUUsageTotal    = %+v\n", usage)
	}
}

func TestCPUUsagePer(t *testing.T) {
	usage, err := system.CPUUsagePer()
	if err != nil {
		t.Fatalf("err = %s\n", err)
	} else {
		t.Logf("CPUUsagePer    = %+v\n", usage)
	}
}

func TestCPUTimesTotal(t *testing.T) {
	times, err := system.CPUTimesTotal()
	if err != nil {
		t.Fatalf("err = %s\n", err)
	} else {
		t.Logf("CPUTimesTotal    = %+v\n", times)
	}
}

func TestCPUTimesPer(t *testing.T) {
	times, err := system.CPUTimesPer()
	if err != nil {
		t.Fatalf("err = %s\n", err)
	} else {
		t.Logf("CPUTimesPer    = %+v\n", times)
	}
}
