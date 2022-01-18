package zhttp

import "net/http"

// SetHeader set request header， like contentType = "application/json;charset=utf-8" and so on.
func SetHeader(req *http.Request, data map[string]string) {
	for k, v := range data {
		req.Header.Set(k, v)
	}
}
