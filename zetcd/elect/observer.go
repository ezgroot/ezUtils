package elect

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type Observer struct {
	client    *clientv3.Client
	position  string
	curLeader string
	CurTerm   string
	ttl       int
	session   *concurrency.Session
	election  *concurrency.Election
	cancel    context.CancelFunc
	closeCh   chan struct{}
}

func NewObserver(c *clientv3.Client, p string, ttl int) *Observer {
	o := &Observer{
		client:   c,
		position: p,
		ttl:      ttl,
		closeCh:  make(chan struct{}, 1),
	}

	return o
}

func (o *Observer) Init() error {
	s, err := concurrency.NewSession(o.client)
	if err != nil {
		return err
	}

	e := concurrency.NewElection(s, o.position)

	o.session = s
	o.election = e

	return nil
}

func (o *Observer) Start(ctx context.Context) {
	ctxAlways, cancel := context.WithCancel(ctx)
	o.cancel = cancel

	go func() {
		respCh := o.election.Observe(ctxAlways)

		for {
			select {
			case <-ctxAlways.Done():
				{
					fmt.Printf("ctx done\n")
					return
				}
			case info := <-respCh:
				{
					for _, kv := range info.Kvs {
						o.CurTerm = string(kv.Key)
						o.curLeader = string(kv.Value)
					}
				}
			case <-o.closeCh:
				{
					fmt.Printf("close\n")
					return
				}
			}
		}
	}()
}

func (o *Observer) Stop() {
	o.closeCh <- struct{}{}
	o.cancel()
	o.session.Close()
}
