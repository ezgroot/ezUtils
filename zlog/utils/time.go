package utils

import (
	"strconv"
	"time"

	"github.com/ezgroot/ezUtils/zlog/common"
)

// GetTime get time
func GetTime(nanoSec int64) time.Time {
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

// GetLogTimeStr format time to string
func GetLogTimeStr(log *common.OneLog) string {
	t := GetTime(log.Timestamp)
	timestamp := strconv.FormatInt(log.Timestamp, 10)
	var timeTmp string
	if len(timestamp) >= 19 {
		timeTmp = "." + string(timestamp[10:19])
	}
	timeNow := t.Format("2006-01-02 15:04:05")

	return timeNow + timeTmp
}
