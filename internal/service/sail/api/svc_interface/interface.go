package svc_interface

import (
	"context"
	"github.com/HYY-yu/sail/internal/service/sail/model"
)

// PublishSystem 发布系统
type PublishSystem interface {
	// EnterPublish
	// ConfigSystem 判断本次更新配置是否需要进入发布系统（判断条件：编辑的命名空间是否需要发布），进入发布系统则不走原来的配置编辑逻辑。
	EnterPublish(ctx context.Context, projectID, namespaceID, configID int, content string) error

	ListPublishConfig(ctx context.Context, projectID, namespaceID int) ([]model.PublishConfig, string, error)

	DeletePublish(ctx context.Context, projectID, namespaceID int, newStatus int) error
}

// ConfigSystem 配置系统
type ConfigSystem interface {
	// ConfigEdit 配置变更回调，有历史记录
	// 做一个配置覆盖编辑，如果是回滚，则用发布前版本覆盖
	// 如果是全量发布，则用发布内容覆盖
	ConfigEdit()

	// GetConfig 根据 configID 获取 config
	GetConfig(ctx context.Context, configID int) (*model.Config, error)

	// ConfigKey 获取配置 key 格式
	ConfigKey(isPublic bool, projectGroupID int, projectKey string, namespaceName string, configName string, configType model.ConfigType) string

	// GetConfigProjectAndNamespace 获取 project 和 namespace 的关键信息
	GetConfigProjectAndNamespace(ctx context.Context, projectID int, namespaceID int) (*model.Project, *model.Namespace, error)
}
