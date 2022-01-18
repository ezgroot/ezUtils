package lock

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

const (
	defaultLockTTL         = 3
	defaultLockOwnerPrefix = "/lock/ownerOf"
	defaultLockIDPrefix    = "/lock/implOf"

	sigOfCloseLocker = 1
	sigOfSessionDone = 2
)

// listenLockQuit 监听关闭锁和锁的会话异常状态，进行对应处理
func (l *Locker) listenLockQuit() {
	for {
		sig := <-l.quitCh

		switch sig {
		case sigOfCloseLocker:
			{
				fmt.Printf("%s lock close\n", l.lockID)
				l.clear(context.Background())
				return
			}
		case sigOfSessionDone:
			{
				fmt.Printf("%s session done\n", l.lockID)
				l.clear(context.Background())
				break
			}
		}
	}
}

// impl 执行抢占锁过程
func (l *Locker) impl(ctx context.Context) error {
	if l.client == nil {
		return fmt.Errorf("%s etcd client is nil", l.lockName)
	}

	ctxTemp, cancelTemp := context.WithTimeout(ctx, time.Duration(1500)*time.Millisecond)
	defer cancelTemp()

	_, lease, err := CreateLease(ctxTemp, l.client, int64(l.ttl))
	if err != nil {
		return fmt.Errorf("%s create lease error = %s", l.lockName, err)
	}

	ctxAlways, cancel := context.WithCancel(ctx)

	session, err := concurrency.NewSession(l.client, concurrency.WithTTL(l.ttl),
		concurrency.WithLease(lease), concurrency.WithContext(ctxAlways))
	if err != nil {
		cancel()
		return fmt.Errorf("%s create session error = %s", l.lockName, err)
	}

	mutex := concurrency.NewMutex(session, l.lockName)
	err = mutex.Lock(ctxAlways)
	if err != nil {
		session.Close()
		cancel()
		return fmt.Errorf("%s lock error = %s", l.lockName, err)
	}

	if l.ownerID != "" {
		ctxTemp2, cancelTemp2 := context.WithTimeout(ctx, time.Duration(1500)*time.Millisecond)
		defer cancelTemp2()
		_, err = PutWithLease(
			ctxTemp2,
			l.client,
			fmt.Sprintf("%s/%s/%x", defaultLockOwnerPrefix, l.originalName, lease),
			l.ownerID,
			lease,
		)
		if err != nil {
			mutex.Unlock(ctx)
			session.Close()
			cancel()
			return fmt.Errorf("%s set owner error = %s", l.lockName, err)
		}
	}

	l.setMutex.Lock()
	defer l.setMutex.Unlock()

	l.leaseID = lease
	l.cancel = cancel
	l.session = session
	l.mutexLock = mutex

	l.lockID = l.mutexLock.Key()
	l.myRev = l.mutexLock.Header().Revision

	// 监听锁的会话状态并通知
	go func() {
		for {
			<-l.session.Done()
			l.quitCh <- sigOfSessionDone
			return
		}
	}()

	return nil
}

func (l *Locker) clear(ctx context.Context) {
	fmt.Printf("clear begin !!!\n")
	l.setMutex.Lock()
	defer l.setMutex.Unlock()

	ctxTemp, cancel := context.WithTimeout(ctx, time.Duration(1500)*time.Millisecond)
	defer cancel()

	if l.mutexLock != nil {
		err := l.mutexLock.Unlock(ctxTemp)
		if err != nil {
			fmt.Printf("%s mutex unlock error = %s\n", l.lockID, err)
		}
	}

	if l.session != nil {
		err := l.session.Close()
		if err != nil {
			// NOTE: 此处告警信息不用关注，因为上面已经调用unlock，所以此处错误context canceled
			// 必然出现，只是为确保会话结束而调用session.Close()
			// fmt.Printf( "%s session close error = %s\n", l.lockID, err)
		}
	}

	if l.cancel != nil {
		l.cancel()
	}

	l.leaseID = clientv3.NoLease
	l.cancel = nil
	l.session = nil
	l.mutexLock = nil

	l.lockID = ""
	l.myRev = 0

	fmt.Printf("clear finished !!!\n")
}

func (l *Locker) lock(ctx context.Context) error {
	err := l.impl(ctx)
	if err != nil {
		return err
	}

	isLocked, revision, err := GetRemoteLockState(ctx, l.client, l.lockName)
	if err != nil {
		fmt.Printf("%s get remote state error = %s\n", l.lockName, err)
		return err
	}

	fmt.Printf("isLocked = %t,revision = %d, l.myRev = %d\n", isLocked, revision, l.myRev)

	if isLocked {
		// 自己持有锁
		if l.myRev == revision {
			return nil
		}

		//　他人持有锁
		l.clear(context.Background())
		return ErrLocked
	}

	// 上锁失败，异常
	l.clear(context.Background())

	return ErrLockFailed
}

func (l *Locker) checkLockState(ctx context.Context) error {
	isLocked, revision, err := GetRemoteLockState(ctx, l.client, l.lockName)
	if err != nil {
		fmt.Printf("%s get remote state error = %s\n", l.lockName, err)
		return err
	}

	if isLocked {
		if l.myRev == revision {
			return nil
		}

		l.clear(ctx)
		return concurrency.ErrLocked
	}

	return ErrLockRemoteNoneNeedTodo
}
