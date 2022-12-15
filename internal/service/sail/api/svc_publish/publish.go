package svc_publish

import (
	"fmt"
	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
	"github.com/HYY-yu/sail/internal/service/sail/storage"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/encrypt"
)

// PublishSystem 发布系统
type PublishSystem interface {
	// EnterPublish
	// ConfigSystem 判断本次更新配置是否需要进入发布系统（判断条件：编辑的命名空间是否需要发布），进入发布系统则不走原来的配置编辑逻辑。
	// 进入发布，如果namespace尚未处于发布期，则自动进入发布期
	// 将 config 加入 namespace 的发布期
	// 如果 config 已加入，则更新 config 内容
	// 如果 config 未加入，则加入
	// 如果 config 已被锁定，则返回无法进入
	EnterPublish(projectID, namespaceID, configID int, content string) error

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

type PublishSvc struct {
	DB    db.Repo
	Store storage.Repo

	ConfigRepo    repo.ConfigRepo
	ProjectRepo   repo.ProjectRepo
	NamespaceRepo repo.NamespaceRepo
}

func NewPublishSvc(
	db db.Repo,
	store storage.Repo,
	cr repo.ConfigRepo,
	pr repo.ProjectRepo,
	nr repo.NamespaceRepo,
) *PublishSvc {
	svc := &PublishSvc{
		DB:            db,
		Store:         store,
		ConfigRepo:    cr,
		ProjectRepo:   pr,
		NamespaceRepo: nr,
	}
	return svc
}

func (p *PublishSvc) EnterPublish(projectID, namespaceID, configID int, content string) error {

	return nil
}

func (p *PublishSvc) QueryPublish(configID int) {

}

// initPublish 进入发布期
// 1. 生成 token
// 2. 写入 token
// 3. 加密配置，设计加密字符串
// 4. 写入数据库，发布状态为发布期
func (p *PublishSvc) initPublish() {

}

func (p *PublishSvc) queryPublish() {

}

func generatePublishToken(projectID, namespaceID int) string {
	return encrypt.SHA256WithEncoding(fmt.Sprintf("%d-%d-%s", projectID, namespaceID, encrypt.Nonce(5)), encrypt.NewBase32Human())
}
