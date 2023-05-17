package config_server

import (
	"context"
	"github.com/go-logr/logr"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"strings"
	"time"
)

// ConfigServer 是 Operator 控制 ConfigMap 的重要组成部分，
// 它负责保持与 ETCD 服务器的连接，从 ETCD 获取指定的配置以及持续监听这个配置
// 当配置发生变更，它将自动更新到 ConfigMap。
// 它会在本地维护 ConfigMap 和 ETCD 配置间的关系，以及一份配置缓存，以便 Get 重复获取。
// - 对于 Reconclier 来说，它只需要启动 InitAndWatch 向 ConfigServer 请求相关配置即可。
// - Reconclier 也可以 Get 配置检查是否下载成功，是否正在 Watch，并且检查配置和 ConfigMap 是否对应。
// - 总之，Reconclier把对 ConfigMap 的管理全部委托给 ConfigServer。
type ConfigServer interface {
	// InitAndWatch 用于从配置系统获取配置，写入到 ConfigMap，并 Watch 配置保证配置是最新的
	InitAndWatch()

	// Get 检查配置是否成功下载到 ConfigMap，Watch 连接是否正常
	Get()
}

type configServer struct {
	L logr.Logger

	metaConfig MetaConfig
}

func NewConfigServer(l logr.Logger, metaConfig MetaConfig) ConfigServer {
	return &configServer{
		L:          l.WithName("ConfigServer"),
		metaConfig: metaConfig,
	}
}

func (s *configServer) etcdConnect() (*clientv3.Client, error) {
	s.L.V(1).Info("start to connect etcd. ")
	etcdEndpoints := strings.Split(s.metaConfig.ETCDEndpoints, ";")

	v3cfg := &clientv3.Config{
		Endpoints:            etcdEndpoints,
		AutoSyncInterval:     time.Minute,
		DialTimeout:          10 * time.Second,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 20 * time.Second,
		Username:             s.metaConfig.ETCDUsername,
		Password:             s.metaConfig.ETCDPassword,
		PermitWithoutStream:  true,
		DialOptions:          []grpc.DialOption{grpc.WithBlock()},
	}

	return clientv3.New(*v3cfg)
}

func (c *configServer) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (c *configServer) InitAndWatch() {
	//TODO implement me
	panic("implement me")
}

func (c *configServer) Get() {
	//TODO implement me
	panic("implement me")
}
