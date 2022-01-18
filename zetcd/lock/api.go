package lock

import (
	"context"
	"errors"
	"fmt"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var (
	ErrLockerClosed           = errors.New("locker has closed")
	ErrLocked                 = errors.New("locked by another session")
	ErrSessionExpired         = errors.New("session is expired")
	ErrLockFailed             = errors.New("lock failed")
	ErrLockRemoteNoneNeedTodo = errors.New("remote lock is none, can try to grab")
)

type Locker struct {
	client       *clientv3.Client
	ownerID      string
	ttl          int
	originalName string
	lockName     string

	setMutex     sync.Mutex
	runningMutex sync.RWMutex
	grabMutex    sync.RWMutex
	quitCh       chan int

	leaseID   clientv3.LeaseID
	cancel    context.CancelFunc
	session   *concurrency.Session
	mutexLock *concurrency.Mutex

	isClosed bool
	lockID   string
	myRev    int64
}

// NewLocker Create a new distributed lock instance, note that if it is no longer used, please call Close.
func NewLocker(c *clientv3.Client, owner string, lock string, ttl int) *Locker {
	if ttl <= 0 {
		ttl = defaultLockTTL
	}

	l := &Locker{
		client:       c,
		ownerID:      owner,
		ttl:          ttl,
		originalName: lock,
		lockName:     defaultLockIDPrefix + lock,

		quitCh: make(chan int, 1),

		leaseID:   clientv3.NoLease,
		cancel:    nil,
		session:   nil,
		mutexLock: nil,

		isClosed: false,
	}

	go l.listenLockQuit()

	return l
}

func (l *Locker) Lock(ctx context.Context) error {
	l.runningMutex.RLock()
	if l.isClosed {
		l.runningMutex.RUnlock()
		return ErrLockerClosed
	}
	l.runningMutex.RUnlock()

	l.runningMutex.Lock()
	defer l.runningMutex.Unlock()

	if l.isClosed {
		return ErrLockerClosed
	}

	err := l.checkLockState(ctx)
	fmt.Printf("0000 checkLockState err = %v\n", err)
	if err == nil {
		return nil
	} else if err == ErrLockRemoteNoneNeedTodo {
		l.clear(ctx)
		fmt.Printf("ready to begin lock\n")
		return l.lock(ctx)
	} else if err == concurrency.ErrLocked {
		l.clear(ctx)
		return ErrLocked
	} else {
		// NOTE: this means GetRemoteLockState() error, nothing to do
	}

	fmt.Printf("finished Lock\n")

	return err
}

// Unlock Unlock, continue to Lock
func (l *Locker) Unlock(ctx context.Context) error {
	l.runningMutex.RLock()
	if l.isClosed {
		l.runningMutex.RUnlock()
		return ErrLockerClosed
	}
	l.runningMutex.RUnlock()

	l.runningMutex.Lock()
	defer l.runningMutex.Unlock()

	if l.isClosed {
		return ErrLockerClosed
	}

	l.clear(ctx)

	return nil
}

// Close Close lock, no longer available
func (l *Locker) Close() {
	l.runningMutex.Lock()
	defer l.runningMutex.Unlock()

	l.isClosed = true
	l.quitCh <- sigOfCloseLocker
}

func (l *Locker) GetLockID() string {
	return l.lockID
}

func (l *Locker) GetLockRev() int64 {
	return l.myRev
}
