package config_server

type MetaConfig struct {
	// Namespace 限制 ConfigServer 管理的 namespace 范围，专注于单个 namespace
	// 这里的 namespace 和 ConfigMapRequest 的 namespace 字段不是一个概念
	// 此处 namespace 指的是 kubernetes 的 namespace
	Namespace string `json:"namespace"`

	ETCDEndpoints string `json:"etcd_endpoints" toml:"endpoints"` // 分号分隔的ETCD地址，0.0.0.0:2379;0.0.0.0:12379;0.0.0.0:22379
	ETCDUsername  string `json:"etcd_username" toml:"etcd_username"`
	ETCDPassword  string `json:"etcd_password" toml:"etcd_password"`
}
