package zhttp

import (
	"bytes"
	"io"
	"net/http"
)

// NewRequest create a request
func NewRequest(method string, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return req, err
}

func GetRequest2(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return req, err
}
