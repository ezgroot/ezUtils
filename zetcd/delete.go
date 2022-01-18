package zetcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func Delete(ctx context.Context, client *clientv3.Client, key string) (*clientv3.DeleteResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	response, err := client.Delete(ctx, key)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func DeletePrefix(ctx context.Context, client *clientv3.Client, key string) (*clientv3.DeleteResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	response, err := client.Delete(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteRange delete key-value which key in range of [startKey, endKey).
func DeleteRange(ctx context.Context, client *clientv3.Client, startKey string, endKey string) (*clientv3.DeleteResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	response, err := client.Delete(ctx, startKey, clientv3.WithRange(endKey))
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteFromKey Removes all key-value pairs that key match >= given key (using byte comparison).
func DeleteFromKey(ctx context.Context, client *clientv3.Client, key string) (*clientv3.DeleteResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	response, err := client.Delete(ctx, key, clientv3.WithFromKey())
	if err != nil {
		return nil, err
	}

	return response, nil
}
