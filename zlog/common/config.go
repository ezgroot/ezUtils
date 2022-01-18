package common

// Config is the configuration of file and screen log impl.
type Config struct {
	LogFilePath   string   `json:"logfilePath"`   // log file path, default use "./"
	LogLevels     int      `json:"logLevels"`     // effect log level, default 63, means all
	LogFilePrefix string   `json:"logFilePrefix"` // log file prefix, default use "server"
	Modules       []string `json:"modules"`       // effect log modules, default "all"
	IsTimeSplit   bool     `json:"isTimeSplit"`   // is split log file depend on time, default true
	SplitPeriod   int64    `json:"splitPeriod"`   // the period of split log file on time，per - second, default 60*60*24
	IsSizeSplit   bool     `json:"isSizeSplit"`   // is split log file depend on file size, default true
	SplitSize     int64    `json:"splitSize"`     // the size of split log file on size，per - byte, default 1024*1024*5 byte
	IsClear       bool     `json:"isClear"`       // is clear expired file
	SavePeriod    int64    `json:"savePeriod"`    // log file save period, per - day, default 7
	UnifyTo       int      `json:"unifyTo"`       // unify all log to screen or file, default 0, means not unify.
}
