// +build linux

package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"
)

const (
	DataScanPeriodOfMin = time.Duration(100) * time.Millisecond
)

type PipeServer struct {
	closeNotifyCh chan struct{}
	dataCh        chan []byte
	filePath      string
	fileHandle    *os.File
	scanPeriod    time.Duration
}

func (ps *PipeServer) checkPipeFile() error {
	_, err := os.Stat(ps.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("[WARN] file not exist, begin create\n")
			err = syscall.Mkfifo(ps.filePath, 0666)
			if err != nil {
				fmt.Printf("[ERROR] create file error = %s\n", err)
				return err
			}
		} else {
			fmt.Printf("[ERROR] check file stat error = %s\n", err)
			return err
		}
	}

	return nil
}

func (ps *PipeServer) run() {
	var err error

	ps.fileHandle, err = os.OpenFile(ps.filePath, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Printf("[WARN] open pipe file error = %s\n", err)
		return
	}

	timer := time.NewTimer(ps.scanPeriod)
	defer timer.Stop()

	reader := bufio.NewReader(ps.fileHandle)

	for {
		select {
		case <-ps.closeNotifyCh:
			{
				return
			}
		case <-timer.C:
			{
				line, err := reader.ReadBytes('\n')
				if err != nil {
					if err != io.EOF {
						fmt.Printf("[WARN] read line error = %s\n", err)
					}

					timer.Reset(ps.scanPeriod)
					continue
				}

				ps.dataCh <- bytes.Trim(line, "\n")
				timer.Reset(ps.scanPeriod)
			}
		}
	}
}

func (ps *PipeServer) Start() error {
	err := ps.checkPipeFile()
	if err != nil {
		return err
	}

	go ps.run()

	return nil
}

func (ps *PipeServer) GetDataCh() <-chan []byte {
	return ps.dataCh
}

func (ps *PipeServer) Close() {
	ps.closeNotifyCh <- struct{}{}

	err := os.Remove(ps.filePath)
	if err != nil {
		fmt.Printf("[WARN] %s file delete error = %s\n", ps.filePath, err)
	}
}

func NewPipeServer(file string, period time.Duration) *PipeServer {
	ps := &PipeServer{
		closeNotifyCh: make(chan struct{}, 1),
		dataCh:        make(chan []byte, 1),
		filePath:      file,
		scanPeriod:    period,
	}

	if ps.scanPeriod <= DataScanPeriodOfMin {
		ps.scanPeriod = DataScanPeriodOfMin
	}

	return ps
}
