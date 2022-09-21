package model

import (
	"errors"
)

type ProjectTree struct {
	NamespaceID int    `json:"namespace_id"`
	Name        string `json:"name"`
	RealTime    bool   `json:"real_time"` // 是否需发布

	Nodes []ConfigNode `json:"nodes"`
}

type ConfigNode struct {
	ConfigID int    `json:"config_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

type ConfigInfo struct {
	ConfigID     int    `json:"config_id"`
	ConfigKey    string `json:"config_key"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Content      string `json:"content"`
	IsPublic     bool   `json:"is_public"`
	IsLinkPublic bool   `json:"is_link_public"`
	IsEncrypt    bool   `json:"is_encrypt"`
}

type AddConfig struct {
	Name           string `json:"name" v:"required|regex:^[a-zA-Z][\\w_\\-.]{1,9}"`
	ProjectGroupID int    `json:"project_group_id" v:"required"`
	ProjectID      int    `json:"project_id" ` // 公共配置可以不传projectID
	NamespaceID    int    `json:"namespace_id" v:"required"`
	IsPublic       bool   `json:"is_public"`
	IsLinkPublic   bool   `json:"is_link_public"`

	IsEncrypt bool       `json:"is_encrypt"`
	Type      ConfigType `json:"type" v:"required"`
	Content   string     `json:"content" v:"required"`

	PublicConfigID int `json:"public_config_id"`
}

var ErrNotEncryptNamespace = errors.New("ErrNotEncryptNamespace")

type ConfigType string

const (
	ConfigTypeCustom = "custom"
	ConfigTypeToml   = "toml"
	ConfigTypeYaml   = "yaml"
	ConfigTypeJson   = "json"
	ConfigTypeIni    = "ini"
)

func (c ConfigType) Valid() bool {
	for _, e := range []ConfigType{ConfigTypeCustom, ConfigTypeToml, ConfigTypeYaml, ConfigTypeJson, ConfigTypeIni} {
		if e == c {
			return true
		}
	}
	return false
}

type EditConfig struct {
	ConfigID       int    `json:"config_id"`
	IsPublicConfig bool   `json:"is_public_config"`
	Content        string `json:"content"`
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
