package elect

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type Candidate struct {
	client   *clientv3.Client
	myName   string
	position string
	ttl      int
	session  *concurrency.Session
	election *concurrency.Election
}

func NewCandidate(client *clientv3.Client, position string, name string, ttl int) *Candidate {
	c := &Candidate{
		client:   client,
		position: position,
		myName:   name,
		ttl:      ttl,
	}

	return c
}

func (c *Candidate) Init() error {
	s, err := concurrency.NewSession(c.client)
	if err != nil {
		return err
	}

	e := concurrency.NewElection(s, c.position)

	c.session = s
	c.election = e

	return nil
}

// Campaign When starting the campaign, ctx only controls the waiting time for the campaign.
// If elected, the status of becoming the leader will always be maintained, and it has nothing to do with ctx.
func (c *Candidate) Campaign(ctx context.Context) error {
	// Campaign In the blocking mode, the execution will continue after the election is successful, timeout or error occurs.
	err := c.election.Campaign(ctx, c.myName)
	if err != nil {
		return err
	}

	return nil
}

// Proclaim If you are the leader, you can change the name of yourself and announce it to its members during this period;
// If you are not the current leader, report an error and return concurrency.ErrElectionNotLeader.
func (c *Candidate) Proclaim(ctx context.Context, newName string) error {
	c.myName = newName

	err := c.election.Proclaim(ctx, newName)
	if err != nil {
		return err
	}

	return nil
}

// Resign If you are the leader, you can voluntarily give up being the leader during this period;
// If you are not the current leader, it will have no effect and will not return err.
// If you call Campaign() later, you can campaign again.
func (c *Candidate) Resign(ctx context.Context) error {
	err := c.election.Resign(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Candidate) MyTerm() string {
	return c.election.Key()
}

func (c *Candidate) MyName() string {
	return c.myName
}

func (c *Candidate) IsCurLeaderAtMe(ctx context.Context) (bool, error) {
	if c.election.Key() == "" {
		return false, nil
	}

	term, _, revision, err := c.GetCurLeaderInfo(ctx)
	if err == nil {
		if c.election.Key() == term && c.election.Rev() == revision {
			return true, nil
		}
	} else if err == concurrency.ErrElectionNoLeader {
		return false, nil
	}

	return false, err
}

// GetCurLeaderInfo Get the information about the current leader in the cluster.
// If no one is elected as the leader in the current cluster, then return concurrency.ErrElectionNoLeader error.
func (c *Candidate) GetCurLeaderInfo(ctx context.Context) (string, string, int64, error) {
	resp, err := c.election.Leader(ctx)
	if err != nil {
		return "", "", 0, err
	}

	var term string
	var leader string
	var createRevision int64
	for _, kv := range resp.Kvs {
		term = string(kv.Key)
		leader = string(kv.Value)
		createRevision = kv.CreateRevision
	}

	return term, leader, createRevision, nil
}

// Close After closing the session, the current election or election will be withdrawn,
// and this instance will no longer be available. Calling Campaign again will report an error.
func (c *Candidate) Close() error {
	return c.session.Close()
}
