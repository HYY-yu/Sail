package storage

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/errors/gerror"
	"go.etcd.io/etcd/client/v3/concurrency"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ETCDConfig struct {
	Endpoints            []string
	DialTimeout          time.Duration
	DialKeepAlive        time.Duration
	DialKeepAliveTimeout time.Duration
	Username             string
	Password             string
}

func New(cfg *ETCDConfig) (Repo, error) {
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
		Endpoints:            cfg.Endpoints,
		AutoSyncInterval:     time.Minute,
		DialTimeout:          cfg.DialTimeout,
		DialKeepAliveTime:    cfg.DialKeepAlive,
		DialKeepAliveTimeout: cfg.DialKeepAliveTimeout,
		Username:             cfg.Username,
		Password:             cfg.Password,
		PermitWithoutStream:  true,
	})
}

type etcdRepo struct {
	client *clientv3.Client
}

func (e *etcdRepo) Set(ctx context.Context, key string, value string) SetResponse {
	var result SetResponse
	resp, err := e.client.Put(ctx, key, value)
	if err != nil {
		result.Err = err
		return result
	}
	revision := resp.Header.GetRevision()
	result.Revision = int(revision)

	return result
}

// AtomicBatchSet 保证 kv 数组写入是事务的
func (e *etcdRepo) AtomicBatchSet(ctx context.Context, key []string, value []string, callback ...func()) SetResponse {
	if len(key) != len(value) {
		return SetResponse{
			Err: errors.New("key len must equal value len"),
		}
	}
	txRes, err := concurrency.NewSTM(e.client, func(stm concurrency.STM) error {
		for i, k := range key {
			select {
			case <-ctx.Done():
				// ctx 到期，事务中断
				return gerror.New("ctx is done, end stm. ")
			default:
			}

			v := value[i]
			stm.Put(k, v)

			if len(callback) > 0 {
				callback[0]()
			}
		}
		return nil
	})

	sr := SetResponse{
		Err: err,
	}
	if txRes != nil && txRes.Header != nil {
		sr.Revision = int(txRes.Header.GetRevision())
	}
	return sr
}

const ConcurrentSet = "/SAIL/ConcurrentSet"

// ConcurrentSet 保证同时写入只有一个能写入成功
func (e *etcdRepo) ConcurrentSet(ctx context.Context, key string, value string) SetResponse {
	session, _ := concurrency.NewSession(e.client, concurrency.WithTTL(5))
	mux := concurrency.NewMutex(session, ConcurrentSet)

	ctx, c := context.WithTimeout(ctx, time.Second*5)
	_ = c // 消除警告
	err := mux.TryLock(ctx)

	res := SetResponse{}
	defer func() {
		_ = mux.Unlock(ctx)
	}()

	if err != nil {
		res.Err = err
		if err == concurrency.ErrLocked {
			// 如果 err 是 Locked，外部调用通过判断 res.Err 得知自己是否更新成功
		}
		return res
	}

	return e.Set(ctx, key, value)
}

func (e *etcdRepo) Get(ctx context.Context, key string) GetResponse {
	var result GetResponse

	resp, err := e.client.Get(ctx, key)
	if err != nil {
		result.Err = err
		return result
	}

	if len(resp.Kvs) == 0 {
		return result
	}
	result.Value = string(resp.Kvs[0].Value)
	result.Revision = int(resp.Header.GetRevision())
	return result
}

func (e *etcdRepo) GetWithReversion(ctx context.Context, key string, reversion int) GetResponse {
	var result GetResponse

	resp, err := e.client.Get(ctx, key, clientv3.WithRev(int64(reversion)))
	if err != nil {
		result.Err = err
		return result
	}

	if len(resp.Kvs) == 0 {
		return result
	}
	result.Value = string(resp.Kvs[0].Value)
	result.Revision = int(resp.Header.GetRevision())
	return result
}

func (e *etcdRepo) Del(ctx context.Context, key string) error {
	_, err := e.client.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (e *etcdRepo) Close() error {
	if e.client != nil {
		return e.client.Close()
	}
	return nil
}
