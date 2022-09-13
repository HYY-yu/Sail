package svc

import (
	"context"

	"github.com/HYY-yu/seckill.pkg/db"

	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type BaseSvc struct {
}

func (s *BaseSvc) CheckStaffGroup(ctx context.Context, projectGroupId int) ([]int, model.Role) {
	//return nil, model.RoleAdmin
	data := ctx.Value(model.StaffGroupRelKey)
	sgArr, ok := data.([]model.StaffGroup)
	if !ok {
		return nil, 0
	}

	pgArr := make([]int, len(sgArr))
	var resultRole model.Role = 99
	for i, e := range sgArr {
		pgArr[i] = e.ProjectGroupID

		if e.ProjectGroupID == projectGroupId {
			resultRole = e.Role
		}
		if e.Role == model.RoleAdmin {
			return nil, model.RoleAdmin
		}
	}
	return pgArr, resultRole
}

func (s *BaseSvc) GetCreateByName(ctx context.Context, db db.Repo, staffRepo repo.StaffRepo, createBy int) string {
	mgr := staffRepo.Mgr(ctx, db.GetDb(ctx))

	staff, _ := mgr.WithOptions(mgr.WithID(createBy)).
		WithSelects(
			model.StaffColumns.ID,
			model.StaffColumns.Name,
		).Get()
	return staff.Name
}
