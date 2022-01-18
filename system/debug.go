package system

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// DebugServerMode Turn on debugging information logging, and you can view
// the specific usage of program memory, coroutines, etc. in the link.
// Note you need at main package import _ "net/http/pprof".
// More, recommend you add import _ "github.com/mkevac/debugcharts" to see
// more char info.
func DebugServerMode(port int) {
	go func() {
		fmt.Printf("Remote see at terminal:\n"+
			"# go tool pprof -http=:8081 http://remoteDebugHostIP:%d/debug/pprof/heap\n"+
			"Web:\n"+
			"1、http://localhost:%d/debug/pprof\n"+
			"2、http://localhost:%d/debug/charts\n"+
			"if get some problem, please make sure you have\n"+
			"import '_ net/http/pprof' and '_ github.com/mkevac/debugcharts'\n",
			port, port, port)
		err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
		if err != nil {
			fmt.Printf("ListenAndServe pprof error = %s\n", err)
			return
		}
	}()
}

// Caller returns information of the assigned level calling function.
// level = 0 - get Caller() itself info;
// level = 1 - get Caller()'s father info;
// level = 2 - get Caller()'s grandfather info;
func Caller(level int) (file string, line int, pkgName string, funcName string) {
	pc, file, line, ok := runtime.Caller(level)
	if !ok {
		return "", 0, "", ""
	}

	pcInfoStr := runtime.FuncForPC(pc).Name()

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

	return file, line, pkgName, funcName
}
