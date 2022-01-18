package zetcd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	DefaultDialTimeout          = time.Duration(2000) * time.Millisecond
	DefaultDialKeepAlivePeriod  = time.Duration(10000) * time.Millisecond
	DefaultDialKeepAliveTimeout = time.Duration(3000) * time.Millisecond
)

type Config struct {
	NodeList             []string      `json:"nodeList"`
	UseTLS               bool          `json:"userTLS"`
	CaFile               string        `json:"ca"`
	CertFile             string        `json:"cert"`
	CertKeyFile          string        `json:"key"`
	ServerName           string        `json:"serverName"`
	DialTimeout          time.Duration `json:"timeout"`          // per - Millisecond
	DialKeepAlivePeriod  time.Duration `json:"keepAlivePeriod"`  // per - Millisecond
	DialKeepAliveTimeout time.Duration `json:"keepAliveTimeout"` // per - Millisecond
}

func Client(config Config) (*clientv3.Client, error) {
	var clientCfg clientv3.Config

	if config.DialTimeout <= 0 {
		config.DialTimeout = DefaultDialTimeout
	}

	if config.DialKeepAlivePeriod <= 0 {
		config.DialKeepAlivePeriod = DefaultDialKeepAlivePeriod
	}

	if config.DialKeepAliveTimeout <= 0 {
		config.DialKeepAliveTimeout = DefaultDialKeepAliveTimeout
	}

	clientLogConfig := &zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.ErrorLevel),
		Development:   false,
		Sampling:      &zap.SamplingConfig{Initial: 100, Thereafter: 100},
		Encoding:      "json",
		EncoderConfig: zap.NewProductionEncoderConfig(),
		// Use "/dev/null" to discard all
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var dialOptions []grpc.DialOption // TODO: config dial option

	if config.UseTLS {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.CertKeyFile)
		if err != nil {
			return nil, err
		}

		caData, err := ioutil.ReadFile(config.CaFile)
		if err != nil {
			return nil, err
		}

		pool := x509.NewCertPool()
		ok := pool.AppendCertsFromPEM(caData)
		if !ok {
			return nil, fmt.Errorf("cert pool append cert from pem failed")
		}

		var tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
			ServerName:   config.ServerName,
		}

		clientCfg = clientv3.Config{Endpoints: config.NodeList,
			DialTimeout:          config.DialTimeout * time.Millisecond,
			DialKeepAliveTime:    config.DialKeepAlivePeriod * time.Millisecond,
			DialKeepAliveTimeout: config.DialKeepAliveTimeout * time.Millisecond,
			TLS:                  tlsConfig,
			LogConfig:            clientLogConfig,
			DialOptions:          dialOptions}
	} else {
		clientCfg = clientv3.Config{Endpoints: config.NodeList,
			DialTimeout:          config.DialTimeout * time.Millisecond,
			DialKeepAliveTime:    config.DialKeepAlivePeriod * time.Millisecond,
			DialKeepAliveTimeout: config.DialKeepAliveTimeout * time.Millisecond,
			LogConfig:            clientLogConfig,
			DialOptions:          dialOptions}
	}

	client, err := clientv3.New(clientCfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Close(client *clientv3.Client) {
	if client == nil {
		return
	}

	err := client.Close()
	if err != nil {
		fmt.Printf("close client error = %s", err)
	}
}
