package repo

import (
	"time"
)

// Config [...]
type Config struct {
	ID             int    `gorm:"primaryKey;column:id;type:int;not null"`
	Name           string `gorm:"column:name;type:varchar(50);not null"`
	ProjectGroupID int    `gorm:"column:project_group_id;type:int;not null"`
	NamespaceID    int    `gorm:"column:namespace_id;type:int;not null"`
	IsPublic       bool   `gorm:"column:is_public;type:tinyint(1);not null"`
	IsLinkPublic   bool   `gorm:"column:is_link_public;type:tinyint(1);not null"`
	IsEncrypt      bool   `gorm:"column:is_encrypt;type:tinyint(1);not null"`
	ConfigType     int    `gorm:"column:config_type;type:int;not null"`
	ConfigKey      string `gorm:"column:config_key;type:varchar(50);not null"`
}

// ConfigColumns get sql column name.获取数据库列名
var ConfigColumns = struct {
	ID             string
	Name           string
	ProjectGroupID string
	NamespaceID    string
	IsPublic       string
	IsLinkPublic   string
	IsEncrypt      string
	ConfigType     string
	ConfigKey      string
}{
	ID:             "id",
	Name:           "name",
	ProjectGroupID: "project_group_id",
	NamespaceID:    "namespace_id",
	IsPublic:       "is_public",
	IsLinkPublic:   "is_link_public",
	IsEncrypt:      "is_encrypt",
	ConfigType:     "config_type",
	ConfigKey:      "config_key",
}

// ConfigHistory [...]
type ConfigHistory struct {
	ID         int       `gorm:"primaryKey;column:id;type:int;not null"`
	ConfigID   int       `gorm:"column:config_id;type:int;not null"`
	Reversion  int       `gorm:"column:reversion;type:int;not null"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;not null"`
	CreateBy   int       `gorm:"column:create_by;type:int;not null"`
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
	ID             int `gorm:"primaryKey;column:id;type:int;not null"`
	ConfigID       int `gorm:"column:config_id;type:int;not null"`
	PublicConfigID int `gorm:"column:public_config_id;type:int;not null"`
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
	ID             int       `gorm:"primaryKey;column:id;type:int;not null"`
	ProjectGroupID int       `gorm:"column:project_group_id;type:int;not null"`
	Name           string    `gorm:"column:name;type:varchar(50);not null"`
	RealTime       bool      `gorm:"column:real_time;type:tinyint(1);not null"` // 是否是实时发布
	CreateTime     time.Time `gorm:"column:create_time;type:timestamp;not null"`
	CreateBy       int       `gorm:"column:create_by;type:int;not null"`
	DeleteTime     int       `gorm:"column:delete_time;type:int;not null;default:0"`
}

// NamespaceColumns get sql column name.获取数据库列名
var NamespaceColumns = struct {
	ID             string
	ProjectGroupID string
	Name           string
	RealTime       string
	CreateTime     string
	CreateBy       string
	DeleteTime     string
}{
	ID:             "id",
	ProjectGroupID: "project_group_id",
	Name:           "name",
	RealTime:       "real_time",
	CreateTime:     "create_time",
	CreateBy:       "create_by",
	DeleteTime:     "delete_time",
}

// Project [...]
type Project struct {
	ID             int       `gorm:"primaryKey;column:id;type:int;not null"`
	ProjectGroupID int       `gorm:"column:project_group_id;type:int;not null"`
	Key            string    `gorm:"column:key;type:varchar(50);not null"`
	Name           string    `gorm:"column:name;type:varchar(50);not null"`
	CreateTime     time.Time `gorm:"column:create_time;type:timestamp;not null"`
	CreateBy       int       `gorm:"column:create_by;type:int;not null"`
	DeleteTime     int       `gorm:"column:delete_time;type:int;not null;default:0"`
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
	ID         int       `gorm:"primaryKey;column:id;type:int;not null"`
	Name       string    `gorm:"column:name;type:varchar(50);not null"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;not null"`
	CreateBy   int       `gorm:"column:create_by;type:int;not null"`
	DeleteTime int       `gorm:"column:delete_time;type:int;not null;default:0"`
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
	ID          int       `gorm:"primaryKey;column:id;type:int;not null"`
	NamespaceID int       `gorm:"column:namespace_id;type:int;not null"`
	Status      int       `gorm:"column:status;type:int;not null"`
	CreateTime  time.Time `gorm:"column:create_time;type:timestamp;not null"`
	UpdateTime  time.Time `gorm:"column:update_time;type:timestamp;not null"`
}

// PublishConfigColumns get sql column name.获取数据库列名
var PublishConfigColumns = struct {
	ID          string
	NamespaceID string
	Status      string
	CreateTime  string
	UpdateTime  string
}{
	ID:          "id",
	NamespaceID: "namespace_id",
	Status:      "status",
	CreateTime:  "create_time",
	UpdateTime:  "update_time",
}

// Staff [...]
type Staff struct {
	ID           int       `gorm:"primaryKey;column:id;type:int;not null"`
	Name         string    `gorm:"column:name;type:varchar(30);not null"`
	StaffGroupID int       `gorm:"column:staff_group_id;type:int;not null"`
	RoleType     int       `gorm:"column:role_type;type:int;not null"` // 权限角色
	CreateTime   time.Time `gorm:"column:create_time;type:timestamp;not null"`
	CreateBy     int       `gorm:"column:create_by;type:int;not null"`
	DeleteTime   int       `gorm:"column:delete_time;type:int;not null;default:0"`
}

// StaffColumns get sql column name.获取数据库列名
var StaffColumns = struct {
	ID           string
	Name         string
	StaffGroupID string
	RoleType     string
	CreateTime   string
	CreateBy     string
	DeleteTime   string
}{
	ID:           "id",
	Name:         "name",
	StaffGroupID: "staff_group_id",
	RoleType:     "role_type",
	CreateTime:   "create_time",
	CreateBy:     "create_by",
	DeleteTime:   "delete_time",
}

// StaffGroup [...]
type StaffGroup struct {
	ID         int       `gorm:"primaryKey;column:id;type:int;not null"`
	Name       string    `gorm:"column:name;type:varchar(50);not null"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;not null"`
	CreateBy   int       `gorm:"column:create_by;type:int;not null"`
	DeleteTime int       `gorm:"column:delete_time;type:int;not null;default:0"`
}

// StaffGroupColumns get sql column name.获取数据库列名
var StaffGroupColumns = struct {
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

// StaffGroupProject [...]
type StaffGroupProject struct {
	ID           int `gorm:"primaryKey;column:id;type:int;not null"`
	StaffGroupID int `gorm:"column:staff_group_id;type:int;not null"`
	ProjectID    int `gorm:"column:project_id;type:int;not null"`
}

// StaffGroupProjectColumns get sql column name.获取数据库列名
var StaffGroupProjectColumns = struct {
	ID           string
	StaffGroupID string
	ProjectID    string
}{
	ID:           "id",
	StaffGroupID: "staff_group_id",
	ProjectID:    "project_id",
}

// StaffGroupRel [...]
type StaffGroupRel struct {
	ID           int `gorm:"primaryKey;column:id;type:int;not null"`
	StaffGroupID int `gorm:"column:staff_group_id;type:int;not null"`
	StaffID      int `gorm:"column:staff_id;type:int;not null"`
}

// StaffGroupRelColumns get sql column name.获取数据库列名
var StaffGroupRelColumns = struct {
	ID           string
	StaffGroupID string
	StaffID      string
}{
	ID:           "id",
	StaffGroupID: "staff_group_id",
	StaffID:      "staff_id",
}
