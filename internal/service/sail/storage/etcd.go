package storage

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ETCDConfig struct {
	Endpoints   []string
	DialTimeout time.Duration
	Username    string
	Password    string
}

func New(cfg *ETCDConfig) (*etcdRepo, error) {
	client, err := etcdConnect(cfg)
	if err != nil {
		return nil, err
	}

	return &etcdRepo{
		client: client,
	}, nil
}

func etcdConnect(cfg *ETCDConfig) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:        cfg.Endpoints,
		AutoSyncInterval: time.Minute,
		DialTimeout:      cfg.DialTimeout,
		Username:         cfg.Username,
		Password:         cfg.Password,
	})
}

type etcdRepo struct {
	client *clientv3.Client
}

func (e *etcdRepo) Set(ctx context.Context, key string, value string) (int, error) {
	resp, err := e.client.Put(ctx, key, value)
	if err != nil {
		return 0, err
	}
	revision := resp.Header.GetRevision()
	return int(revision), nil
}

func (e *etcdRepo) Get(ctx context.Context, key string) {
	e.client.Get(ctx, key)
}

func (e *etcdRepo) GetWithReversion(ctx context.Context, key string, reversion int) {
	panic("implement me")
}

func (e *etcdRepo) Del(ctx context.Context, key string) {
	panic("implement me")
}

func (e *etcdRepo) Close() error {
	if e.client != nil {
		return e.client.Close()
	}
	return nil
}
