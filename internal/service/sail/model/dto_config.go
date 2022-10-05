package model

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/gogf/gf/v2/encoding/gini"
	"github.com/gogf/gf/v2/encoding/gproperties"
	"github.com/gogf/gf/v2/encoding/gtoml"
	"github.com/gogf/gf/v2/encoding/gxml"
	"github.com/gogf/gf/v2/encoding/gyaml"
)

type ProjectTree struct {
	NamespaceID int    `json:"namespace_id"`
	Name        string `json:"name"`
	RealTime    bool   `json:"real_time"`  // 是否需发布
	CanSecret   bool   `json:"can_secret"` // 是否能加密

	Title  string `json:"title"`
	Spread bool   `json:"spread"`

	Nodes []ConfigNode `json:"children"`
}

type ConfigNode struct {
	ConfigID int    `json:"config_id"`
	Title    string `json:"title"`
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

	IsEncrypt bool          `json:"is_encrypt"`
	Type      ConfigType    `json:"type" v:"required"`
	Content   ConfigContent `json:"content" `

	PublicConfigID int `json:"public_config_id"`
}

type ConfigContent string

func (c ConfigContent) Valid(t ConfigType) bool {
	// 为空不校验
	if len(strings.TrimSpace(string(c))) == 0 {
		return true
	}
	cb := []byte(c)

	switch t {
	case ConfigTypeToml:
		_, err := gtoml.Decode(cb)
		return err == nil
	case ConfigTypeJson:
		return json.Valid(cb)
	case ConfigTypeYaml:
		_, err := gyaml.Decode(cb)
		return err == nil
	case ConfigTypeIni:
		_, err := gini.Decode(cb)
		return err == nil
	case ConfigTypeXml:
		_, err := gxml.Decode(cb)
		return err == nil
	case ConfigTypeProperties:
		_, err := gproperties.Decode(cb)
		return err == nil
	default:
		return true
	}
}

var ErrNotEncryptNamespace = errors.New("ErrNotEncryptNamespace")

type ConfigType string

const (
	ConfigTypeCustom     ConfigType = "custom"
	ConfigTypeToml       ConfigType = "toml"
	ConfigTypeYaml       ConfigType = "yaml"
	ConfigTypeJson       ConfigType = "json"
	ConfigTypeIni        ConfigType = "ini"
	ConfigTypeXml        ConfigType = "xml"
	ConfigTypeProperties ConfigType = "properties"
)

func (c ConfigType) AllConfigType() []ConfigType {
	return []ConfigType{ConfigTypeCustom, ConfigTypeToml, ConfigTypeYaml, ConfigTypeJson, ConfigTypeIni, ConfigTypeXml, ConfigTypeProperties}
}

func (c ConfigType) Valid() bool {
	for _, e := range []ConfigType{
		ConfigTypeCustom,
		ConfigTypeToml,
		ConfigTypeYaml,
		ConfigTypeJson,
		ConfigTypeIni,
		ConfigTypeXml,
		ConfigTypeProperties,
	} {
		if e == c {
			return true
		}
	}
	return false
}

type EditConfig struct {
	ConfigID int           `json:"config_id" v:"required"`
	Content  ConfigContent `json:"content" v:"required"`
}

type ConfigCopy struct {
	ConfigID int `json:"config_id" v:"required"`
	Op       int `json:"op" v:"required"` // 1 转为副本 2关联公共配置
}

type ConfigHistoryList struct {
	ConfigID int `json:"config_id" `

	CreateBy     int    `json:"create_by"`
	CreateByName string `json:"create_by_name"`
	CreateTime   int64  `json:"create_time"`
	Reversion    int    `json:"reversion"`
	OpType       int    `json:"op_type"`
	OpTypeStr    string `json:"op_type_str"`
}

type ConfigHistoryOpType int

const (
	OpTypeAdd ConfigHistoryOpType = iota + 1
	OpTypeEdit
	OpTypeRollback
	OpTypeLink
)

func (o ConfigHistoryOpType) String() string {
	return [...]string{"未知", "新增", "编辑", "回滚", "关联"}[o]
}

type RollbackConfig struct {
	ConfigID  int `json:"config_id" v:"required"`
	Reversion int `json:"reversion" v:"required"`
}
