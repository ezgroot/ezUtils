package lock

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func CreateLease(ctx context.Context, client *clientv3.Client, leaseTTL int64) (*clientv3.LeaseGrantResponse, clientv3.LeaseID, error) {
	if client == nil {
		return nil, 0, fmt.Errorf("client is nil")
	}

	response, err := client.Grant(ctx, leaseTTL)
	if err != nil {
		return nil, 0, err
	}

	leaseID := response.ID

	return response, leaseID, nil
}

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

func GetFirstCrt(ctx context.Context, client *clientv3.Client, key string) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	ops := clientv3.WithFirstCreate()
	response, err := client.Get(ctx, key, ops...)
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetRemoteLockState 获取远程锁的状态
func GetRemoteLockState(ctx context.Context, client *clientv3.Client, lockName string) (bool, int64, error) {
	ctxTmp, cancel := context.WithTimeout(ctx, time.Duration(1500)*time.Millisecond)
	defer cancel()

	response, data, err := GetFirstCrt(ctxTmp, client, lockName)
	if err != nil {
		return false, 0, err
	}

	if len(data) == 0 {
		return false, 0, nil
	}

	return true, response.Kvs[0].CreateRevision, nil
}
