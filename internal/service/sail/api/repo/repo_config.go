package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type ConfigRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_ConfigMgr
}

type configRepo struct {
}

func NewConfigRepo() ConfigRepo {
	return &configRepo{}
}

func (*configRepo) Mgr(ctx context.Context, db *gorm.DB) *_ConfigMgr {
	mgr := ConfigMgr(ctx, db)
	return mgr
}

// ------- 自定义方法 -------

func (obj *_ConfigMgr) ListConfig(
	limit, offset int,
	sort string,
) (result []model.Config, err error) {
	err = obj.
		sort(sort, model.ConfigColumns.ID+" desc").
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}
