package filelog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ezgroot/ezUtils/zlog/common"
	"github.com/ezgroot/ezUtils/zlog/utils"
)

// LoggerImpl file log impl
type LoggerImpl struct {
	log             *log.Logger
	mutex           sync.Mutex
	curFileTime     int64
	logFileFullPath string

	logQueue chan common.OneLog

	logFilePath   string
	logFilePrefix string
	isTimeSplit   bool
	splitPeriod   int64
	isSizeSplit   bool
	splitSize     int64
	isClear       bool
	savePeriod    int64
}

func newLoggerImpl() *LoggerImpl {
	var modulesTemp = make([]string, 0)

	modulesTemp = append(modulesTemp, common.ModulesAll)

	impl := &LoggerImpl{
		logFilePath: common.DefaultLogFilePath,
		isTimeSplit: true,
		splitPeriod: common.DefaultSplitPeriod,
		isSizeSplit: true,
		splitSize:   common.DefaultSplitSize,
		logQueue:    make(chan common.OneLog, common.FileLogQueueMaxNumber),
		isClear:     true,
		savePeriod:  common.DefaultLogFileSavePeriod}

	go impl.listenLogQueue()

	go impl.clearAndRecycle()

	return impl
}

func (f *LoggerImpl) writeStr(log common.OneLog) {
	var prefixStr string

	t := utils.GetTime(log.Timestamp)
	timestamp := strconv.FormatInt(log.Timestamp, 10)
	var timeTmp = string(timestamp[10:19])
	timeNow := t.Format("2006-01-02 15:04:05")
	timeStr := timeNow + " @ " + timeTmp

	prefixStr += fmt.Sprintf("%s:%d %s() %s ", log.CallerFile, log.CallerLine, log.CallerName, timeStr)

	prefixStr = fmt.Sprintf("%s%s%s:%d ◆ %s ★ %s",
		common.LevelMap[log.Level],
		" "+timeStr+" ⇔ ",
		log.CallerFile,
		log.CallerLine,
		log.CallerPkg,
		log.CallerName+" ▶ ")

	f.log.SetPrefix(prefixStr)
	f.log.Printf(log.Format, log.Args...)
}

func (f *LoggerImpl) writeJSON(log common.OneLog) {
	timeStr := utils.GetLogTimeStr(&log)

	str1 := "\n" + "    " + utils.GetJSONStr("level", common.LevelMap[log.Level]) + ",\n" +
		"    " + utils.GetJSONStr("time", timeStr) + ",\n" +
		"    " + utils.GetJSONStr("file", log.CallerFile) + ",\n" +
		"    " + utils.GetJSONStr("line", log.CallerLine) + ",\n" +
		"    " + utils.GetJSONStr("pkg", log.CallerPkg) + ",\n" +
		"    " + utils.GetJSONStr("func", log.CallerName) + ",\n" +
		"    " + utils.GetJSONStr("msg", log.Format) + ",\n"

	var key string
	var ok bool
	for i, value := range log.Args {
		if i%2 == 0 {
			key, ok = value.(string)
			if !ok {

			}
		} else {
			if i == len(log.Args)-1 {
				str1 += ("    " + utils.GetJSONStr(key, value) + "\n")
			} else {
				str1 += ("    " + utils.GetJSONStr(key, value) + ",\n")
			}
		}
	}

	str1 = "{" + str1 + "}"
	f.log.SetPrefix("")
	f.log.Printf("%s", str1)
}

func (f *LoggerImpl) write(log common.OneLog) {
	if log.FormatType == common.FormatToString {
		f.writeStr(log)
	} else if log.FormatType == common.FormatToJSON {
		f.writeJSON(log)
	} else {
		f.writeStr(log)
	}
}

func (f *LoggerImpl) listenLogQueue() {
	for l := range f.logQueue {
		if f.isSplitLogFile() {
			err := f.initFileLogImpl()
			if err != nil {
				fmt.Printf("[WARN] init fileLog impl error = %s\n", err)
				f.logQueue <- l
				continue
			}
		}

		f.write(l)
	}
}

func (f *LoggerImpl) isSplitLogFile() bool {
	t := time.Now()
	timeNow := t.UTC().Unix()

	if (timeNow - f.curFileTime) >= f.splitPeriod {
		return true
	}

	fileInfo, err := utils.IsExist(f.logFileFullPath)
	if err != nil {
		fmt.Printf("[WARN] check log file error = %s\n", err)
		return false
	}

	if fileInfo.Size() >= f.splitSize {
		return true
	}

	return false
}

func (f *LoggerImpl) initFileLogImpl() error {
	t := time.Now()
	f.curFileTime = t.UTC().Unix()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	var timeTmp = string(timestamp[10:19])
	timeNow := t.Format("2006-01-02 15:04:05")
	timeNow = strings.Replace(timeNow, " ", "@", -1)
	filenameTmp := strings.Replace(timeNow, ":", "-", -1)

	var fileName string
	if utils.GetOS() == utils.Windows {
		fileName = f.logFilePath + "\\" + f.logFilePrefix + "=" + filenameTmp + "@" + timeTmp + ".log"
	} else if utils.GetOS() == utils.Linux {
		if strings.HasSuffix(f.logFilePath, "/") {
			fileName = f.logFilePath + f.logFilePrefix + "=" + filenameTmp + "@" + timeTmp + ".log"
		} else {
			fileName = f.logFilePath + "/" + f.logFilePrefix + "=" + filenameTmp + "@" + timeTmp + ".log"
		}
	} else if utils.GetOS() == utils.Darwin {
		fileName = f.logFilePath + "/" + f.logFilePrefix + "=" + filenameTmp + "@" + timestamp + ".log"
	} else {
		return fmt.Errorf("[WARN] error system type = %s", utils.GetOS())
	}

	f.logFileFullPath = fileName

	var logFile io.Writer
	_, err := os.Stat(f.logFilePath)
	if err == nil {
		logFile, err = os.Create(fileName)
		if err != nil {
			fmt.Printf("[WARN] create log file error = %s\n", err)
			return err
		}

		fmt.Printf("[WARN] create new log file = %s\n", fileName)
	} else {
		if os.IsNotExist(err) {
			fmt.Printf("[WARN] log dir is not exist, begin to create\n")
			err = os.MkdirAll(f.logFilePath, 0777)
			if err != nil {
				fmt.Printf("[WARN] create log dir error = %s\n", err)
				return err
			}

			logFile, err = os.Create(fileName)
			if err != nil {
				fmt.Printf("[WARN] create log file error = %s\n", err)
			} else {
				fmt.Printf("[WARN] create new log file = %s\n", fileName)
			}
		} else {
			fmt.Printf("[WARN] check log dir error = %s\n", err)
			return err
		}
	}

	f.log = log.New(logFile, "[INIT]-", 0)

	return err
}

func (f *LoggerImpl) clearAndRecycle() {
	for {
		if f.isClear {
			_, err := os.Stat(f.logFilePath)
			if err != nil {
				time.Sleep(time.Duration(30) * time.Second)
				continue
			}

			rd, err := ioutil.ReadDir(f.logFilePath)
			if err == nil {
				for _, fi := range rd {
					if fi.IsDir() {
						continue
					}

					var fileFullPath string
					if utils.GetOS() == utils.Windows {
						fileFullPath = f.logFilePath + "\\" + fi.Name()
					} else if utils.GetOS() == utils.Linux {
						if strings.HasSuffix(f.logFilePath, "/") {
							fileFullPath = f.logFilePath + fi.Name()
						} else {
							fileFullPath = f.logFilePath + "/" + fi.Name()
						}
					} else if utils.GetOS() == utils.Darwin {
						fileFullPath = f.logFilePath + "/" + fi.Name()
					} else {
						fmt.Printf("[WARN] error system type = %s\n", utils.GetOS())
						continue
					}

					fileInfo, err := utils.IsExist(fileFullPath)
					if err != nil {
						fmt.Printf("[WARN] check file exist error = %s\n", err)
						continue
					}

					if strings.HasPrefix(fileInfo.Name(), f.logFilePrefix) && strings.HasSuffix(fileInfo.Name(), ".log") {
						timeFile := fileInfo.ModTime()
						fileExpireTime := timeFile.Add(time.Duration(f.savePeriod*60*60*24) * time.Second)
						if fileExpireTime.Before(time.Now()) {
							err = os.Remove(fileFullPath)
							if err != nil {
								fmt.Printf("[WARN] delete file error = %s\n", err)
							} else {
								fmt.Printf("[INFO] delete expired log file = %s", fileFullPath)
							}
						}
					}
				}
			} else {
				fmt.Printf("[WARN] read dir = %s error = %s\n", f.logFileFullPath, err)
			}
		}

		time.Sleep(time.Duration(60*60*12) * time.Second)
	}
}

// Add push a file log to queue.
func (f *LoggerImpl) Add(log *common.OneLog) {
	f.logQueue <- *log
}

// SetLoggerConfig set screen logger config
func (f *LoggerImpl) SetLoggerConfig(c common.Config) {
	if c.LogFilePath == "" {
		f.logFilePath = common.DefaultLogFilePath
	} else {
		f.logFilePath = c.LogFilePath
	}

	f.isTimeSplit = c.IsTimeSplit
	if c.SplitPeriod <= 0 {
		f.splitPeriod = common.DefaultSplitPeriod
	} else {
		f.splitPeriod = c.SplitPeriod
	}

	f.isSizeSplit = c.IsSizeSplit
	if c.SplitSize <= 0 {
		f.splitSize = common.DefaultSplitSize
	} else {
		f.splitSize = c.SplitSize
	}

	if !c.IsClear {
		f.isClear = c.IsClear
	}

	if c.SavePeriod <= 0 {
		f.savePeriod = common.DefaultLogFileSavePeriod
	} else {
		f.savePeriod = c.SavePeriod
	}

	if c.LogFilePrefix == "" {
		f.logFilePrefix = common.DefaultLogFilePrefix
	} else {
		f.logFilePrefix = c.LogFilePrefix
	}

	f.log = nil
}
