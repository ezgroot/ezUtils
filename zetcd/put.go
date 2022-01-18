package zetcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Put a key-value pair.
func Put(ctx context.Context, client *clientv3.Client, key string, value string) (*clientv3.PutResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	response, err := client.Put(ctx, key, value)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// PutWithLease a key-value pair with lease.
func PutWithLease(ctx context.Context, client *clientv3.Client, key string, value string, leaseID clientv3.LeaseID) (*clientv3.PutResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	response, err := client.Put(ctx, key, value, clientv3.WithLease(leaseID))
	if err != nil {
		return nil, err
	}

	return response, nil
}
