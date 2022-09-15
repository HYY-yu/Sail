package model

type ProjectList struct {
	ProjectID        int    `json:"project_id"`
	ProjectGroupID   int    `json:"project_group_id"`
	ProjectGroupName string `json:"project_group_name"`
	Key              string `json:"key"`
	Name             string `json:"name"`
	CreateBy         int    `json:"create_by"`
	CreateByName     string `json:"create_by_name"`
	CreateTime       int64  `json:"create_time"`

	Managed bool `json:"managed"`
}

type AddProject struct {
	Name           string `json:"name" v:"required|length:3,10"`
	ProjectGroupID int    `json:"project_group_id" v:"required"`
}

type EditProject struct {
	ProjectId int `json:"project_id" v:"required"`

	Name *string `json:"name" v:"length:3,10"`
}
