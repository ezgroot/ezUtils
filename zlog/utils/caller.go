package utils

import (
	"runtime"
	"strings"

	"github.com/ezgroot/ezUtils/zlog/common"
)

// ThirdCallerInfo third layout caller
func ThirdCallerInfo(log *common.OneLog) {
	level := 3
	pc, file, line, _ := runtime.Caller(level)

	pcInfoStr := runtime.FuncForPC(pc).Name()

	pkgName := ""
	funcName := ""

	lastPathSplitIndex := strings.LastIndex(pcInfoStr, "/")
	if lastPathSplitIndex <= 0 {
		firstfuncSplitIndex := strings.Index(pcInfoStr, ".")
		pkgName = pcInfoStr[:firstfuncSplitIndex]
		funcName = pcInfoStr[firstfuncSplitIndex+1:]
	} else {
		pkgStr := pcInfoStr[:lastPathSplitIndex]
		funcStr := pcInfoStr[lastPathSplitIndex+1:]

		firstfuncSplitIndex := strings.Index(funcStr, ".")

		pkgName = pkgStr + "/" + funcStr[:firstfuncSplitIndex]
		funcName = funcStr[firstfuncSplitIndex+1:]
	}

	if !strings.Contains(funcName, "()") {
		funcName = funcName + "()"
	}

	log.CallerFile = file
	log.CallerLine = line
	log.CallerName = funcName
	log.CallerPkg = pkgName

	return
}
