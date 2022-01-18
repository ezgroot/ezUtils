package zhttp

import (
	"errors"
	"net/http"
	"syscall"
	"time"
)

// default config
const (
	defaultTimeout = time.Duration(3000) * time.Millisecond
)

var ErrorOfEmptyClient = errors.New("client is nil")

func control(network, address string, c syscall.RawConn) error {
	//fmt.Printf("network=%s, address=%s", network, address)
	return nil
}

func checkRedirect(req *http.Request, via []*http.Request) error {

	return nil
}
