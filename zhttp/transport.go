package zhttp

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func Transport() *http.Transport {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,

		DialContext: (&net.Dialer{
			Timeout:   time.Duration(3) * time.Second,
			DualStack: true,
			KeepAlive: time.Duration(30) * time.Second,
			// Resolver
			// Cancel
			Control: control,
		}).DialContext,
		MaxIdleConns:          300,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       time.Duration(60) * time.Second,
		ExpectContinueTimeout: time.Duration(3) * time.Second,
		ForceAttemptHTTP2:     true,
	}

	return tr
}

func TLSTransport(tls *tls.Config) *http.Transport {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialTLSContext: (&net.Dialer{
			Timeout:   time.Duration(3) * time.Second, // the maximum amount of time a dial will wait for a connect to complete.
			DualStack: true,
			KeepAlive: time.Duration(30) * time.Second, // specifies the interval between keep-alive probes for an active network connection.
			// Resolver
			// Cancel
			Control: control, // called after creating the network connection but before actually dialing.
		}).DialContext,
		TLSClientConfig:       tls,
		MaxIdleConns:          300,                             // the maximum number of idle (keep-alive) connections across all hosts. Zero means no limit.
		MaxIdleConnsPerHost:   20,                              //　the maximum idle (keep-alive) connections to keep per-host. If zero,
		IdleConnTimeout:       time.Duration(60) * time.Second, // the maximum amount of time an idle (keep-alive) connection will remain idle before closing itself.
		TLSHandshakeTimeout:   time.Duration(5) * time.Second,
		ExpectContinueTimeout: time.Duration(3) * time.Second,
		ForceAttemptHTTP2:     true,
	}

	return tr
}

func GetTlsConfig(caFilePath string) (*tls.Config, error) {
	caCrt, err := ioutil.ReadFile(caFilePath)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCrt)

	return &tls.Config{RootCAs: pool}, nil
}

func GetSkipTlsConfig() *tls.Config {
	return &tls.Config{InsecureSkipVerify: true}
}
