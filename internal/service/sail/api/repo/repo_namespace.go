package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type NamespaceRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_NamespaceMgr
}

type namespaceRepo struct {
}

func NewNamespaceRepo() NamespaceRepo {
	return &namespaceRepo{}
}

func (*namespaceRepo) Mgr(ctx context.Context, db *gorm.DB) *_NamespaceMgr {
	mgr := NamespaceMgr(ctx, db)
	return mgr
}

// ------- 自定义方法 -------

func (obj *_NamespaceMgr) ListNamespace(
	limit, offset int,
	sort string,
) (result []model.Namespace, err error) {
	err = obj.
		sort(sort, model.NamespaceColumns.ID+" desc").
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}
