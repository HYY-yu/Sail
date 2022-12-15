package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/client/v3/concurrency"
	"strconv"
	"sync"
	"testing"
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
