package loger

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ezgroot/ezUtils/zlog/common"
	"github.com/ezgroot/ezUtils/zlog/filelog"
	"github.com/ezgroot/ezUtils/zlog/screenlog"
	"github.com/ezgroot/ezUtils/zlog/utils"
)

// Manager manage all log
type Manager struct {
	Conf  common.Config
	queue chan *common.OneLog
}

// filterLog filter logs
func (m *Manager) filterLog(level int, module string) bool {
	var retCheck bool

	ret := m.Conf.LogLevels & level
	if ret == 0 {
		return false
	}

	if len(m.Conf.Modules) == 0 {
		return false
	}

	for _, value := range m.Conf.Modules {
		if value == common.ModulesAll {
			retCheck = true
			break
		} else if value == common.ModulesNone {
			retCheck = false
			break
		} else if strings.Contains(value, module) {
			retCheck = true
		} else {
			continue
		}
	}

	return retCheck
}

func (m *Manager) run() {
	for {
		log := <-m.queue

		if log.OutTo == common.UnifyTypeOfScreen {
			screenlog.GetScreenLogImpl().Show(log)
		} else if log.OutTo == common.UnifyTypeOfFile {
			filelog.GetFileLogImpl().Add(log)
		} else {
			screenlog.GetScreenLogImpl().Show(log)
		}
	}
}

// SetConfig set manager config
func (m *Manager) SetConfig(c common.Config) {
	m.Conf = c

	if len(m.Conf.Modules) == 0 {
		m.Conf.Modules = append(m.Conf.Modules, common.ModulesAll)
	}

	if m.Conf.LogLevels == 0 {
		m.Conf.LogLevels = common.All
	}

	filelog.GetFileLogImpl().SetLoggerConfig(m.Conf)
}

// Add add log
func (m *Manager) Add(outTo int, format int, level int, msg string, v ...interface{}) {
	log := &common.OneLog{}

	if !m.filterLog(level, log.CallerPkg) {
		return
	}

	utils.ThirdCallerInfo(log)

	if m.Conf.UnifyTo != 0 {
		log.OutTo = m.Conf.UnifyTo
	} else {
		log.OutTo = outTo
	}

	log.FormatType = format
	log.Level = level
	log.Format = msg
	log.Args = v

	t := time.Now()
	log.Timestamp = t.UnixNano()

	offset := strings.LastIndex(log.CallerFile, "/")
	if offset != 0 {
		log.CallerFile = log.CallerFile[(offset + 1):]
	}

	if (cap(m.queue) - len(m.queue)) <= cap(m.queue)/100 {
		fmt.Printf("manager log queue is nearly full!!!\n")
	}

	m.queue <- log
}

func newManager() *Manager {
	c := common.Config{
		LogFilePath:   common.DefaultLogFilePath,
		LogLevels:     common.All,
		LogFilePrefix: common.DefaultLogFilePrefix,
		Modules:       make([]string, 0, 64),
		IsTimeSplit:   true,
		SplitPeriod:   common.DefaultSplitPeriod,
		IsSizeSplit:   true,
		SplitSize:     common.DefaultSplitSize,
		IsClear:       true,
		SavePeriod:    common.DefaultLogFileSavePeriod,
		UnifyTo:       common.UnifyTypeOfOff,
	}

	c.Modules = append(c.Modules, common.ModulesAll)

	m := &Manager{queue: make(chan *common.OneLog, common.ManagerQueueMaxNumber), Conf: c}
	m.SetConfig(c)

	go m.run()

	return m
}

var instance *Manager
var once sync.Once

// GetManager get manager impl
func GetManager() *Manager {
	once.Do(func() {
		instance = newManager()
	})

	return instance
}
