package svc_publish

// PublishSystem 发布系统
type PublishSystem interface {
	// EnterPublish
	// ConfigSystem 判断本次更新配置是否需要进入发布系统（判断条件：编辑的命名空间是否需要发布），进入发布系统则不走原来的配置编辑逻辑。
	// 进入发布，如果namespace尚未处于发布期，则自动进入发布期
	// 将 config 加入 namespace 的发布期
	// 如果 config 已加入，则更新 config 内容
	// 如果 config 未加入，则加入
	// 如果 config 已被锁定，则返回无法进入
	EnterPublish(projectID, namespaceID, configID int, content string)

	// QueryPublish 查询配置的状态
	QueryPublish(configID int)
}

// ConfigSystem 配置系统
type ConfigSystem interface {
	// ConfigEdit 配置变更回调
	// 做一个配置覆盖编辑，如果是回滚，则用发布前版本覆盖
	// 如果是全量发布，则用发布内容覆盖
	ConfigEdit()
}
