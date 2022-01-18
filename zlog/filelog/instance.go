package filelog

import (
	"sync"
)

var instance *LoggerImpl
var once sync.Once

// GetFileLogImpl get file log impl
func GetFileLogImpl() *LoggerImpl {
	once.Do(func() {
		instance = newLoggerImpl()
	})

	return instance
}
