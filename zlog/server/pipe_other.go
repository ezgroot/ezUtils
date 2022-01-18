// +build !linux

package server

import (
	"time"
)

type PipeServer struct {
	dataCh chan []byte
}

func (ps *PipeServer) Start() error {
	return nil
}

func (ps *PipeServer) GetDataCh() <-chan []byte {
	return ps.dataCh
}

func (ps *PipeServer) Close() {

}

func NewPipeServer(file string, period time.Duration) *PipeServer {
	ps := &PipeServer{
		dataCh: make(chan []byte, 1),
	}

	return ps
}
