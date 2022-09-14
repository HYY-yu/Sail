package model

type NamespaceList struct {
	NamespaceID      int    `json:"namespace_id"`
	ProjectGroupID   int    `json:"project_group_id"`
	ProjectGroupName string `json:"project_group_name"`

	Name         string `json:"name"`
	RealTime     bool   `json:"real_time"` // 是否灰度
	SecretKey    string `json:"secret_key"`
	CreateBy     int    `json:"create_by"`
	CreateByName string `json:"create_by_name"`
	CreateTime   int64  `json:"create_time"`
}

type AddNamespace struct {
	ProjectGroupID int    `json:"project_group_id" v:"required"`
	Name           string `json:"name" v:"required|regex:^[a-zA-Z][\\w_-.]{1,9}"`
	RealTime       bool   `json:"real_time"` // 是否灰度
	Secret         bool   `json:"secret"`    // 是否加密
}

type EditNamespace struct {
	NamespaceId int     `json:"namespace_id" v:"required"`
	Name        *string `json:"name" v:"regex:^[a-zA-Z][\\w_-.]{1,9}"`
	RealTime    *bool   `json:"real_time"`
}
