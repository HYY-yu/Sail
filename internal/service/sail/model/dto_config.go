package model

type ProjectTree struct {
	Nodes []TreeNode `json:"nodes"`
}

type TreeNode struct {
	ID   int    `json:"id"`
	Type int    `json:"type"` // 可能是 Namespace，可能是 Config
	Name string `json:"name"`

	RealTime bool       `json:"real_time"` // false 代表可发布
	Children []TreeNode `json:"children"`
}

type ConfigInfo struct {
	ConfigID     int    `json:"config_id"`
	ConfigKey    string `json:"config_key"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Content      string `json:"content"`
	IsPublic     bool   `json:"is_public"`
	IsLinkPublic bool   `json:"is_link_public"`

	IsCopy    bool `json:"is_copy"`
	IsEncrypt bool `json:"is_encrypt"`
}

type AddConfig struct {
	Name         string `json:"name" v:"required|regex:^[a-zA-Z][\\w_\\-.]{1,9}"`
	ProjectID    int    `json:"project_id" v:"required"`
	NamespaceID  int    `json:"namespace_id" v:"required"` // 传-1代表全部
	IsPublic     bool   `json:"is_public"`
	IsLinkPublic bool   `json:"is_link_public"`

	IsEncrypt bool   `json:"is_encrypt"`
	Type      string `json:"type" v:"required"`
	Content   string `json:"content"`

	PublicConfigID int `json:"public_config_id"`
}

const (
	ConfigTypeCustom = "custom"
	ConfigTypeToml   = "toml"
	ConfigTypeYaml   = "yaml"
	ConfigTypeJson   = "json"
	ConfigTypeIni    = "ini"
)

type EditConfig struct {
	ConfigID int    `json:"config_id"`
	Content  string `json:"content"`
}

type ConfigCopy struct {
	ConfigID int `json:"config_id"`
	Op       int `json:"op"` // 1 转为副本 2关联公共配置
}

type ConfigHistoryList struct {
	ConfigID int `json:"config_id"`

	CreateBy     int    `json:"create_by"`
	CreateByName string `json:"create_by_name"`
	CreateTime   int64  `json:"create_time"`
	Reversion    int    `json:"reversion"`
	Content      string `json:"content"`
}

type RollbackConfig struct {
	HistoryID int `json:"history_id"`
}
