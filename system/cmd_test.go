package system_test

/*
 * go test -v cmd_test.go
 */

import (
	"testing"

	"github.com/ezgroot/ezUtils/system"
)

func TestBashCommand(t *testing.T) {
	result, err := system.BashCommand("ls -la")
	if err != nil {
		t.Fatalf("err = %s\n", err)
	} else {
		t.Logf("%s\n", result)
	}

	result, err = system.BashCommand("pwd")
	if err != nil {
		t.Fatalf("err = %s\n", err)
	} else {
		t.Logf("%s\n", result)
	}
}
