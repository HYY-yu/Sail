package model

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

const (
	PublishStatusRelease = iota + 1 // 发布期
	PublishStatusLock               // 锁定期
	PublishStatusEnd                // 已结束
)
