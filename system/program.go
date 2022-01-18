package system

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

type exitFunc func()

// SecurityExitProcess Listen for system signals, execute exitFunc
// to release program resources, and exit gracefully.
func SecurityExitProcess(exitFunc exitFunc) {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			{
				fmt.Printf("security exit by %s signal.\n", s)
				exitFunc()
				time.Sleep(time.Millisecond * time.Duration(1500))
				os.Exit(0)
			}
		default:
			{
				fmt.Printf("unknown exit by %s signal.\n", s)
				exitFunc()
				time.Sleep(time.Millisecond * time.Duration(1500))
				os.Exit(0)
			}
		}
	}
}

// SecurityExitProcessWithNotify Listen for system signals and notify, execute exitFunc
// to release program resources, and exit gracefully.
func SecurityExitProcessWithNotify(exitFunc exitFunc, cc chan string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		select {
		case s := <-c:
			{
				switch s {
				case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
					{
						fmt.Printf("security exit by %s signal.\n", s)
						exitFunc()
						time.Sleep(time.Millisecond * time.Duration(1500))
						os.Exit(0)
					}
				default:
					{
						fmt.Printf("unknown exit by %s signal.\n", s)
						exitFunc()
						time.Sleep(time.Millisecond * time.Duration(1500))
						os.Exit(0)
					}
				}
			}
		case v := <-cc:
			{
				fmt.Printf("security exit by notify of %s.\n", v)
				exitFunc()
				time.Sleep(time.Millisecond * time.Duration(1500))
				os.Exit(0)
			}
		}
	}
}

// GetGoroutineID Get the coroutine ID (this method affects performance and is only used for debugging).
func GetGoroutineID() int64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		fmt.Printf("%s strconv.ParseUint error = %s\n", string(b), err)
		return -1
	}

	return n
}

func GetPath() (string, error) {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	return path, nil
}

func GetPID() int {
	return os.Getpid()
}
