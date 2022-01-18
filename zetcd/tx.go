package zetcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// TxnPut txn multi write.
func TxnPut(ctx context.Context, client *clientv3.Client, KVList map[string]string) (*clientv3.TxnResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	var opList []clientv3.Op
	for k, v := range KVList {
		opList = append(opList, clientv3.OpPut(k, v))
	}

	response, err := client.Txn(ctx).If().Then(opList...).Else(opList...).Commit()
	if err != nil {
		return nil, err
	}

	return response, nil
}

// TxnPutWithLease txn multi lease.
func TxnPutWithLease(ctx context.Context, client *clientv3.Client, KVList map[string]string, leaseID clientv3.LeaseID) (*clientv3.TxnResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	var opList []clientv3.Op
	for k, v := range KVList {
		opList = append(opList, clientv3.OpPut(k, v, clientv3.WithLease(leaseID)))
	}

	response, err := client.Txn(ctx).If().Then(opList...).Else(opList...).Commit()
	if err != nil {
		return nil, err
	}

	return response, nil
}

// TxnDelete txn multi delete.
func TxnDelete(ctx context.Context, client *clientv3.Client, KVList map[string]string) (*clientv3.TxnResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	var opList []clientv3.Op
	for k := range KVList {
		opList = append(opList, clientv3.OpDelete(k))
	}

	response, err := client.Txn(ctx).If().Then(opList...).Else(opList...).Commit()
	if err != nil {
		return nil, err
	}

	return response, nil
}
