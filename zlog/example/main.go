package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ezgroot/ezUtils/zlog"
	"github.com/ezgroot/ezUtils/zlog/common"
)

func main() {
	var config common.Config
	config.LogFilePath = "./log"
	config.LogLevels = common.All
	config.LogFilePrefix = "test"
	//config.Modules = append(config.Modules, "screen")
	//config.Modules = append(config.Modules, "file")
	//config.Modules = append(config.Modules, "zlog")
	config.Modules = append(config.Modules, common.ModulesAll)
	//config.Modules = append(config.Modules, common.ModulesNone)
	config.IsTimeSplit = true
	config.SplitPeriod = 10
	config.IsSizeSplit = true
	config.SplitSize = 1000 * 1000 * 1000
	config.IsClear = true
	config.SavePeriod = 1
	config.UnifyTo = common.UnifyTypeOfOff

	zlog.Init(config)

	zlog.StartPipeServer()

	zlog.Debug("Debug example %s %d", "string", 222)
	zlog.Info("Info example %s %d", "string", 222)
	zlog.Notice("Notice example %s %d", "string", 222)
	zlog.Warn("Warn example %s %d", "string", 222)
	zlog.Error("Error example %s %d", "string", 222)
	zlog.Critical("Critical example %s %d", "string", 222)

	zlog.DebugJs("DebugJs example", "key1", "value1", "key2", 999)
	zlog.InfoJs("InfoJs example", "key1", "value1", "key2", 999)
	zlog.NoticeJs("NoticeJs example", "key1", "value1", "key2", 999)
	zlog.WarnJs("WarnJs example", "key1", "value1", "key2", 999)
	zlog.ErrorJs("ErrorJs example", "key1", "value1", "key2", 999)
	zlog.CriticalJs("CriticalJs example", "key1", "value1", "key2", 999)

	zlog.Debugf("Debugf example %s %d", "string", 444)
	zlog.Infof("Infof example %s %d", "string", 444)
	zlog.Noticef("Noticef example %s %d", "string", 444)
	zlog.Warnf("Warnf example %s %d", "string", 444)
	zlog.Errorf("Errorf example %s %d", "string", 444)
	zlog.Criticalf("Criticalf example %s %d", "string", 444)

	zlog.DebugfJs("DebugJs example", "key1", "value1", "key2", 999)
	zlog.InfofJs("InfoJs example", "key1", "value1", "key2", 999)
	zlog.NoticefJs("NoticeJs example", "key1", "value1", "key2", 999)
	zlog.WarnfJs("WarnJs example", "key1", "value1", "key2", 999)
	zlog.ErrorfJs("ErrorJs example", "key1", "value1", "key2", 999)
	zlog.CriticalfJs("CriticalJs example", "key1", "value1", "key2", 999)

	securityExitProcess(quit)
}

type exitFunc func()

func securityExitProcess(exitFunc exitFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Printf("\n[ INFO ] (system) - security exit by %s signal.\n", s)
			exitFunc()
		default:
			fmt.Printf("\n[ INFO ] (system) - unknown exit by %s signal.\n", s)
			exitFunc()
		}
	}
}

func quit() {
	os.Exit(0)
}
