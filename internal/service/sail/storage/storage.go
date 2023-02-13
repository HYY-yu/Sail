package storage

import (
	"context"
)

type Repo interface {
	Set(ctx context.Context, key string, value string) SetResponse
	AtomicBatchSet(ctx context.Context, key []string, value []string, callback ...func()) SetResponse
	Get(ctx context.Context, key string) GetResponse
	ConcurrentSet(ctx context.Context, key string, value string) SetResponse
	GetWithReversion(ctx context.Context, key string, reversion int) GetResponse
	Del(ctx context.Context, key string) error
	Close() error
}
type SetResponse struct {
	Revision int
	Err      error
}

type GetResponse struct {
	Value string

	Revision int
	Err      error
}
