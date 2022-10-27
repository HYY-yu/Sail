package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type ConfigHistoryRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_ConfigHistoryMgr
}

type configHistoryRepo struct {
}

func NewConfigHistoryRepo() ConfigHistoryRepo {
	return &configHistoryRepo{}
}

func (*configHistoryRepo) Mgr(ctx context.Context, db *gorm.DB) *_ConfigHistoryMgr {
	mgr := ConfigHistoryMgr(ctx, db)
	return mgr
}

// ------- 自定义方法 -------

func (obj *_ConfigHistoryMgr) ListConfigHistory(
	limit, offset int,
	sort string,
) (result []model.ConfigHistory, err error) {
	err = obj.
		sort(sort, model.ConfigHistoryColumns.ID+" desc").
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}
