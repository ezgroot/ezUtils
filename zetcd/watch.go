package zetcd

import (
	"context"
	"fmt"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	MethodCreate = "create"
	MethodModify = "modify"
	MethodDelete = "delete"
)

type WatchData struct {
	Operate string
	Vision  int64
	Key     string
	Value   string
}

// Watch a key change event.
func Watch(ctx context.Context,
	client *clientv3.Client,
	key string,
	data chan WatchData,
	stopCh chan struct{}) {
	rch := client.Watch(ctx, key)

	for {
		select {
		case wr := <-rch:
			{
				for _, ev := range wr.Events {
					if ev.IsCreate() {
						var watchData WatchData
						watchData.Operate = MethodCreate
						watchData.Vision = wr.Header.Revision
						watchData.Key = string(ev.Kv.Key)
						watchData.Value = string(ev.Kv.Value)
						data <- watchData
					} else if ev.IsModify() {
						var watchData WatchData
						watchData.Operate = MethodModify
						watchData.Vision = wr.Header.Revision
						watchData.Key = string(ev.Kv.Key)
						watchData.Value = string(ev.Kv.Value)
						data <- watchData
					} else if ev.Type == mvccpb.DELETE {
						var watchData WatchData
						watchData.Operate = MethodDelete
						watchData.Vision = wr.Header.Revision
						watchData.Key = string(ev.Kv.Key)
						watchData.Value = string(ev.Kv.Value)
						data <- watchData
					} else {
						fmt.Printf("watch unknown operate type = %d, key = %s, value = %s",
							ev.Type, ev.Kv.Key, ev.Kv.Value)
					}
				}
			}
		case <-stopCh:
			{
				return
			}
		}

	}
}

// Watch supports setting optional configuration. The list of supported configuration items is as follows:
// 1 - func WithRev(rev int64)
// 2 - func WithPrefix()
// 3 - func WithRange(endKey string)
// 4 - func WithFromKey()
// 5 - func WithProgressNotify()
// 6 - func WithCreatedNotify()
// 7 - func WithFilterPut()
// 8 - func WithFilterDelete()
// 9 - func WithPrevKV()
// 10 - func WithFragment()
// 11 - func WithIgnoreValue()
// 12 - func WithIgnoreLease()

// WatchPrefix a prefix key change event.
func WatchPrefix(ctx context.Context,
	client *clientv3.Client,
	key string,
	data chan WatchData,
	stopCh chan struct{}) {
	rch := client.Watch(ctx, key, clientv3.WithPrefix())

	for {
		select {
		case wr := <-rch:
			{
				for _, ev := range wr.Events {
					if ev.IsCreate() {
						var watchData WatchData
						watchData.Operate = MethodCreate
						watchData.Vision = wr.Header.Revision
						watchData.Key = string(ev.Kv.Key)
						watchData.Value = string(ev.Kv.Value)
						data <- watchData
					} else if ev.IsModify() {
						var watchData WatchData
						watchData.Operate = MethodModify
						watchData.Vision = wr.Header.Revision
						watchData.Key = string(ev.Kv.Key)
						watchData.Value = string(ev.Kv.Value)
						data <- watchData
					} else if ev.Type == mvccpb.DELETE {
						var watchData WatchData
						watchData.Operate = MethodDelete
						watchData.Vision = wr.Header.Revision
						watchData.Key = string(ev.Kv.Key)
						watchData.Value = string(ev.Kv.Value)
						data <- watchData
					} else {
						fmt.Printf("watch unknown operate type = %d, key = %s, value = %s",
							ev.Type, ev.Kv.Key, ev.Kv.Value)
					}
				}
			}
		case <-stopCh:
			{
				return
			}
		}
	}
}

// WatchRange a key range change event.
func WatchRange(ctx context.Context,
	client *clientv3.Client,
	keyOne string,
	keyTwo string,
	data chan WatchData,
	stopCh chan struct{}) {
	rch := client.Watch(ctx, keyOne, clientv3.WithRange(keyTwo))

	for {
		select {
		case wr := <-rch:
			{
				for _, ev := range wr.Events {
					if ev.IsCreate() {
						var watchData WatchData
						watchData.Operate = MethodCreate
						watchData.Vision = wr.Header.Revision
						watchData.Key = string(ev.Kv.Key)
						watchData.Value = string(ev.Kv.Value)
						data <- watchData
					} else if ev.IsModify() {
						var watchData WatchData
						watchData.Operate = MethodModify
						watchData.Vision = wr.Header.Revision
						watchData.Key = string(ev.Kv.Key)
						watchData.Value = string(ev.Kv.Value)
						data <- watchData
					} else if ev.Type == mvccpb.DELETE {
						var watchData WatchData
						watchData.Operate = MethodDelete
						watchData.Vision = wr.Header.Revision
						watchData.Key = string(ev.Kv.Key)
						watchData.Value = string(ev.Kv.Value)
						data <- watchData
					} else {
						fmt.Printf("watch unknown operate type = %d, key = %s, value = %s",
							ev.Type, ev.Kv.Key, ev.Kv.Value)
					}
				}
			}
		case <-stopCh:
			{
				return
			}
		}
	}
}

// WatchWithProgressNotify with progress notify ???
func WatchWithProgressNotify(ctx context.Context, client *clientv3.Client, key string, data chan WatchData) {
	rch := client.Watch(ctx, key, clientv3.WithProgressNotify())
	wr := <-rch

	fmt.Printf("watch revision = %d", wr.Header.Revision)
	fmt.Printf("watch isProgressNotify = %t", wr.IsProgressNotify())
}
