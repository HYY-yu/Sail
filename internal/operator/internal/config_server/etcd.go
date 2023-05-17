package config_server

type MetaConfig struct {
	ETCDEndpoints string `json:"etcd_endpoints" toml:"endpoints"` // 分号分隔的ETCD地址，0.0.0.0:2379;0.0.0.0:12379;0.0.0.0:22379
	ETCDUsername  string `json:"etcd_username" toml:"etcd_username"`
	ETCDPassword  string `json:"etcd_password" toml:"etcd_password"`
}
