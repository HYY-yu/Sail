package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type PublishConfigRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_PublishConfigMgr
}

type publishConfigRepo struct {
}

func NewPublishConfigRepo() PublishConfigRepo {
	return &publishConfigRepo{}
}

func (*publishConfigRepo) Mgr(ctx context.Context, db *gorm.DB) *_PublishConfigMgr {
	mgr := PublishConfigMgr(ctx, db)
	return mgr
}

// ------- 自定义方法 -------

func (obj *_PublishConfigMgr) ListPublishConfig(
	limit, offset int,
	sort string,
) (result []model.PublishConfig, err error) {
	err = obj.
		sort(sort, model.PublishConfigColumns.ID+" desc").
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}
