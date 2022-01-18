package system_test

import (
	"testing"
	"time"

	"github.com/ezgroot/ezUtils/system"
)

func TestTimestampNow(t *testing.T) {
	t.Logf("TimestampNow    = %d\n", system.TimestampNow())
}

func TestTimestampMilNow(t *testing.T) {
	t.Logf("TimestampMilNow = %d\n", system.TimestampMilNow())
}

func TestTimestampMicNow(t *testing.T) {
	t.Logf("TimestampMicNow = %d\n", system.TimestampMicNow())
}

func TestTimestampNanoNow(t *testing.T) {
	t.Logf("TimestampMicNow = %d\n", system.TimestampNanoNow())
}

func TestTimestampToTime(t *testing.T) {
	t.Logf("TimestampToTime = %v\n", system.TimestampToTime(162194457))
	t.Logf("TimestampToTime = %v\n", system.TimestampToTime(1621944573))
	t.Logf("TimestampToTime = %v\n", system.TimestampToTime(1603287383074255335))
	t.Logf("TimestampToTime = %v\n", system.TimestampToTime(8603287383074255335))
}

func TestTimeDataNow(t *testing.T) {
	t.Logf("TimeDataNow = %s\n", system.TimeDataNow())
}

func TestTimeToTimestamp(t *testing.T) {
	ts, err := system.TimeToTimestamp("2020-10-21T20:03:46.123456789+08:00", time.RFC1123)
	if err == nil {
		t.Logf("TimeToTimestamp = %d\n", ts)
	} else {
		t.Fatalf("err = %s\n", err)
		return
	}

	ts, err = system.TimeToTimestamp("2020-10-21 20:03:46.123456789 +0800 CST",
		"2006-01-02 15:04:05.999999999 -0700 MST")
	if err == nil {
		t.Logf("TimeToTimestamp = %d\n", ts)
	} else {
		t.Fatalf("err = %s\n", err)
	}
}
