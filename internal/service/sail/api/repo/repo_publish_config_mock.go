package repo

import (
	"context"
	"gorm.io/gorm"
)

type mockPublishRepo struct {
	PM *_PublishMgr // 可以用 gomonkey 替换 PM 的Method
}

func NewMockPublishRepo() PublishRepo {
	return &mockPublishRepo{PM: &_PublishMgr{}}
}

func (m mockPublishRepo) Mgr(ctx context.Context, db *gorm.DB) *_PublishMgr {
	return m.PM
}
