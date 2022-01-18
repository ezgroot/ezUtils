package system_test

import (
	"testing"

	"github.com/ezgroot/ezUtils/system"
)

func TestDiskInfo(t *testing.T) {
	info, err := system.DiskInfo()
	if err != nil {
		t.Fatalf("DiskInfo error = %s", err)
		return
	}

	t.Logf("DiskInfo = %+v", info)
}
