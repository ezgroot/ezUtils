package zetcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// Locker 分布式锁
type Locker struct {
	client   *clientv3.Client
	ownerID  string
	ttl      int
	lockName string

	cancel    context.CancelFunc
	session   *concurrency.Session
	mutexLock *concurrency.Mutex
}

func NewLocker(client *clientv3.Client, lock string, ttl int, owner string) *Locker {
	l := &Locker{
		client:   client,
		lockName: lock,
		ttl:      ttl,
		ownerID:  owner,
	}

	return l
}

func (l *Locker) Init() error {
	s, err := concurrency.NewSession(l.client, concurrency.WithTTL(l.ttl))
	if err != nil {
		return err
	}

	l.session = s

	mutex := concurrency.NewMutex(l.session, l.lockName)
	l.mutexLock = mutex

	return nil
}

func (l *Locker) TryLock(ctx context.Context) error {
	return l.mutexLock.TryLock(ctx)
}

func (l *Locker) Lock(ctx context.Context) error {
	return l.mutexLock.Lock(ctx)
}

func (l *Locker) IsCurLockAtMe(ctx context.Context) (bool, error) {
	response, data, err := GetFirstCrt(ctx, l.client, l.lockName)
	if err != nil {
		return false, err
	}

	if len(data) == 0 {
		return false, nil
	}

	var curLocker string
	for k := range data {
		curLocker = k
	}

	var curLockerRev int64
	for _, kv := range response.Kvs {
		curLockerRev = kv.CreateRevision
	}

	if curLocker == l.mutexLock.Key() &&
		curLockerRev == l.mutexLock.Header().Revision {
		return true, nil
	}

	return false, nil
}

func (l *Locker) Unlock(ctx context.Context) error {
	return l.mutexLock.Unlock(ctx)
}

func (l *Locker) Close() error {
	return l.session.Close()
}
