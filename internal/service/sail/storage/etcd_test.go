package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/client/v3/concurrency"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_etcdRepo_ConcurrentSet(t *testing.T) {
	e, err := New(&ETCDConfig{
		Endpoints: []string{"127.0.0.1:2379", "127.0.0.1:12379", "127.0.0.1:22379"},
	})
	assert.NoError(t, err)

	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TEST1",
			args: args{
				key:   "/TEST",
				value: "-1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := sync.WaitGroup{}
			resultC := make(chan int, 1)

			for i := 0; i < 10; i++ {
				go func(x int) {
					wg.Add(1)
					defer func() {
						wg.Done()
					}()

					sr := e.ConcurrentSet(context.Background(), tt.args.key, strconv.Itoa(x))
					if sr.Err == concurrency.ErrLocked {
						return
					}
					resultC <- x
				}(i)
			}
			wg.Wait()
			r := <-resultC
			t.Logf("write is %d", r)

			gr := e.Get(context.Background(), tt.args.key)
			assert.Equal(t, strconv.Itoa(r), gr.Value)

			// 清理
			err = e.Del(context.Background(), tt.args.key)
			assert.NoError(t, err)
		})
	}
}

func Test_etcdRepo_AtomicBatchSet(t *testing.T) {
	e, err := New(&ETCDConfig{
		Endpoints: []string{"127.0.0.1:2379", "127.0.0.1:12379", "127.0.0.1:22379"},
	})
	assert.NoError(t, err)

	type args struct {
		timeout   time.Duration
		sleepTime time.Duration
		keys      []string
		values    []string
	}
	tests := []struct {
		name       string
		args       args
		wantHasKey bool
	}{
		// 完成 3 个 key 要 3 秒，ctx 1 秒后过期，此事务无法成功，查询 ETCD 肯定没有 key
		{
			name: "TEST1",
			args: args{
				timeout:   time.Second,
				sleepTime: time.Second,
				keys:      []string{"/TEST1", "/TEST2", "/TEST3"},
				values:    []string{"1", "2", "3"},
			},
			wantHasKey: false,
		},
		//  完成 3 个 key 要 3 秒，ctx 5 秒后过期，此事务肯定成功，查询 ETCD 肯定有 key
		{
			name: "TEST2",
			args: args{
				timeout:   time.Second * 5,
				sleepTime: time.Second,
				keys:      []string{"/TEST1", "/TEST2", "/TEST3"},
				values:    []string{"1", "2", "3"},
			},
			wantHasKey: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.args.timeout)
			_ = cancel

			e.AtomicBatchSet(ctx, tt.args.keys, tt.args.values, func() {
				time.Sleep(tt.args.sleepTime)
			})

			gp := e.Get(context.Background(), tt.args.keys[0])
			if tt.wantHasKey {
				assert.NoError(t, gp.Err)
				assert.Equal(t, tt.args.values[0], gp.Value)
			} else {
				assert.Equal(t, GetResponse{}, gp)
			}

			// 清理
			for _, key := range tt.args.keys {
				err := e.Del(context.Background(), key)
				assert.NoError(t, err)
			}
		})
	}
}
