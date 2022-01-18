package zetcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func Get(ctx context.Context, client *clientv3.Client, key string) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.Get(ctx, key)
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

func GetWithLimit(ctx context.Context, client *clientv3.Client, key string, num int64) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithLimit(num))
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

func GetWithRevision(ctx context.Context, client *clientv3.Client, key string, revision int64) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.Get(ctx, key, clientv3.WithRev(revision))
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

func GetWithSorted(ctx context.Context, client *clientv3.Client, key string) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

func GetWithPrefix(ctx context.Context, client *clientv3.Client, key string) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetWithRange get keys which in [startKey, endKey).
func GetWithRange(ctx context.Context, client *clientv3.Client, startKey string, endKey string) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.Get(ctx, startKey, clientv3.WithRange(endKey))
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetFromKey get keys that are greater than or equal to the given key using byte compare.
func GetFromKey(ctx context.Context, client *clientv3.Client, key string) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.Get(ctx, key, clientv3.WithFromKey())
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetKeysOnly get key return only the keys and the corresponding values will be omitted.
func GetKeysOnly(ctx context.Context, client *clientv3.Client, key string) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetCountOnly get keys return only the count of keys.
func GetCountOnly(ctx context.Context, client *clientv3.Client, key string) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetMinModRevision get key filters out keys for Get with modification revisions less than the given revision.
func GetMinModRevision(ctx context.Context, client *clientv3.Client, key string, revision int64) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithMinModRev(revision))
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetMaxModRevision get key filters out keys for Get with modification revisions greater than the given revision.
func GetMaxModRevision(ctx context.Context, client *clientv3.Client, key string, revision int64) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithMaxModRev(revision))
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetMinCrtRevision get key filters out keys for Get with creation revisions less than the given revision.
func GetMinCrtRevision(ctx context.Context, client *clientv3.Client, key string, revision int64) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithMinCreateRev(revision))
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetMaxCrtRevision get key filters out keys for Get with creation revisions greater than the given revision.
func GetMaxCrtRevision(ctx context.Context, client *clientv3.Client, key string, revision int64) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithMaxCreateRev(revision))
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetFirstCrt gets the key with the oldest creation revision in the request range.
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

// GetLastCrt gets the key with the latest creation revision in the request range.
func GetLastCrt(ctx context.Context, client *clientv3.Client, key string) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	ops := clientv3.WithLastCreate()
	response, err := client.Get(ctx, key, ops...)
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetFirst gets the lexically first key in the request range.
func GetFirst(ctx context.Context, client *clientv3.Client, key string) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	ops := clientv3.WithFirstKey()
	response, err := client.Get(ctx, key, ops...)
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetLast gets the lexically last key in the request range.
func GetLast(ctx context.Context, client *clientv3.Client, key string) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	ops := clientv3.WithLastKey()
	response, err := client.Get(ctx, key, ops...)
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetFirstRevision gets the key with the oldest modification revision in the request range.
func GetFirstRevision(ctx context.Context, client *clientv3.Client, key string) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	ops := clientv3.WithFirstRev()
	response, err := client.Get(ctx, key, ops...)
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}

// GetLastRevision gets the key with the latest modification revision in the request range.
func GetLastRevision(ctx context.Context, client *clientv3.Client, key string) (*clientv3.GetResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	ops := clientv3.WithLastRev()
	response, err := client.Get(ctx, key, ops...)
	if err != nil {
		return nil, data, err
	}

	for _, ev := range response.Kvs {
		data[string(ev.Key)] = string(ev.Value)
	}

	return response, data, nil
}
