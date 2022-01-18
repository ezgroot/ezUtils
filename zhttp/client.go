package zhttp

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/publicsuffix"
)

func Client(tr *http.Transport, timeout time.Duration) *http.Client {
	if tr == nil {
		tr = &http.Transport{}
	}

	if timeout <= 0 {
		timeout = defaultTimeout
	}

	// The error returned by the internal package implementation cookiejar.New() must be empty
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	client := &http.Client{
		// Transport strategy, specifying the mechanism for executing independent, single HTTP requests
		Transport: tr,
		// Specify the strategy for handling redirection, the client will call this function field before executing the redirection.
		CheckRedirect: checkRedirect,
		// Jar specifies the cookie manager. If Jar is nil, no cookie will be sent in the request, and the cookie in the reply will be ignored.
		Jar: jar,
		// Timeout = connection time + all redirects + read response body, Timeout zero value means no timeout is set.
		Timeout: timeout,
	}

	return client
}
