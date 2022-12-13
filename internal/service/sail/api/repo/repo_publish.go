package repo

import (
	"context"
	"github.com/HYY-yu/sail/internal/service/sail/model"
	"gorm.io/gorm"
)

type PublishRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_PublishMgr
}

type publishRepo struct {
}

func NewPublishRepo() PublishRepo {
	return &publishRepo{}
}

func (*publishRepo) Mgr(ctx context.Context, db *gorm.DB) *_PublishMgr {
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
