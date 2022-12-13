package repo

import (
	"context"
	"github.com/HYY-yu/sail/internal/service/sail/model"
	"gorm.io/gorm"
)

type PublishStrategyRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_PublishStrategyMgr
}

type publishStrategyRepo struct {
}

func NewPublishStrategyRepo() PublishStrategyRepo {
	return &publishStrategyRepo{}
}

func (*publishStrategyRepo) Mgr(ctx context.Context, db *gorm.DB) *_PublishStrategyMgr {
	mgr := PublishStrategyMgr(ctx, db)
	return mgr
}

// ------- 自定义方法 -------

func (obj *_PublishStrategyMgr) ListPublishStrategy(
	limit, offset int,
	sort string,
) (result []model.PublishStrategy, err error) {
	err = obj.
		sort(sort, model.PublishStrategyColumns.ID+" desc").
		WithContext(obj.ctx).
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}
