package repo

import (
	"context"

	"gorm.io/gorm"
)

type StaffGroupRelRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_StaffGroupRelMgr
}

type staffGroupRelRepo struct {
}

func NewStaffGroupRelRepo() StaffGroupRelRepo {
	return &staffGroupRelRepo{}
}

func (*staffGroupRelRepo) Mgr(ctx context.Context, db *gorm.DB) *_StaffGroupRelMgr {
	mgr := StaffGroupRelMgr(ctx, db)
	return mgr
}

// ------- 自定义方法 -------
