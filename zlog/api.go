package zlog

import (
	"github.com/ezgroot/ezUtils/zlog/common"
	"github.com/ezgroot/ezUtils/zlog/loger"
	"github.com/ezgroot/ezUtils/zlog/server"
)

// Init init file and screen log impl, otherwise will use default config.
func Init(config common.Config) {
	loger.GetManager().SetConfig(config)
}

func StartPipeServer() {
	go server.GetInstance().Start()
}

func Debug(format string, v ...interface{}) {
	if format == "" {
		length := len(v)
		for i := 0; i < length; i++ {
			format += "%v "
		}
	}

	loger.GetManager().Add(common.UnifyTypeOfScreen, common.FormatToString, common.Debug, format, v...)
}

func Info(format string, v ...interface{}) {
	if format == "" {
		length := len(v)
		for i := 0; i < length; i++ {
			format += "%v "
		}
	}

	loger.GetManager().Add(common.UnifyTypeOfScreen, common.FormatToString, common.Info, format, v...)
}

func Notice(format string, v ...interface{}) {
	if format == "" {
		length := len(v)
		for i := 0; i < length; i++ {
			format += "%v "
		}
	}

	loger.GetManager().Add(common.UnifyTypeOfScreen, common.FormatToString, common.Notice, format, v...)
}

func Warn(format string, v ...interface{}) {
	if format == "" {
		length := len(v)
		for i := 0; i < length; i++ {
			format += "%v "
		}
	}

	loger.GetManager().Add(common.UnifyTypeOfScreen, common.FormatToString, common.Warn, format, v...)
}

func Error(format string, v ...interface{}) {
	if format == "" {
		length := len(v)
		for i := 0; i < length; i++ {
			format += "%v "
		}
	}

	loger.GetManager().Add(common.UnifyTypeOfScreen, common.FormatToString, common.Error, format, v...)
}

func Critical(format string, v ...interface{}) {
	if format == "" {
		length := len(v)
		for i := 0; i < length; i++ {
			format += "%v "
		}
	}

	loger.GetManager().Add(common.UnifyTypeOfScreen, common.FormatToString, common.Critical, format, v...)
}

func Debugf(format string, v ...interface{}) {
	if format == "" {
		length := len(v)
		for i := 0; i < length; i++ {
			format += "%v "
		}
	}

	loger.GetManager().Add(common.UnifyTypeOfFile, common.FormatToString, common.Debug, format, v...)
}

func Infof(format string, v ...interface{}) {
	if format == "" {
		length := len(v)
		for i := 0; i < length; i++ {
			format += "%v "
		}
	}

	loger.GetManager().Add(common.UnifyTypeOfFile, common.FormatToString, common.Info, format, v...)
}

func Noticef(format string, v ...interface{}) {
	if format == "" {
		length := len(v)
		for i := 0; i < length; i++ {
			format += "%v "
		}
	}

	loger.GetManager().Add(common.UnifyTypeOfFile, common.FormatToString, common.Notice, format, v...)
}

func Warnf(format string, v ...interface{}) {
	if format == "" {
		length := len(v)
		for i := 0; i < length; i++ {
			format += "%v "
		}
	}

	loger.GetManager().Add(common.UnifyTypeOfFile, common.FormatToString, common.Warn, format, v...)
}

func Errorf(format string, v ...interface{}) {
	if format == "" {
		length := len(v)
		for i := 0; i < length; i++ {
			format += "%v "
		}
	}

	loger.GetManager().Add(common.UnifyTypeOfFile, common.FormatToString, common.Error, format, v...)
}

func Criticalf(format string, v ...interface{}) {
	if format == "" {
		length := len(v)
		for i := 0; i < length; i++ {
			format += "%v "
		}
	}

	loger.GetManager().Add(common.UnifyTypeOfFile, common.FormatToString, common.Critical, format, v...)
}

func DebugJs(msg string, v ...interface{}) {
	loger.GetManager().Add(common.UnifyTypeOfScreen, common.FormatToJSON, common.Debug, msg, v...)
}

func InfoJs(msg string, v ...interface{}) {
	loger.GetManager().Add(common.UnifyTypeOfScreen, common.FormatToJSON, common.Info, msg, v...)
}

func NoticeJs(msg string, v ...interface{}) {
	loger.GetManager().Add(common.UnifyTypeOfScreen, common.FormatToJSON, common.Notice, msg, v...)
}

func WarnJs(msg string, v ...interface{}) {
	loger.GetManager().Add(common.UnifyTypeOfScreen, common.FormatToJSON, common.Warn, msg, v...)
}

func ErrorJs(msg string, v ...interface{}) {
	loger.GetManager().Add(common.UnifyTypeOfScreen, common.FormatToJSON, common.Error, msg, v...)
}

func CriticalJs(msg string, v ...interface{}) {
	loger.GetManager().Add(common.UnifyTypeOfScreen, common.FormatToJSON, common.Critical, msg, v...)
}

func DebugfJs(msg string, v ...interface{}) {
	loger.GetManager().Add(common.UnifyTypeOfFile, common.FormatToJSON, common.Debug, msg, v...)
}

func InfofJs(msg string, v ...interface{}) {
	loger.GetManager().Add(common.UnifyTypeOfFile, common.FormatToJSON, common.Info, msg, v...)
}

func NoticefJs(msg string, v ...interface{}) {
	loger.GetManager().Add(common.UnifyTypeOfFile, common.FormatToJSON, common.Notice, msg, v...)
}

func WarnfJs(msg string, v ...interface{}) {
	loger.GetManager().Add(common.UnifyTypeOfFile, common.FormatToJSON, common.Warn, msg, v...)
}

func ErrorfJs(msg string, v ...interface{}) {
	loger.GetManager().Add(common.UnifyTypeOfFile, common.FormatToJSON, common.Error, msg, v...)
}

func CriticalfJs(msg string, v ...interface{}) {
	loger.GetManager().Add(common.UnifyTypeOfFile, common.FormatToJSON, common.Critical, msg, v...)
}
