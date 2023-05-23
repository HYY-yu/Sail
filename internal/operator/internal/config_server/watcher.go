package config_server

import (
	"context"
	"github.com/HYY-yu/sail/internal/operator/api/v1beta1"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Watcher interface {
	Run()
}

type etcdWatcher struct {
	s *configServer

	namespaceSecret string

	ctx  context.Context
	spec *v1beta1.ConfigMapRequestSpec
}

func NewWatcher(ctx context.Context, s *configServer, namespaceSecret string, spec *v1beta1.ConfigMapRequestSpec) Watcher {
	etcdW := &etcdWatcher{
		s:               s,
		ctx:             ctx,
		namespaceSecret: namespaceSecret,
		spec:            spec,
	}
	return etcdW
}

func (e *etcdWatcher) Run() {
	if e.s.etcdClient.Watcher == nil {
		return
	}

	// Only Watch e.spec.ProjectKey/e.spec.Namespace (Config)
	wc := e.s.etcdClient.Watch(
		e.ctx,
		getETCDKeyPrefix(e.spec),
		clientv3.WithPrefix(),
	)

	go func() {
		for {
			select {
			case we := <-wc:
				for _, ev := range we.Events {
					switch ev.Type {
					case mvccpb.PUT:
						isPublish, _ := checkPublish(ev.Kv.Value)
						if isPublish {
							// 在 Watch 配置变更时，忽略 Publish 消息推送
							continue
						}

						e.dealETCDMsg(string(ev.Kv.Key), ev.Kv.Value)
					case mvccpb.DELETE:

					}
				}
			case <-e.ctx.Done():
				e.s.l.Info("close etcd watch, bye~ ")
				return
			}
		}
	}()
}

func (e *etcdWatcher) dealETCDMsg(key string, value []byte) {
	e.s.l.V(1).Info("got a event by: ", "key", key)
	if len(value) == 0 {
		return
	}
	// TODO 检查 K8s 是否有这个 ConfigMap的 key，没有则忽略

	configFileKey := getConfigFileKeyFrom(key)
	newValue := e.s.tryDecryptConfigContent(string(value), e.namespaceSecret)

	_ = configFileKey // TODO
	_ = newValue
}
