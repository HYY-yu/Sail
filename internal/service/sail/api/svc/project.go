package svc

import (
	"net/http"

	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/HYY-yu/seckill.pkg/pkg/util"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type ProjectSvc struct {
	BaseSvc
	DB db.Repo

	ProjectRepo      repo.ProjectRepo
	StaffRepo        repo.StaffRepo
	ProjectGroupRepo repo.ProjectGroupRepo
}

func NewProjectSvc(
	db db.Repo,
	projectRepo repo.ProjectRepo,
	pgRepo repo.ProjectGroupRepo,
	staffRepo repo.StaffRepo,
) *ProjectSvc {
	svc := &ProjectSvc{
		DB:               db,
		ProjectRepo:      projectRepo,
		ProjectGroupRepo: pgRepo,
		StaffRepo:        staffRepo,
	}
	return svc
}

func (s *ProjectSvc) List(sctx core.SvcContext, pr *page.PageRequest) (*page.Page, error) {
	ctx := sctx.Context()
	mgr := s.ProjectRepo.Mgr(ctx, s.DB.GetDb(ctx))
	pgMgr := s.ProjectGroupRepo.Mgr(ctx, s.DB.GetDb(ctx))

	projectGroupArr, role := s.CheckStaffGroup(ctx, 0)
	if len(projectGroupArr) == 0 && role != model.RoleAdmin {
		return page.NewPage(
			0,
			[]model.ProjectList{},
		), nil
	}

	limit, offset := pr.GetLimitAndOffset()
	pr.AddAllowSortField(model.ProjectColumns.CreateTime)
	sort, _ := pr.Sort()

	op := make([]repo.Option, 0)
	if v, ok := pr.Filter["project_id"]; ok && util.IsNotZero(v) {
		op = append(op, mgr.WithID(gconv.Int(v)))
	}
	if v, ok := pr.Filter["project_name"]; ok && util.IsNotZero(v) {
		op = append(op, mgr.WithName(
			util.WrapSqlLike(gconv.String(v)),
			" LIKE ?",
		))
	}
	op = append(op, mgr.WithDeleteTime(0))
	if role > model.RoleAdmin {
		op = append(op, mgr.WithProjectGroupID(projectGroupArr, " IN ?"))
	}

	data, err := mgr.WithOptions(op...).ListProject(limit, offset, sort)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	var count int64
	mgr.Count(&count)
	pgMgr.UpdateDB(pgMgr.WithPrepareStmt())

	result := make([]model.ProjectList, len(data))
	for i, e := range data {
		pg, _ := pgMgr.WithOptions(pgMgr.WithID(e.ProjectGroupID)).WithSelects(
			model.ProjectGroupColumns.ID,
			model.ProjectGroupColumns.Name,
		).Get()
		_, mr := s.CheckStaffGroup(ctx, pg.ID)

		r := model.ProjectList{
			ProjectID:        e.ID,
			ProjectGroupID:   e.ProjectGroupID,
			ProjectGroupName: pg.Name,
			Key:              e.Key,
			Name:             e.Name,
			CreateBy:         e.CreateBy,
			CreateByName:     s.GetCreateByName(ctx, s.DB, s.StaffRepo, e.CreateBy),
			Managed:          mr <= model.RoleOwner,
		}
		result[i] = r
	}
	return page.NewPage(
		count,
		result,
	), nil
}
