package system

import (
	"strconv"
	"time"
)

// TimestampNow Current timestamp, accurate to the second.
func TimestampNow() int64 {
	return time.Now().Unix()
}

// TimestampMilNow Current timestamp, accurate to the millisecond.
func TimestampMilNow() int64 {
	return time.Now().UnixNano() / (1000 * 1000)
}

// TimestampMicNow Current timestamp, accurate to the micsecond.
func TimestampMicNow() int64 {
	return time.Now().UnixNano() / 1000
}

// TimestampNanoNow Current timestamp, accurate to nanosecond.
func TimestampNanoNow() int64 {
	return time.Now().UnixNano()
}

// TimestampToTime converts a timestamp nano to time.Time.
func TimestampToTime(nanoSec int64) time.Time {
	var sec int64
	var nsec int64

	timestamp := strconv.FormatInt(nanoSec, 10)
	length := len(timestamp)
	if length <= 10 {
		sec = nanoSec
		nsec = 0
	} else {
		sec = nanoSec / (1000 * 1000 * 1000)

		var timeTmp = string(timestamp[10:length])
		nsec, _ = strconv.ParseInt(timeTmp, 10, 64)
	}

	return time.Unix(sec, nsec)
}

// TimeDataNow The current date and time.
func TimeDataNow() string {
	t := time.Now()
	layout := "2006-01-02 15:04:05.999999999 -0700 MST"
	return t.Format(layout)
}

// TimeToTimestamp Convert a datetime string into a timestamp by layout.
// Layout see in go/src/time/format.go, eq:time.RFC3339Nano.
func TimeToTimestamp(timeStr string, layout string) (int64, error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return -1, err
	}

	if layout == "" {
		layout = "2006-01-02 15:04:05.999999999 -0700 MST"
	}

	times, err := time.ParseInLocation(layout, timeStr, loc)
	if err != nil {
		return -1, err
	}

	return times.UnixNano(), nil
}

func GetNowGMTDatetime() (string, error) {
	loc, err := time.LoadLocation("GMT")
	if err != nil {
		return "", err
	}

	return time.Now().In(loc).Format(time.RFC1123), nil
}
