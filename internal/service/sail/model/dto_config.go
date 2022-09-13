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
	Content      string `json:"content"`
	IsPublic     bool   `json:"is_public"`
	IsLinkPublic bool   `json:"is_link_public"`

	IsCopy    bool `json:"is_copy"`
	IsEncrypt bool `json:"is_encrypt"`
}

type AddConfig struct {
	Name         string `json:"name"`
	ProjectID    int    `json:"project_id"`
	NamespaceID  int    `json:"namespace_id"` // 传-1代表全部
	IsPublic     bool   `json:"is_public"`
	IsLinkPublic bool   `json:"is_link_public"`

	IsEncrypt bool   `json:"is_encrypt"`
	Type      string `json:"type"`
	Content   string `json:"content"`

	PublicConfigID int `json:"public_config_id"`
}

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
