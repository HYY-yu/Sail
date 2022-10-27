package model

type ProjectGroupList struct {
	ProjectGroupID int    `json:"project_group_id"`
	Name           string `json:"name"`
	CreateBy       int    `json:"create_by"`
	CreateByName   string `json:"create_by_name"`
	CreateTime     int64  `json:"create_time"`

	Managed bool `json:"managed"`
}

type AddProjectGroup struct {
	Name string `json:"name" v:"required|length:3,10"`
}

type EditProjectGroup struct {
	ProjectGroupID int     `json:"project_group_id" v:"required"`
	Name           *string `json:"name" v:"length:3,10"`
}
