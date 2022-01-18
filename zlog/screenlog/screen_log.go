package screenlog

import (
	"fmt"

	"github.com/ezgroot/ezUtils/zlog/common"
	"github.com/ezgroot/ezUtils/zlog/utils"
)

const (
	backgroundColour = 0
)

// LoggerImpl screen log impl
type LoggerImpl struct{}

// newLoggerImpl create screen log impl
func newLoggerImpl() *LoggerImpl {
	impl := &LoggerImpl{}

	return impl
}

// /home/groot/Work/9-openSource/gin/logger.go
const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

func getColour(log *common.OneLog) (string, string) {
	var isHighLight int
	var background int
	var foreground int

	if log.Level == common.Debug {
		background = backgroundColour
		foreground = 37
	} else if log.Level == common.Info {
		background = backgroundColour
		foreground = 32
	} else if log.Level == common.Notice {
		background = backgroundColour
		foreground = 34
	} else if log.Level == common.Warn {
		background = backgroundColour
		foreground = 33
	} else if log.Level == common.Error {
		background = backgroundColour
		foreground = 35
	} else if log.Level == common.Critical {
		background = backgroundColour
		foreground = 31
	} else {
		background = backgroundColour
		foreground = 36
	}

	var colourBegin = ""
	var colourEnd = ""
	if utils.GetOS() == utils.Windows {
		colourBegin = ""
		colourEnd = "\n"
	} else if utils.GetOS() == utils.Linux {
		colourBegin = fmt.Sprintf("%c[%d;%d;%dm", 0x1B, isHighLight, background, foreground)

		tempEnd := [5]byte{0x1b, 0x5b, 0x30, 0x6d, 0x0a}
		colourEnd = string(tempEnd[:]) // fmt.Sprintf("%c[0m\n", 0x1B)
	} else if utils.GetOS() == utils.Darwin {
		colourBegin = fmt.Sprintf("%c[%d;%d;%dm", 0x1B, isHighLight, background, foreground)
		tempEnd := [5]byte{0x1b, 0x5b, 0x30, 0x6d, 0x0a}
		colourEnd = string(tempEnd[:]) // fmt.Sprintf("%c[0m\n", 0x1B)
	} else {
		fmt.Printf("error system type = %s\n", utils.GetOS())
	}

	return colourBegin, colourEnd
}

func (s *LoggerImpl) showLogStr(log *common.OneLog) {
	colourBegin, colourEnd := getColour(log)

	timeStr := utils.GetLogTimeStr(log)

	var formatLog = ""

	formatLog = fmt.Sprintf("%s%s%s%s:%d ◆ %s ★ %s%s%s",
		colourBegin,
		common.LevelMap[log.Level],
		" "+timeStr+" ⇔ ", // ⊗  ⇔
		log.CallerFile,
		log.CallerLine,
		log.CallerPkg,
		log.CallerName+" ▶ ",
		log.Format,
		colourEnd)

	fmt.Printf(formatLog, log.Args...)
}

func (s *LoggerImpl) showLogJSON(log *common.OneLog) {
	colourBegin, colourEnd := getColour(log)
	timeStr := utils.GetLogTimeStr(log)

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

	fmt.Printf("%s%s%s", colourBegin, str1, colourEnd)
}

// Show output a screen log
func (s *LoggerImpl) Show(log *common.OneLog) {
	if log.FormatType == common.FormatToString {
		s.showLogStr(log)
	} else if log.FormatType == common.FormatToJSON {
		s.showLogJSON(log)
	} else {
		s.showLogStr(log)
	}
}
