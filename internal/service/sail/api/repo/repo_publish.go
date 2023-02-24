package repo

import (
	"context"
	"github.com/HYY-yu/sail/internal/service/sail/model"
	"gorm.io/gorm"
)

type PublishRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) PublishMgrInter
}

type publishRepo struct {
}

func NewPublishRepo() PublishRepo {
	return &publishRepo{}
}

func (*publishRepo) Mgr(ctx context.Context, db *gorm.DB) PublishMgrInter {
	mgr := PublishMgr(ctx, db)
	return mgr
}

// ------- 自定义方法 -------

func (obj *_PublishMgr) ListPublish(
	limit, offset int,
	sort string,
) (result []model.Publish, err error) {
	err = obj.
		sort(sort, model.PublishColumns.ID+" desc").
		WithContext(obj.ctx).
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}

// -------- 提取接口 ---------

// PublishMgrInter 用于 Mock Repo 层做单元测试
type PublishMgrInter interface {
	ListPublish(
		limit, offset int,
		sort string,
	) (result []model.Publish, err error)
	WithSelects(idName string, selects ...string) *_PublishMgr
	WithOptions(opts ...Option) *_PublishMgr
	GetTableName() string
	Tx(tx *gorm.DB) *_PublishMgr
	WithPrepareStmt()
	Reset() *_PublishMgr
	Get() (result model.Publish, err error)
	Gets() (results []model.Publish, err error)
	Catch() (results model.Publish, err error)
	Count() (count int64, err error)
	HasRecord() (bool, error)
	WithID(id interface{}, cond ...string) Option
	WithProjectID(projectID interface{}, cond ...string) Option
	WithNamespaceID(namespaceID interface{}, cond ...string) Option
	WithPublishToken(publishToken interface{}, cond ...string) Option
	WithStatus(status interface{}, cond ...string) Option
	WithCreateTime(createTime interface{}, cond ...string) Option
	WithUpdateTime(updateTime interface{}, cond ...string) Option
	CreatePublish(bean *model.Publish) (err error)
	UpdatePublish(bean *model.Publish) (err error)
	DeletePublish(bean *model.Publish) (err error)
	FetchIndexByProjectID(projectID int, namespaceID int, status int8) (results []*model.Publish, err error)
}
