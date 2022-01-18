package common

// FileLogQueueMaxNumber file log queue max length.
var FileLogQueueMaxNumber = 300 * 10000 // 300w

// ManagerQueueMaxNumber manager log queue max length
var ManagerQueueMaxNumber = 500 * 10000 // 500w

// default config
const (
	ModulesAll  = "all"  // display all modules log
	ModulesNone = "none" // hide all modules log

	DefaultLogFilePath   = "./"
	DefaultLogFilePrefix = "server"

	DefaultSplitPeriod       = 60 * 60 * 24
	DefaultSplitSize         = 1024 * 1024 * 20
	DefaultLogFileSavePeriod = 7

	DefaultHTTPServerPort = 60000

	UnifyTypeOfOff    = 0
	UnifyTypeOfScreen = 1
	UnifyTypeOfFile   = 2

	FormatToString = 0
	FormatToJSON   = 1
)
