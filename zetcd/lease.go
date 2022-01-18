package zetcd

import (
	"context"
	"fmt"

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

func RevokeLease(ctx context.Context, client *clientv3.Client, leaseID clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	response, err := client.Revoke(ctx, leaseID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// KeepAliveAways keep a lease aways.
func KeepAliveAways(ctx context.Context, client *clientv3.Client, leaseID clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	ch, err := client.KeepAlive(ctx, leaseID)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

// KeepAliveOnce keep a lease once.
func KeepAliveOnce(ctx context.Context, client *clientv3.Client, leaseID clientv3.LeaseID) (*clientv3.LeaseKeepAliveResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	response, err := client.KeepAliveOnce(ctx, leaseID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// LeaseTimeToLive retrieves the lease information of the given lease ID.
func LeaseTimeToLive(ctx context.Context, client *clientv3.Client, leaseID clientv3.LeaseID) (*clientv3.LeaseTimeToLiveResponse, map[string]string, error) {
	if client == nil {
		return nil, nil, fmt.Errorf("client is nil")
	}

	data := make(map[string]string)

	response, err := client.TimeToLive(ctx, leaseID, clientv3.WithAttachedKeys())
	if err != nil {
		return nil, data, err
	}

	if response.ID != leaseID {
		return nil, data, fmt.Errorf("leaseID expected %d, got %d", leaseID, response.ID)
	}

	if response.TTL == 0 || response.TTL > response.GrantedTTL {
		return nil, data, fmt.Errorf("unexpected TTL %d (granted %d)", response.TTL, response.GrantedTTL)
	}

	ks := make([]string, len(response.Keys))
	for i := range response.Keys {
		ks[i] = string(response.Keys[i])
	}

	for _, key := range ks {
		data[key] = ""
	}

	return response, data, nil
}
