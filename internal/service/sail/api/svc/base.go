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
	data := ctx.Value(model.StaffGroupRelKey)
	// sgArr 员工的所有 ProjectGroup 对应的权限
	sgArr, ok := data.([]model.StaffGroup)
	if !ok {
		return nil, 0
	}

	// pgArr 返回此员工拥有的所有 ProjectGroup
	pgArr := make([]int, len(sgArr))
	// resultRole 返回此员工在 projectGroupId 对应的 ProjectGroup
	// 的具体权限， 99 代表无权限。
	var resultRole model.Role = 99
	for i, e := range sgArr {
		pgArr[i] = e.ProjectGroupID

		if e.ProjectGroupID == projectGroupId {
			resultRole = e.Role
		}
		if e.Role == model.RoleAdmin {
			// 如果此员工拥有 RoleAdmin ，则不管他在任何 ProjectGroup
			// 所有 ProjectGroup 都对他开放
			return nil, model.RoleAdmin
		}
	}
	return pgArr, resultRole
}

func (s *BaseSvc) GetCreateByName(ctx context.Context, db db.Repo, staffRepo repo.StaffRepo, createBy int) string {
	mgr := staffRepo.Mgr(ctx, db.GetDb())
	mgr.WithPrepareStmt()

	staff, _ := mgr.WithOptions(mgr.WithID(createBy)).
		WithSelects(
			model.StaffColumns.ID,
			model.StaffColumns.Name,
		).Get()
	return staff.Name
}
