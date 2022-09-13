package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type ProjectGroupRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_ProjectGroupMgr
}

type projectGroupRepo struct {
}

func NewProjectGroupRepo() ProjectGroupRepo {
	return &projectGroupRepo{}
}

func (*projectGroupRepo) Mgr(ctx context.Context, db *gorm.DB) *_ProjectGroupMgr {
	mgr := ProjectGroupMgr(ctx, db)
	return mgr
}

// ------- 自定义方法 -------

func (obj *_ProjectGroupMgr) ListProjectGroup(
	limit, offset int,
	sort string,
) (result []model.ProjectGroup, err error) {
	err = obj.
		sort(sort, model.ProjectGroupColumns.ID+" desc").
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}
