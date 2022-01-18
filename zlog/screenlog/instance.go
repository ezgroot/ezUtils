package screenlog

import "sync"

var instance *LoggerImpl
var once sync.Once

// GetScreenLogImpl get screen log impl
func GetScreenLogImpl() *LoggerImpl {
	once.Do(func() {
		instance = newLoggerImpl()
	})

	return instance
}
