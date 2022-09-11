package svc

import (
	"net/http"
	"time"

	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/mysqlerr_helper"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/HYY-yu/seckill.pkg/pkg/util"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type ProjectGroupSvc struct {
	BaseSvc
	DB db.Repo

	PGRepo    repo.ProjectGroupRepo
	StaffRepo repo.StaffRepo
}

func NewProjectGroupSvc(
	db db.Repo,
	pgRepo repo.ProjectGroupRepo,
	staffRepo repo.StaffRepo,
) *ProjectGroupSvc {
	svc := &ProjectGroupSvc{
		DB:        db,
		PGRepo:    pgRepo,
		StaffRepo: staffRepo,
	}
	return svc
}

func (s *ProjectGroupSvc) List(sctx core.SvcContext, pr *page.PageRequest) (*page.Page, error) {
	ctx := sctx.Context()
	mgr := s.PGRepo.Mgr(ctx, s.DB.GetDb(ctx))

	limit, offset := pr.GetLimitAndOffset()
	pr.AddAllowSortField(model.ProjectGroupColumns.CreateTime)
	sort, _ := pr.Sort()

	// Filter
	op := make([]repo.Option, 0)
	if v, ok := pr.Filter["project_group_id"]; ok && util.IsNotZero(v) {
		op = append(op, mgr.WithID(gconv.Int(v)))
	}
	if v, ok := pr.Filter["project_group_name"]; ok && util.IsNotZero(v) {
		op = append(op, mgr.WithName(
			util.WrapSqlLike(gconv.String(v)),
			" LIKE ?",
		))
	}
	op = append(op, mgr.WithDeleteTime(0))

	data, err := mgr.WithOptions(op...).ListProjectGroup(limit, offset, sort)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	var count int64
	mgr.Count(&count)

	var result = make([]model.ProjectGroupList, len(data))

	for i, e := range data {
		r := model.ProjectGroupList{
			ProjectGroupID: e.ID,
			Name:           e.Name,
			CreateBy:       e.CreateBy,
			CreateByName:   s.GetCreateByName(ctx, s.DB, s.StaffRepo, e.CreateBy),
			CreateTime:     e.CreateTime.Unix(),
		}

		result[i] = r
	}
	return page.NewPage(
		count,
		result,
	), nil
}

func (s *ProjectGroupSvc) Add(sctx core.SvcContext, param *model.AddProjectGroup) error {
	ctx := sctx.Context()
	mgr := s.PGRepo.Mgr(ctx, s.DB.GetDb(ctx))

	_, role := s.CheckStaffGroup(ctx, 0)
	if role != model.RoleAdmin {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"只有管理员可以访问此接口",
		)
	}

	bean := &model.ProjectGroup{
		Name:       param.Name,
		CreateBy:   int(sctx.UserId()),
		CreateTime: time.Now(),
	}

	err := mgr.CreateProjectGroup(bean)
	if err != nil {
		if mysqlerr_helper.IsMysqlDupEntryError(err) {
			return response.NewErrorWithStatusOk(
				response.ParamBindError,
				"已存在相同的ProjectGroup",
			)
		}
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (s *ProjectGroupSvc) Edit(sctx core.SvcContext, param *model.EditProjectGroup) error {
	ctx := sctx.Context()
	mgr := s.PGRepo.Mgr(ctx, s.DB.GetDb(ctx))

	_, role := s.CheckStaffGroup(ctx, 0)
	if role != model.RoleAdmin {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"只有管理员可以访问此接口",
		)
	}

	bean := &model.ProjectGroup{
		ID: param.ProjectGroupID,
	}

	updateColumns := make([]string, 0)

	if param.Name != nil && !g.IsEmpty(*param.Name) {
		bean.Name = *param.Name
		updateColumns = append(updateColumns, model.ProjectGroupColumns.Name)
	}

	err := mgr.WithSelects(model.ProjectGroupColumns.ID, updateColumns...).UpdateProjectGroup(bean)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (s *ProjectGroupSvc) Delete(sctx core.SvcContext, projectGroupID int) error {
	ctx := sctx.Context()
	mgr := s.PGRepo.Mgr(ctx, s.DB.GetDb(ctx))

	_, role := s.CheckStaffGroup(ctx, 0)
	if role != model.RoleAdmin {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"只有管理员可以访问此接口",
		)
	}

	bean := &model.ProjectGroup{
		ID:         projectGroupID,
		DeleteTime: int(time.Now().Unix()),
	}

	err := mgr.UpdateProjectGroup(bean)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}
