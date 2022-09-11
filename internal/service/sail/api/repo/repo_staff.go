package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type StaffRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_StaffMgr
}

type staffRepo struct {
}

func NewStaffRepo() StaffRepo {
	return &staffRepo{}
}

func (*staffRepo) Mgr(ctx context.Context, db *gorm.DB) *_StaffMgr {
	mgr := StaffMgr(ctx, db)
	return mgr
}

// ------- 自定义方法 -------

func (obj *_StaffMgr) ListStaff(
	limit, offset int,
	sort string,
) (result []model.Staff, err error) {
	err = obj.
		sort(sort, model.StaffColumns.ID+" desc").
		WithContext(obj.ctx).
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}
