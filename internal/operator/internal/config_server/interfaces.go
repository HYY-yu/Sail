package config_server

import (
	"context"
	"fmt"
	"github.com/HYY-yu/sail/internal/operator/api/v1beta1"
	"github.com/HYY-yu/seckill.pkg/pkg/encrypt"
	"github.com/go-logr/logr"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"path/filepath"
	"sort"
	"strconv"
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
	InitAndWatch(ctx context.Context, namespaceSecretKey string, resp *v1beta1.ConfigMapRequestSpec) error

	// Get 检查配置是否成功下载到 ConfigMap，Watch 连接是否正常
	Get(ctx context.Context, spec *v1beta1.ConfigMapRequestSpec)
}

type configServer struct {
	l         logr.Logger
	namespace string

	etcdClient *clientv3.Client
	clientSet  *kubernetes.Clientset

	metaConfig MetaConfig
	restConfig *rest.Config
}

func NewConfigServer(l logr.Logger, restConfig *rest.Config, metaConfig MetaConfig) ConfigServer {
	return &configServer{
		l: l.WithName("ConfigServer"),

		metaConfig: metaConfig,
		namespace:  metaConfig.Namespace,
		restConfig: restConfig,
	}
}

func (c *configServer) etcdConnect() (*clientv3.Client, error) {
	c.l.V(1).Info("start to connect etcd. ")
	etcdEndpoints := strings.Split(c.metaConfig.ETCDEndpoints, ";")

	v3cfg := &clientv3.Config{
		Endpoints:            etcdEndpoints,
		AutoSyncInterval:     time.Minute,
		DialTimeout:          10 * time.Second,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 20 * time.Second,
		Username:             c.metaConfig.ETCDUsername,
		Password:             c.metaConfig.ETCDPassword,
		PermitWithoutStream:  true,
		DialOptions:          []grpc.DialOption{grpc.WithBlock()},
	}

	return clientv3.New(*v3cfg)
}

func (c *configServer) Start(_ context.Context) error {
	etcdClient, err := c.etcdConnect()
	if err != nil {
		return err
	}
	c.etcdClient = etcdClient

	// Connect Kubernetes
	cs, err := kubernetes.NewForConfig(c.restConfig)
	if err != nil {
		return err
	}

	c.clientSet = cs
	return nil
}

func (c *configServer) InitAndWatch(ctx context.Context, namespaceSecretKey string, spec *v1beta1.ConfigMapRequestSpec) error {
	etcdConfigMap, err := c.pullETCDConfig(ctx, namespaceSecretKey, spec)
	if err != nil {
		return err
	}

	_ = etcdConfigMap // TODO

	// 写入 kubernetes

	c.clientSet.CoreV1().ConfigMaps(c.namespace).Create()

	if *spec.Watched {
		c.l.V(1).Info("start to watch config. ")
		NewWatcher(ctx, c, namespaceSecretKey, spec).Run()
	}
	return nil
}

type ConfigKey string
type ConfigValue []byte

func (c *configServer) pullETCDConfig(ctx context.Context, namespaceSecretKey string, spec *v1beta1.ConfigMapRequestSpec) (map[ConfigKey]ConfigValue, error) {
	keyPrefix := getETCDKeyPrefix(spec)
	c.l.V(1).Info("pull config key", "keys", spec.Configs)

	if len(spec.Configs) == 0 {
		// 取所有 config
		getResp, err := c.etcdClient.Get(ctx,
			keyPrefix,
			clientv3.WithPrefix(),
			clientv3.WithKeysOnly(),
		)
		if err != nil {
			return nil, fmt.Errorf("read config from etcd err: %w ", err)
		}
		etcdKeys := make([]string, 0, len(getResp.Kvs))
		for _, e := range getResp.Kvs {
			etcdKeys = append(etcdKeys, getConfigFileKeyFrom(string(e.Key)))
		}
		spec.Configs = etcdKeys
	}
	sort.Strings(spec.Configs)
	if len(spec.Configs) == 0 {
		// 还是没有 config，直接退出
		return nil, fmt.Errorf("no config found. ")
	}
	fromKey := spec.Configs[0]

	getResp, err := c.etcdClient.Get(ctx,
		keyPrefix+fromKey,
		clientv3.WithFromKey(),
		clientv3.WithLimit(int64(len(spec.Configs))),
	)
	if err != nil {
		return nil, fmt.Errorf("read config from etcd err: %w ", err)
	}
	etcdKeys := make([]string, 0, len(spec.Configs))
	for _, e := range getResp.Kvs {
		etcdKeys = append(etcdKeys, getConfigFileKeyFrom(string(e.Key)))
	}
	if len(etcdKeys) == 0 {
		return nil, fmt.Errorf("read empty config from etcd! ")
	}

	// insETCDKeys 取 etcdKeys 和 spec.Configs 的交集
	// 这是因为 etcdKeys 可能会有些配置被删除了，这时 CMR 尚未更新
	// 我们需要处理这种情况
	insETCDKeys := intersectionSortStringArr(etcdKeys, spec.Configs)
	c.l.V(1).Info("real config key", "keys", insETCDKeys)
	insETCDKeyMap := make(map[string]struct{})
	for _, e := range insETCDKeys {
		insETCDKeyMap[e] = struct{}{}
	}

	result := make(map[ConfigKey]ConfigValue)
	for _, e := range getResp.Kvs {
		configFileKey := getConfigFileKeyFrom(string(e.Key))
		if _, ok := insETCDKeyMap[configFileKey]; ok {
			isPublish, reversion := checkPublish(e.Value)
			if isPublish {
				newValue, err := c.readFromReversion(ctx, e.Key, int64(reversion))
				if err != nil {
					return nil, err
				}
				e.Value = newValue
			}

			// 尝试解密内容
			newValue := c.tryDecryptConfigContent(string(e.Value), namespaceSecretKey)

			// Add ConfigKey and Value
			result[ConfigKey(configFileKey)] = ConfigValue(newValue)
		}
	}
	return result, nil
}

func checkPublish(etcdValue []byte) (isPublish bool, reversion int) {
	etcdValueStr := string(etcdValue)

	if strings.HasPrefix(etcdValueStr, "PUBLISH") {
		publishStrArr := strings.Split(etcdValueStr, "&")

		if len(publishStrArr) != 5 {
			return false, 0
		}

		reversion, _ := strconv.Atoi(publishStrArr[3])
		return true, reversion
	}
	return false, 0
}

func (c *configServer) readFromReversion(ctx context.Context, etcdKey []byte, reversion int64) ([]byte, error) {
	getResp, err := c.etcdClient.Get(ctx,
		string(etcdKey),
		clientv3.WithRev(reversion),
	)
	if err != nil {
		return nil, err
	}
	if len(getResp.Kvs) == 0 {
		return nil, fmt.Errorf("cannot find etcdKey: %s", string(etcdKey))
	}

	return getResp.Kvs[0].Value, nil
}

func (c *configServer) Get(ctx context.Context, spec *v1beta1.ConfigMapRequestSpec) {
	// configs 和 ConfigMap 一一对应
	// merged: true 全部的 configs 对应一个 ConfigMap

	// 先检查 ConfigMap 是否存在

	//
}

// /conf/{project_key}/namespace/config_name.config.type
func getETCDKeyPrefix(spec *v1beta1.ConfigMapRequestSpec) string {
	b := strings.Builder{}

	b.WriteString("/conf")

	b.WriteByte('/')
	b.WriteString(spec.ProjectKey)

	b.WriteByte('/')
	b.WriteString(spec.Namespace)

	b.WriteByte('/')
	return b.String()
}

func getETCDKey(keyPrefix string, configFile string) string {
	return keyPrefix + "/" + configFile
}

func getConfigFileKeyFrom(etcdKey string) string {
	_, result := filepath.Split(etcdKey)
	return result
}

func intersectionSortStringArr(a []string, b []string) []string {
	if (len(a) == 0) || (len(b) == 0) {
		return []string{}
	}

	result := make([]string, 0, len(a))
	i, j := 0, 0
	for i != len(a) && j != len(b) {
		if a[i] > b[j] {
			j++
		} else if a[i] < b[j] {
			i++
		} else {
			result = append(result, a[i])
			i++
			j++
		}
	}
	return result
}

// tryDecryptConfigContent 只是会尝试解密内容，解密失败了就把 content 返回
func (s *configServer) tryDecryptConfigContent(content string, namespaceSecretKey string) string {
	_, err := encrypt.NewBase64Encoding().DecodeString(content)
	if err == nil {
		// 能被 Base64 解码，却不能被解密，那就把 content 原样返回
		decryptContent, err := decryptConfigContent(content, namespaceSecretKey)
		if err == nil {
			content = decryptContent
		}
	}
	return content
}

func decryptConfigContent(content string, namespaceKey string) (string, error) {
	if namespaceKey == "" {
		return "", nil
	}

	goAES := encrypt.NewGoAES(namespaceKey, encrypt.AES192)
	decryptContent, err := goAES.WithModel(encrypt.ECB).WithEncoding(encrypt.NewBase64Encoding()).Decrypt(content)
	if err != nil {
		return "", err
	}
	return decryptContent, nil
}
