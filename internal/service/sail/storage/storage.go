package storage

import (
	"context"
)

type Repo interface {
	Set(ctx context.Context, key string, value string) SetResponse
	Get(ctx context.Context, key string) GetResponse
	GetWithReversion(ctx context.Context, key string, reversion int) GetResponse
	Del(ctx context.Context, key string) bool
}

type SetResponse struct {
	Revision int
	Err      error
}

type GetResponse struct {
	Value string

	Revision int
}
