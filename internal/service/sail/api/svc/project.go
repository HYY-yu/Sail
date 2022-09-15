package svc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/encrypt"
	"github.com/HYY-yu/seckill.pkg/pkg/mysqlerr_helper"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/HYY-yu/seckill.pkg/pkg/util"
	"github.com/gogf/gf/v2/frame/g"
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

func (s *ProjectSvc) Add(sctx core.SvcContext, param *model.AddProject) error {
	ctx := sctx.Context()
	userId := sctx.UserId()
	tx := s.DB.GetDb(ctx).Begin()
	defer tx.Rollback()

	mgr := s.ProjectRepo.Mgr(ctx, s.DB.GetDb(ctx))
	mgr.Tx(tx)

	_, role := s.CheckStaffGroup(ctx, param.ProjectGroupID)
	if role > model.RoleOwner {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}

	bean := &model.Project{
		ProjectGroupID: param.ProjectGroupID,
		Name:           param.Name,
		CreateTime:     time.Now(),
		CreateBy:       int(userId),
	}
	err := mgr.CreateProject(bean)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	// key
	k := fmt.Sprintf("%d-%d", bean.ID, bean.ProjectGroupID)
	bean.Key = encrypt.MD5(k)

	err = mgr.WithOptions(mgr.WithID(bean.ID)).Update(model.ProjectColumns.Key, bean.Key).Error
	if err != nil {
		if mysqlerr_helper.IsMysqlDupEntryError(err) {
			return response.NewErrorWithStatusOk(
				response.ParamBindError,
				"已经存在相同的ProjectKey，请保证Key唯一",
			)
		}
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	tx.Commit()
	return nil
}

func (s *ProjectSvc) Edit(sctx core.SvcContext, param *model.EditProject) error {
	ctx := sctx.Context()
	mgr := s.ProjectRepo.Mgr(ctx, s.DB.GetDb(ctx))

	project, err := mgr.WithOptions(mgr.WithID(param.ProjectId)).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	_, role := s.CheckStaffGroup(ctx, project.ProjectGroupID)
	if role > model.RoleOwner {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}

	bean := &model.Project{
		ID: project.ID,
	}
	updateColumns := make([]string, 0)

	if param.Name != nil && !g.IsEmpty(*param.Name) && *param.Name != project.Name {
		bean.Name = *param.Name
		updateColumns = append(updateColumns, model.ProjectColumns.Name)
	}

	err = mgr.WithSelects(model.ProjectGroupColumns.ID, updateColumns...).UpdateProject(bean)
	if err != nil {
		if mysqlerr_helper.IsMysqlDupEntryError(err) {
			return response.NewErrorWithStatusOk(
				response.ParamBindError,
				"已经存在相同的ProjectKey，请保证Key唯一",
			)
		}
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (s *ProjectSvc) Delete(sctx core.SvcContext, projectID int) error {
	ctx := sctx.Context()
	mgr := s.ProjectRepo.Mgr(ctx, s.DB.GetDb(ctx))
	project, err := mgr.WithOptions(mgr.WithID(projectID)).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	_, role := s.CheckStaffGroup(ctx, project.ProjectGroupID)
	if role > model.RoleOwner {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}

	bean := &model.Project{
		ID:         project.ID,
		DeleteTime: int(time.Now().Unix()),
	}

	err = mgr.UpdateProject(bean)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}
