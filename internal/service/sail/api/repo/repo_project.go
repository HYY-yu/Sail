package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type ProjectRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_ProjectMgr
}

type projectRepo struct {
}

func NewProjectRepo() ProjectRepo {
	return &projectRepo{}
}

func (*projectRepo) Mgr(ctx context.Context, db *gorm.DB) *_ProjectMgr {
	mgr := ProjectMgr(ctx, db)
	return mgr
}

// ------- 自定义方法 -------

func (obj *_ProjectMgr) ListProject(
	limit, offset int,
	sort string,
) (result []model.Project, err error) {
	err = obj.
		sort(sort, model.ProjectColumns.ID+" desc").
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}
