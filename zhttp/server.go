package zhttp

import (
	"net/http"
	"time"
)

func StartServer(addr string, handler http.Handler, timeout time.Duration) error {
	listener, err := NewListener(addr, timeout, timeout)
	if err != nil {
		return err
	}

	if err = http.Serve(listener, handler); err != nil {
		return err
	}

	return nil
}
