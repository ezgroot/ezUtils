package zetcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Compact compacts etcd KV history before the given rev.
// To avoid accumulating an infinite amount of historical data,
// it is important to compress past revisions. After compaction,
// etcd deletes historical revisions, freeing resources for future use.
// The superseded data of all revisions prior to the compressed revision will be inaccessible.
func Compact(ctx context.Context, client *clientv3.Client, version int64) (response *clientv3.CompactResponse, err error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	response, err = client.Compact(ctx, version)
	if err != nil {
		return nil, err
	}

	return response, nil
}
