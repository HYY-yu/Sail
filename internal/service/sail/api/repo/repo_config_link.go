package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type ConfigLinkRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_ConfigLinkMgr
}

type configLinkRepo struct {
}

func NewConfigLinkRepo() ConfigLinkRepo {
	return &configLinkRepo{}
}

func (*configLinkRepo) Mgr(ctx context.Context, db *gorm.DB) *_ConfigLinkMgr {
	mgr := ConfigLinkMgr(ctx, db)
	return mgr
}

// ------- 自定义方法 -------

func (obj *_ConfigLinkMgr) ListConfigLink(
	limit, offset int,
	sort string,
) (result []model.ConfigLink, err error) {
	err = obj.
		sort(sort, model.ConfigLinkColumns.ID+" desc").
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}
