package config_server

import (
	"context"
	"github.com/HYY-yu/sail/internal/operator/api/v1beta1"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sort"
)

type EtcdWatcher struct {
	s *configServer

	namespaceSecret string

	ctx context.Context

	spec *v1beta1.ConfigMapRequestSpec

	watching      bool
	stopWatchChan chan struct{}

	managedConfigMap map[ConfigKey]ConfigValue
}

func NewWatcher(
	ctx context.Context,
	s *configServer,
	namespaceSecret string,
	spec *v1beta1.ConfigMapRequestSpec,
	cm map[ConfigKey]ConfigValue,
) *EtcdWatcher {

	etcdW := &EtcdWatcher{
		s:                s,
		ctx:              ctx,
		namespaceSecret:  namespaceSecret,
		spec:             spec,
		managedConfigMap: cm,
		stopWatchChan:    make(chan struct{}),
	}

	if *spec.Watch {
		etcdW.Run()
	}
	return etcdW
}

func (e EtcdWatcher) ManagedConfigs() []ConfigKey {
	result := make([]ConfigKey, 0)
	for k := range e.managedConfigMap {
		result = append(result, k)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].String() < result[j].String()
	})

	return result
}

func (e *EtcdWatcher) Watching() bool {
	return e.watching
}

func (e *EtcdWatcher) ShouldWatch(w bool) {
	switch w {
	case true:
		if e.watching {
			return
		}
		// start new watch
		close(e.stopWatchChan)
		e.stopWatchChan = make(chan struct{})

		e.Run()
	case false:
		if !e.watching {
			return
		}

		// stop old watch
		close(e.stopWatchChan)
	}
}

func (e *EtcdWatcher) Run() {
	if e.s.etcdClient.Watcher == nil {
		return
	}

	e.s.l.V(1).Info("start etcd watch")

	managedConfigKeys := e.ManagedConfigs()
	keyPrefix := getETCDKeyPrefix(e.spec)

	fromKey := managedConfigKeys[0]
	wc := e.s.etcdClient.Watch(
		e.ctx,
		keyPrefix+fromKey.String(),
		clientv3.WithFromKey(),
	)

	e.watching = true
	go func() {
		defer func() {
			e.s.l.V(1).Info("close etcd watch, bye~ ")
			e.watching = false
		}()
		for {
			select {
			case we := <-wc:
				for _, ev := range we.Events {
					switch ev.Type {
					case mvccpb.PUT:
						// 过滤不监听的 key
						if _, ok := e.managedConfigMap[ConfigKey(getETCDKeyPrefix(e.spec)+string(ev.Kv.Key))]; !ok {
							continue
						}

						isPublish, _ := checkPublish(ev.Kv.Value)
						if isPublish {
							// 在 Watch 配置变更时，忽略 Publish 消息推送
							continue
						}

						e.dealETCDMsg(string(ev.Kv.Key), ev.Kv.Value)
					case mvccpb.DELETE:
						// 删除时，忽略此事件，防止产生不好的影响。
						// 程序若要删除 ConfigMap，应做好处理，然后删除 CMR
					}
				}
			case <-e.ctx.Done():
				return
			case <-e.stopWatchChan:
				return
			}
		}
	}()
}

func (e *EtcdWatcher) dealETCDMsg(key string, value []byte) {
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
