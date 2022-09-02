package model

type ProjectGroupList struct {
	ProjectGroupID int    `json:"project_group_id"`
	Name           string `json:"name"`
	CreateBy       int    `json:"create_by"`
	CreateTime     int    `json:"create_time"`
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
	CreateTime       int    `json:"create_time"`
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

	Name       string `json:"name"`
	RealTime   bool   `json:"real_time"` // 是否灰度
	CreateBy   int    `json:"create_by"`
	CreateTime int    `json:"create_time"`
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
