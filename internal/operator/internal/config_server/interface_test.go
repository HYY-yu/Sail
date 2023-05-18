package config_server

import (
	"context"
	"github.com/HYY-yu/sail/internal/operator/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
)

func Test_intersectionSortStringArr(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test1",
			args: args{
				a: []string{"mysql.toml", "redis.yaml"},
				b: []string{"cfg.json", "mysql.toml", "redis.yaml"},
			},
			want: []string{"mysql.toml", "redis.yaml"},
		},
		{
			name: "test2",
			args: args{
				a: []string{"mysql.toml", "redis.yaml"},
				b: []string{},
			},
			want: []string{},
		},
		{
			name: "test3",
			args: args{
				a: []string{"mysql.toml", "redis.yaml"},
				b: []string{"ca.cert"},
			},
			want: []string{},
		},
		{
			name: "test4",
			args: args{
				a: []string{".bachrc", "mysql.toml", "redis.yaml"},
				b: []string{"mysql.toml", "redis.yaml"},
			},
			want: []string{"mysql.toml", "redis.yaml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intersectionSortStringArr(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intersectionSortStringArr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSail_getETCDKeyPrefix(t *testing.T) {
	type fields struct {
		metaConfig *v1beta1.ConfigMapRequestSpec
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test1",
			fields: fields{
				metaConfig: &v1beta1.ConfigMapRequestSpec{
					ProjectKey: "test_project_key",
					Namespace:  "test",
				},
			},
			want: "/conf/test_project_key/test/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getETCDKeyPrefix(tt.fields.metaConfig); got != tt.want {
				t.Errorf("getETCDKeyPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getConfigFileKeyFrom(t *testing.T) {
	type args struct {
		etcdKey string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test",
			args: args{
				etcdKey: "/conf/project_key/test/mysql.toml",
			},
			want: "mysql.toml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConfigFileKeyFrom(tt.args.etcdKey); got != tt.want {
				t.Errorf("getConfigFileKeyFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSail_checkPublish(t *testing.T) {
	type args struct {
		etcdValue []byte
	}
	tests := []struct {
		name          string
		args          args
		wantIsPublish bool
		wantReversion int
	}{
		{
			name:          "test1",
			args:          struct{ etcdValue []byte }{etcdValue: []byte("PUBLISH&THIS_IS_TOKEN&1&22&SecretData==")},
			wantIsPublish: true,
			wantReversion: 22,
		},
		{
			name:          "test2",
			args:          struct{ etcdValue []byte }{etcdValue: []byte("")},
			wantIsPublish: false,
			wantReversion: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsPublish, gotReversion := checkPublish(tt.args.etcdValue)
			assert.Equalf(t, tt.wantIsPublish, gotIsPublish, "checkPublish(%v)", tt.args.etcdValue)
			assert.Equalf(t, tt.wantReversion, gotReversion, "checkPublish(%v)", tt.args.etcdValue)
		})
	}
}

func TestSail_pullETCDConfig(t *testing.T) {
	tests := []struct {
		name            string
		namespaceSecret string
		configs         []string
		response        *clientv3.GetResponse
		wantResp        map[ConfigKey]ConfigValue
	}{
		{
			name:            "TEST1",
			namespaceSecret: "NTUZNTNQNUKYEL4GP5SGVDV9LEYZAWBD",
			configs:         []string{"mysql.toml", "redis.properties"},
			response: &clientv3.GetResponse{
				Kvs: []*mvccpb.KeyValue{
					{
						Key:   []byte("/conf/test_project_key/test/mysql.toml"),
						Value: []byte("database=\"127.0.0.1:3306\""),
					},
					// redis.properties
					// host=0.0.0.0
					// port=6379
					{
						Key:   []byte("/conf/test_project_key/test/redis.properties"),
						Value: []byte("I9IfkJSBekxeYbQJSX6zQsvZJwlfj3VyZ6RrtRF4LFI="),
					},
				},
			},
			wantResp: map[ConfigKey]ConfigValue{
				"mysql.toml":       ConfigValue("database=\"127.0.0.1:3306\""),
				"redis.properties": ConfigValue("host=0.0.0.0\nport=6379"),
			},
		},
		{
			name:            "TEST2",
			configs:         []string{"mysql.toml"},
			namespaceSecret: "NTUZNTNQNUKYEL4GP5SGVDV9LEYZAWBD",
			response: &clientv3.GetResponse{
				Kvs: []*mvccpb.KeyValue{
					{
						Key:   []byte("/conf/test_project_key/test/mysql.toml"),
						Value: []byte("database=\"127.0.0.1:3306\""),
					},
					{
						Key:   []byte("/conf/test_project_key/test/redis.properties"),
						Value: []byte("I9IfkJSBekxeYbQJSX6zQsvZJwlfj3VyZ6RrtRF4LFI="),
					},
				},
			},
			wantResp: map[ConfigKey]ConfigValue{
				"mysql.toml": ConfigValue("database=\"127.0.0.1:3306\""),
			},
		}, {
			name:            "TEST3",
			configs:         []string{"cfg.custom"},
			namespaceSecret: "",
			response: &clientv3.GetResponse{
				Kvs: []*mvccpb.KeyValue{
					{
						Key:   []byte("/conf/test_project_key/test/cfg.custom"),
						Value: []byte("CA"),
					},
				},
			},
			wantResp: map[ConfigKey]ConfigValue{
				"cfg.custom": ConfigValue("CA"),
			},
		},
	}

	csInter := NewConfigServer(zap.New(), nil, MetaConfig{})
	cs := csInter.(*configServer)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs.etcdClient = &clientv3.Client{
				KV: &mockKV{KV: clientv3.NewKVFromKVClient(nil, nil), response: tt.response},
			}

			mapValue, err := cs.pullETCDConfig(context.Background(), tt.namespaceSecret, &v1beta1.ConfigMapRequestSpec{
				ProjectKey: "test_project_key",
				Namespace:  "test",
				Configs:    tt.configs,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.wantResp, mapValue)
		})
	}
}

type mockKV struct {
	clientv3.KV
	response *clientv3.GetResponse
}

func (kv *mockKV) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	return kv.response, nil
}
