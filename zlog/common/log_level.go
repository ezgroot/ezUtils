package common

// log level.
const (
	Debug    = 1      // 1
	Info     = 1 << 1 // 2
	Notice   = 1 << 2 // 4
	Warn     = 1 << 3 // 8
	Error    = 1 << 4 // 16
	Critical = 1 << 5 // 32

	All = Debug | Info | Notice | Warn | Error | Critical // 63
)

// LevelMap log level map
var LevelMap = make(map[int]string)

func init() {
	LevelMap[Debug] = "[ DEBUG  ]"
	LevelMap[Info] = "[  INFO  ]"
	LevelMap[Notice] = "[ NOTICE ]"
	LevelMap[Warn] = "[  WARN  ]"
	LevelMap[Error] = "[ ERROR  ]"
	LevelMap[Critical] = "[CRITICAL]"
}
