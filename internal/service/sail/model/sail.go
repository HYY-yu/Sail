package model

import (
	"time"
)

// Config [...]
type Config struct {
	ID             int    `gorm:"primaryKey;column:id;type:int(11);not null"`
	Name           string `gorm:"column:name;type:varchar(50);not null"`
	ProjectID      int    `gorm:"column:project_id;type:int(11);not null"`
	ProjectGroupID int    `gorm:"column:project_group_id;type:int(11);not null"` // 公共配置只有project_group_id
	NamespaceID    int    `gorm:"column:namespace_id;type:int(11);not null"`
	IsPublic       bool   `gorm:"column:is_public;type:tinyint(1);not null"`
	IsLinkPublic   bool   `gorm:"column:is_link_public;type:tinyint(1);not null"`
	IsEncrypt      bool   `gorm:"column:is_encrypt;type:tinyint(1);not null"`
	ConfigType     string `gorm:"column:config_type;type:varchar(10);not null"`
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
	ConfigID   int       `gorm:"column:config_id;type:int(11);not null"`
	Reversion  int       `gorm:"column:reversion;type:int(11);not null"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	CreateBy   int       `gorm:"column:create_by;type:int(11);not null"`
}

// ConfigHistoryColumns get sql column name.获取数据库列名
var ConfigHistoryColumns = struct {
	ID         string
	ConfigID   string
	Reversion  string
	CreateTime string
	CreateBy   string
}{
	ID:         "id",
	ConfigID:   "config_id",
	Reversion:  "reversion",
	CreateTime: "create_time",
	CreateBy:   "create_by",
}

// ConfigLink [...]
type ConfigLink struct {
	ID             int `gorm:"primaryKey;column:id;type:int(11);not null"`
	ConfigID       int `gorm:"column:config_id;type:int(11);not null"`
	PublicConfigID int `gorm:"column:public_config_id;type:int(11);not null"`
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
	ProjectGroupID int       `gorm:"column:project_group_id;type:int(11);not null"`
	Name           string    `gorm:"column:name;type:varchar(50);not null"`
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
	ProjectGroupID int       `gorm:"column:project_group_id;type:int(11);not null"`
	Key            string    `gorm:"column:key;type:varchar(100);not null"`
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

// PublishConfig [...]
type PublishConfig struct {
	ID               int       `gorm:"primaryKey;column:id;type:int(11);not null"`
	ProjectID        int       `gorm:"column:project_id;type:int(11);not null"`
	NamespaceID      int       `gorm:"column:namespace_id;type:int(11);not null"`
	PublishType      int       `gorm:"column:publish_type;type:int(11);not null"`     // 发布方式
	PublishData      string    `gorm:"column:publish_data;type:varchar(20);not null"` // 发布数据
	PublishConfigIDs string    `gorm:"column:publish_config_ids;type:varchar(100);not null"`
	Status           int       `gorm:"column:status;type:int(11);not null"`
	CreateTime       time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdateTime       time.Time `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

// PublishConfigColumns get sql column name.获取数据库列名
var PublishConfigColumns = struct {
	ID               string
	ProjectID        string
	NamespaceID      string
	PublishType      string
	PublishData      string
	PublishConfigIDs string
	Status           string
	CreateTime       string
	UpdateTime       string
}{
	ID:               "id",
	ProjectID:        "project_id",
	NamespaceID:      "namespace_id",
	PublishType:      "publish_type",
	PublishData:      "publish_data",
	PublishConfigIDs: "publish_config_ids",
	Status:           "status",
	CreateTime:       "create_time",
	UpdateTime:       "update_time",
}

// Staff [...]
type Staff struct {
	ID           int       `gorm:"primaryKey;column:id;type:int(11);not null"`
	Name         string    `gorm:"column:name;type:varchar(30);not null"`
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
	ID             int `gorm:"primaryKey;column:id;type:int(11);not null"`
	ProjectGroupID int `gorm:"column:project_group_id;type:int(11);not null"`
	StaffID        int `gorm:"column:staff_id;type:int(11);not null"`
	RoleType       int `gorm:"column:role_type;type:int(11);not null"` // 权限角色
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
