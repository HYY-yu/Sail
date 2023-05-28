package config_server

import (
	"context"
	"github.com/HYY-yu/sail/internal/operator/api/v1beta1"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"sync"
	"time"
)

type EtcdWatcher struct {
	s *configServer

	namespaceSecret string

	ctx context.Context

	spec *v1beta1.ConfigMapRequestSpec

	watching      bool
	stopWatchChan chan struct{}

	managedConfigMap            map[ConfigKey]ConfigValue
	managedConfigLastUpdateTime map[ConfigKey]*time.Time

	rwLock sync.RWMutex
}

func NewWatcher(
	ctx context.Context,
	s *configServer,
	namespaceSecret string,
	spec *v1beta1.ConfigMapRequestSpec,
	cm map[ConfigKey]ConfigValue,
) *EtcdWatcher {

	etcdW := &EtcdWatcher{
		s:                           s,
		ctx:                         ctx,
		namespaceSecret:             namespaceSecret,
		spec:                        spec,
		managedConfigMap:            cm,
		managedConfigLastUpdateTime: make(map[ConfigKey]*time.Time),
		stopWatchChan:               make(chan struct{}),
		rwLock:                      sync.RWMutex{},
	}

	// init LastUpdateTime
	nowTime := time.Now()
	for k := range cm {
		etcdW.managedConfigLastUpdateTime[k] = &nowTime
	}

	if *spec.Watch {
		etcdW.Run()
	}
	return etcdW
}

func (e *EtcdWatcher) ManagedConfigKeys() []ConfigKey {
	result := make([]ConfigKey, 0)
	e.rwLock.RLock()
	defer e.rwLock.Unlock()

	for k := range e.managedConfigMap {
		result = append(result, k)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].String() < result[j].String()
	})

	return result
}

type ConfigManagedInfo struct {
	Value          ConfigValue
	LastUpdateTime *time.Time
}

func (e *EtcdWatcher) ManagedConfig() map[ConfigKey]ConfigManagedInfo {
	result := make(map[ConfigKey]ConfigManagedInfo)

	e.rwLock.RLock()
	defer e.rwLock.Unlock()
	for k, v := range e.managedConfigMap {
		result[k] = ConfigManagedInfo{
			Value:          v,
			LastUpdateTime: e.managedConfigLastUpdateTime[k],
		}
	}
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

	managedConfigKeys := e.ManagedConfigKeys()
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
						e.rwLock.RLock()
						configFileKey := getConfigKeyFrom(string(ev.Kv.Key))
						if _, ok := e.managedConfigMap[configFileKey]; !ok {
							return
						}
						e.rwLock.Unlock()

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
	e.s.l.V(1).Info("got a event by etcd ", "key", key)
	if len(value) == 0 {
		return
	}
	var failReturn bool
	defer func() {
		if failReturn {
			e.s.l.Info("failed to deal etcd msg, will try again in next watch event. ")
		}
	}()

	configMapName := getConfigMapName(e.spec)

	_, err := e.s.clientSet.CoreV1().ConfigMaps(e.s.namespace).Get(e.ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		e.s.l.Error(err, "failed to get configMap")
		failReturn = true
		return
	}

	newValue := tryDecryptConfigContent(string(value), e.namespaceSecret)
	configFileKey := getConfigKeyFrom(key)
	nowTime := time.Now()

	// 更新 configMap
	e.rwLock.Lock()
	e.managedConfigMap[configFileKey] = ConfigValue(newValue)
	e.managedConfigLastUpdateTime[configFileKey] = &nowTime
	e.rwLock.Unlock()

	configMapData, err := makeConfigMapData(e.spec, e.managedConfigMap)
	if err != nil {
		e.s.l.Error(err, "failed to make configMap data")
		failReturn = true
		return
	}
	configMap := generateConfigMap(e.spec, configMapName, configMapData)
	_, err = e.s.clientSet.CoreV1().ConfigMaps(e.s.namespace).Update(e.ctx, configMap, metav1.UpdateOptions{})
	return
}
