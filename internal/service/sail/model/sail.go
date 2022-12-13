package model

import (
	"gorm.io/datatypes"
	"time"
)

// Config [...]
type Config struct {
	ID             int    `gorm:"primaryKey;column:id;type:int(11);not null"`
	Name           string `gorm:"uniqueIndex:project_id;column:name;type:varchar(50);not null"`
	ProjectID      int    `gorm:"uniqueIndex:project_id;column:project_id;type:int(11);not null"`
	ProjectGroupID int    `gorm:"uniqueIndex:project_id;column:project_group_id;type:int(11);not null"` // 公共配置只有project_group_id
	NamespaceID    int    `gorm:"uniqueIndex:project_id;column:namespace_id;type:int(11);not null"`
	IsPublic       bool   `gorm:"column:is_public;type:tinyint(1);not null"`
	IsLinkPublic   bool   `gorm:"column:is_link_public;type:tinyint(1);not null"`
	IsEncrypt      bool   `gorm:"column:is_encrypt;type:tinyint(1);not null"`
	ConfigType     string `gorm:"uniqueIndex:project_id;column:config_type;type:varchar(10);not null"`
}

// ConfigColumns get sql column name.获取数据库列名
var ConfigColumns = struct {
	ID             string
	Name           string
	ProjectID      string
	ProjectGroupID string
	NamespaceID    string
	IsPublic       string
	IsLinkPublic   string
	IsEncrypt      string
	ConfigType     string
}{
	ID:             "id",
	Name:           "name",
	ProjectID:      "project_id",
	ProjectGroupID: "project_group_id",
	NamespaceID:    "namespace_id",
	IsPublic:       "is_public",
	IsLinkPublic:   "is_link_public",
	IsEncrypt:      "is_encrypt",
	ConfigType:     "config_type",
}

// ConfigHistory [...]
type ConfigHistory struct {
	ID         int       `gorm:"primaryKey;column:id;type:int(11);not null"`
	ConfigID   int       `gorm:"uniqueIndex:config_id;column:config_id;type:int(11);not null"`
	Reversion  int       `gorm:"uniqueIndex:config_id;column:reversion;type:int(11);not null"`
	OpType     int8      `gorm:"column:op_type;type:tinyint(4);not null"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	CreateBy   int       `gorm:"column:create_by;type:int(11);not null"`
}

// ConfigHistoryColumns get sql column name.获取数据库列名
var ConfigHistoryColumns = struct {
	ID         string
	ConfigID   string
	Reversion  string
	OpType     string
	CreateTime string
	CreateBy   string
}{
	ID:         "id",
	ConfigID:   "config_id",
	Reversion:  "reversion",
	OpType:     "op_type",
	CreateTime: "create_time",
	CreateBy:   "create_by",
}

// ConfigLink [...]
type ConfigLink struct {
	ID             int `gorm:"primaryKey;column:id;type:int(11);not null"`
	ConfigID       int `gorm:"uniqueIndex:config_id;column:config_id;type:int(11);not null"`
	PublicConfigID int `gorm:"uniqueIndex:config_id;index:public_config_id;column:public_config_id;type:int(11);not null"`
}

// ConfigLinkColumns get sql column name.获取数据库列名
var ConfigLinkColumns = struct {
	ID             string
	ConfigID       string
	PublicConfigID string
}{
	ID:             "id",
	ConfigID:       "config_id",
	PublicConfigID: "public_config_id",
}

// Namespace [...]
type Namespace struct {
	ID             int       `gorm:"primaryKey;column:id;type:int(11);not null"`
	ProjectGroupID int       `gorm:"uniqueIndex:project_group_id;column:project_group_id;type:int(11);not null"`
	Name           string    `gorm:"uniqueIndex:project_group_id;column:name;type:varchar(50);not null"`
	RealTime       bool      `gorm:"column:real_time;type:tinyint(1);not null"` // 是否是实时发布
	SecretKey      string    `gorm:"column:secret_key;type:varchar(100);not null"`
	CreateTime     time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	CreateBy       int       `gorm:"column:create_by;type:int(11);not null"`
	DeleteTime     int       `gorm:"column:delete_time;type:int(11);not null;default:0"`
}

// NamespaceColumns get sql column name.获取数据库列名
var NamespaceColumns = struct {
	ID             string
	ProjectGroupID string
	Name           string
	RealTime       string
	SecretKey      string
	CreateTime     string
	CreateBy       string
	DeleteTime     string
}{
	ID:             "id",
	ProjectGroupID: "project_group_id",
	Name:           "name",
	RealTime:       "real_time",
	SecretKey:      "secret_key",
	CreateTime:     "create_time",
	CreateBy:       "create_by",
	DeleteTime:     "delete_time",
}

// Project [...]
type Project struct {
	ID             int       `gorm:"primaryKey;column:id;type:int(11);not null"`
	ProjectGroupID int       `gorm:"index:project_group_id;column:project_group_id;type:int(11);not null"`
	Key            string    `gorm:"unique;column:key;type:varchar(50);not null"`
	Name           string    `gorm:"column:name;type:varchar(50);not null"`
	CreateTime     time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	CreateBy       int       `gorm:"column:create_by;type:int(11);not null"`
	DeleteTime     int       `gorm:"column:delete_time;type:int(11);not null;default:0"`
}

// ProjectColumns get sql column name.获取数据库列名
var ProjectColumns = struct {
	ID             string
	ProjectGroupID string
	Key            string
	Name           string
	CreateTime     string
	CreateBy       string
	DeleteTime     string
}{
	ID:             "id",
	ProjectGroupID: "project_group_id",
	Key:            "key",
	Name:           "name",
	CreateTime:     "create_time",
	CreateBy:       "create_by",
	DeleteTime:     "delete_time",
}

// ProjectGroup [...]
type ProjectGroup struct {
	ID         int       `gorm:"primaryKey;column:id;type:int(11);not null"`
	Name       string    `gorm:"unique;column:name;type:varchar(50);not null"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	CreateBy   int       `gorm:"column:create_by;type:int(11);not null"`
	DeleteTime int       `gorm:"column:delete_time;type:int(11);not null;default:0"`
}

// ProjectGroupColumns get sql column name.获取数据库列名
var ProjectGroupColumns = struct {
	ID         string
	Name       string
	CreateTime string
	CreateBy   string
	DeleteTime string
}{
	ID:         "id",
	Name:       "name",
	CreateTime: "create_time",
	CreateBy:   "create_by",
	DeleteTime: "delete_time",
}

// Publish [...]
type Publish struct {
	ID          int       `gorm:"primaryKey;column:id;type:int(11);not null"`
	ProjectID   int       `gorm:"index:project_id;column:project_id;type:int(11);not null"`
	NamespaceID int       `gorm:"index:project_id;column:namespace_id;type:int(11);not null"`
	Status      int8      `gorm:"index:project_id;column:status;type:tinyint(4);not null"`
	CreateTime  time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdateTime  time.Time `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

// PublishColumns get sql column name.获取数据库列名
var PublishColumns = struct {
	ID          string
	ProjectID   string
	NamespaceID string
	Status      string
	CreateTime  string
	UpdateTime  string
}{
	ID:          "id",
	ProjectID:   "project_id",
	NamespaceID: "namespace_id",
	Status:      "status",
	CreateTime:  "create_time",
	UpdateTime:  "update_time",
}

// PublishConfig [...]
type PublishConfig struct {
	ID                 int       `gorm:"primaryKey;column:id;type:int(11);not null"`
	PublishID          int       `gorm:"index:publish_id;column:publish_id;type:int(11);not null"`
	ConfigID           int       `gorm:"index:config_id;column:config_id;type:int(11);not null"`
	ConfigPreReversion int       `gorm:"column:config_pre_reversion;type:int(11);not null"`
	Status             int       `gorm:"column:status;type:int(11);not null"`
	CreateTime         time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdateTime         time.Time `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

// PublishConfigColumns get sql column name.获取数据库列名
var PublishConfigColumns = struct {
	ID                 string
	PublishID          string
	ConfigID           string
	ConfigPreReversion string
	Status             string
	CreateTime         string
	UpdateTime         string
}{
	ID:                 "id",
	PublishID:          "publish_id",
	ConfigID:           "config_id",
	ConfigPreReversion: "config_pre_reversion",
	Status:             "status",
	CreateTime:         "create_time",
	UpdateTime:         "update_time",
}

// PublishStrategy [...]
type PublishStrategy struct {
	ID         int            `gorm:"primaryKey;column:id;type:int(11);not null"`
	PublishID  int            `gorm:"index:publish_id;column:publish_id;type:int(11);not null"`
	Type       int8           `gorm:"column:type;type:tinyint(4);not null"` // 发布类型
	Data       datatypes.JSON `gorm:"column:data;type:json;not null"`
	Status     int8           `gorm:"column:status;type:tinyint(4);not null"`
	Result     string         `gorm:"column:result;type:varchar(500);not null"`
	CreateTime time.Time      `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdateTime time.Time      `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

// PublishStrategyColumns get sql column name.获取数据库列名
var PublishStrategyColumns = struct {
	ID         string
	PublishID  string
	Type       string
	Data       string
	Status     string
	Result     string
	CreateTime string
	UpdateTime string
}{
	ID:         "id",
	PublishID:  "publish_id",
	Type:       "type",
	Data:       "data",
	Status:     "status",
	Result:     "result",
	CreateTime: "create_time",
	UpdateTime: "update_time",
}

// Staff [...]
type Staff struct {
	ID           int       `gorm:"primaryKey;column:id;type:int(11);not null"`
	Name         string    `gorm:"unique;column:name;type:varchar(30);not null"`
	Password     string    `gorm:"column:password;type:varchar(100);not null"`
	RefreshToken string    `gorm:"column:refresh_token;type:varchar(200);not null;default:''"`
	CreateTime   time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	CreateBy     int       `gorm:"column:create_by;type:int(11);not null"`
}

// StaffColumns get sql column name.获取数据库列名
var StaffColumns = struct {
	ID           string
	Name         string
	Password     string
	RefreshToken string
	CreateTime   string
	CreateBy     string
}{
	ID:           "id",
	Name:         "name",
	Password:     "password",
	RefreshToken: "refresh_token",
	CreateTime:   "create_time",
	CreateBy:     "create_by",
}

// StaffGroupRel [...]
type StaffGroupRel struct {
	ID             int  `gorm:"primaryKey;column:id;type:int(11);not null"`
	ProjectGroupID int  `gorm:"index:project_group_id;column:project_group_id;type:int(11);not null"`
	StaffID        int  `gorm:"index:staff_id;column:staff_id;type:int(11);not null"`
	RoleType       int8 `gorm:"column:role_type;type:tinyint(4);not null"` // 权限角色
}

// StaffGroupRelColumns get sql column name.获取数据库列名
var StaffGroupRelColumns = struct {
	ID             string
	ProjectGroupID string
	StaffID        string
	RoleType       string
}{
	ID:             "id",
	ProjectGroupID: "project_group_id",
	StaffID:        "staff_id",
	RoleType:       "role_type",
}
