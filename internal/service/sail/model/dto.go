package model

type ProjectGroupList struct {
	ProjectGroupID int    `json:"project_group_id"`
	Name           string `json:"name"`
	CreateBy       int    `json:"create_by"`
	CreateByName   string `json:"create_by_name"`
	CreateTime     int64  `json:"create_time"`
}

type AddProjectGroup struct {
	Name string `json:"name" v:"required|length:3,10"`
}

type EditProjectGroup struct {
	ProjectGroupID int     `json:"project_group_id" v:"required"`
	Name           *string `json:"name" v:"length:3,10"`
}

type ProjectList struct {
	ProjectID        int    `json:"project_id"`
	ProjectGroupID   int    `json:"project_group_id"`
	ProjectGroupName string `json:"project_group_name"`
	Key              string `json:"key"`
	Name             string `json:"name"`
	CreateBy         int    `json:"create_by"`
	CreateByName     string `json:"create_by_name"`
	CreateTime       int64  `json:"create_time"`
}

type AddProject struct {
	Name           string `json:"name" v:"required|length:3,10"`
	Key            string `json:"key" v:"required|regex:^[a-zA-Z][\\w-_.]{1,9}"`
	ProjectGroupID int    `json:"project_group_id" v:"required"`
}

type EditProject struct {
	ProjectId int     `json:"project_id" v:"required"`
	Name      *string `json:"name" v:"length:3,10"`
}

type NamespaceList struct {
	NamespaceID      int    `json:"namespace_id"`
	ProjectGroupID   int    `json:"project_group_id"`
	ProjectGroupName string `json:"project_group_name"`

	Name         string `json:"name"`
	RealTime     bool   `json:"real_time"` // 是否灰度
	CreateBy     int    `json:"create_by"`
	CreateByName string `json:"create_by_name"`
	CreateTime   int64  `json:"create_time"`
}

type AddNamespace struct {
	ProjectGroupID int    `json:"project_group_id" v:"required"`
	Name           string `json:"name" v:"required|regex:^[a-zA-Z]{1,9}"`
	RealTime       bool   `json:"real_time"` // 是否灰度
}

type EditNamespace struct {
	NamespaceId int     `json:"namespace_id" v:"required"`
	Name        *string `json:"name" v:"regex:^[a-zA-Z]{1,9}"`
	RealTime    *bool   `json:"real_time"`
}

type StaffList struct {
	StaffID    int    `json:"staff_id"`
	Name       string `json:"name"`
	CreateTime int64  `json:"create_time"`

	Roles []StaffRole `json:"roles"`
}

type Role int

const (
	RoleAdmin Role = iota + 1
	RoleOwner
	RoleManager
)

func (r Role) String() string {
	return [...]string{"Unknown", "Admin", "Owner", "Manager"}[r]
}

type StaffRole struct {
	StaffGroupRelID  int    `json:"staff_group_rel_id"`
	ProjectGroupID   int    `json:"project_group_id"`
	ProjectGroupName string `json:"project_group_name"`
	Role             Role   `json:"role"`
	RoleInfo         string `json:"role_info"`
}

type AddStaff struct {
	Name string `json:"name" v:"required|regex:^[a-zA-Z][a-zA-Z1-9]{1,9}"` // 员工标识
}

type EditStaff struct {
	StaffID int     `json:"staff_id" v:"required"`
	Name    *string `json:"name" v:"regex:^[a-zA-Z][a-zA-Z1-9]{1,9}"`
}

type GrantStaff struct {
	StaffID        int  `json:"staff_id" v:"required"`
	Role           Role `json:"role" v:"required"`
	ProjectGroupID int  `json:"project_group_id"`
}

type StaffGroup struct {
	ProjectGroupID int  `json:"project_group_id"`
	Role           Role `json:"role"`
}

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

type PublishConfigList struct {
	PublishID     int    `json:"publish_id"`
	NamespaceID   int    `json:"namespace_id"`
	NamespaceName string `json:"namespace_name"`

	PublishType    int    `json:"publish_type"`
	PublishTypeStr string `json:"publish_type_str"`

	PublishConfigs []string `json:"publish_configs"`

	CreateBy     int    `json:"create_by"`
	CreateByName string `json:"create_by_name"`
	CreateTime   int64  `json:"create_time"`

	Status    int    `json:"status"`
	StatusStr string `json:"status_str"`
}

type RollbackPublish struct {
	PublishID int `json:"publish_id"`
}

type AddPublish struct {
	ProjectID   int `json:"project_id"`
	NamespaceID int `json:"namespace_id"`

	ConfigIDArr []int  `json:"config_id_arr"`
	PublishType int    `json:"publish_type"`
	PublishData string `json:"publish_data"`
}

type LoginParams struct {
	UserName string `json:"user_name" v:"required|length:3,20"`
	Password string `json:"password" v:"required|length:6,30"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	InitPassword bool   `json:"init_password"`
}
