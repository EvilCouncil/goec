package etcdlib

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.etcd.io/etcd/client/v3"
)

type EtcdAgent struct {
	Client  *clientv3.Client
	Lease   clientv3.Lease
	LeaseID clientv3.LeaseID
}

func NewEtcdAgent(ctx context.Context, servers []string) (*EtcdAgent, error) {
	ea := new(EtcdAgent)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   servers,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return ea, err
	}

	lease := clientv3.NewLease(cli)
	lr, err := lease.Grant(ctx, 10)
	if err != nil {
		return ea, err
	}

	ea.Client = cli
	ea.Lease = lease
	ea.LeaseID = lr.ID

	return ea, nil
}

func (ea *EtcdAgent) KeepAlive(ctx context.Context) error {
	ka, err := ea.Lease.KeepAlive(ctx, ea.LeaseID)
	if err != nil {
		return err
	}

	running := true
	for running {
		select {
		case r := <-ka:
			slog.Info(fmt.Sprintf("KeepAlive: %d, ttl: %d", r.ID, r.TTL))
		case <-ctx.Done():
			running = false
		}
	}

	return nil
}

func (ea *EtcdAgent) Put(ctx context.Context, key string, value string) error {
	_, err := ea.Client.Put(ctx, key, value, clientv3.WithLease(ea.LeaseID))
	if err != nil {
		return err
	}

	return nil
}
