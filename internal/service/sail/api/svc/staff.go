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
	"golang.org/x/crypto/bcrypt"

	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type StaffSvc struct {
	BaseSvc
	DB db.Repo

	StaffRepo      repo.StaffRepo
	StaffGroupRepo repo.StaffGroupRelRepo
	PGRepo         repo.ProjectGroupRepo
}

func NewStaffSvc(
	db db.Repo,
	r repo.StaffRepo,
	sgr repo.StaffGroupRelRepo,
	pgr repo.ProjectGroupRepo,
) *StaffSvc {
	svc := &StaffSvc{
		DB:             db,
		StaffRepo:      r,
		StaffGroupRepo: sgr,
		PGRepo:         pgr,
	}
	return svc
}

func (s *StaffSvc) List(sctx core.SvcContext, pr *page.PageRequest) (*page.Page, error) {
	ctx := sctx.Context()
	mgr := s.StaffRepo.Mgr(ctx, s.DB.GetDb())
	sgMgr := s.StaffGroupRepo.Mgr(ctx, s.DB.GetDb())
	pgMgr := s.PGRepo.Mgr(ctx, s.DB.GetDb())

	_, role := s.CheckStaffGroup(ctx, 0)
	if role != model.RoleAdmin {
		return nil, response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"只有管理员可以访问此接口",
		)
	}

	limit, offset := pr.GetLimitAndOffset()
	pr.AddAllowSortField(model.StaffColumns.CreateTime)
	sort, _ := pr.Sort()

	op := make([]repo.Option, 0)
	if v, ok := pr.Filter["staff_id"]; ok && util.IsNotZero(v) {
		op = append(op, mgr.WithID(gconv.Int(v)))
	}
	if v, ok := pr.Filter["staff_name"]; ok && util.IsNotZero(v) {
		op = append(op, mgr.WithName(
			util.WrapSqlLike(gconv.String(v)),
			" LIKE ?",
		))
	}
	mgr.UpdateDB(mgr.WithPrepareStmt())
	sgMgr.UpdateDB(sgMgr.WithPrepareStmt())
	pgMgr.UpdateDB(pgMgr.WithPrepareStmt())

	data, err := mgr.WithOptions(op...).ListStaff(limit, offset, sort)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	var count int64
	mgr.Count(&count)

	var result = make([]model.StaffList, len(data))

	for i, e := range data {
		staffGroupRel, err := sgMgr.WithOptions(sgMgr.WithStaffID(e.ID)).Gets()
		if err != nil {
			return nil, response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
		roles := make([]model.StaffRole, len(staffGroupRel))

		for j, sg := range staffGroupRel {
			projectGroup, _ := pgMgr.
				WithSelects(
					model.ProjectGroupColumns.ID,
					model.ProjectGroupColumns.Name,
				).
				WithOptions(pgMgr.WithID(sg.ProjectGroupID)).
				Get()
			roles[j] = model.StaffRole{
				StaffGroupRelID:  sg.ID,
				ProjectGroupID:   projectGroup.ID,
				ProjectGroupName: projectGroup.Name,
				Role:             model.Role(sg.RoleType),
				RoleInfo:         model.Role(sg.RoleType).String(),
			}
		}

		r := model.StaffList{
			StaffID:    e.ID,
			Name:       e.Name,
			CreateTime: e.CreateTime.Unix(),
			Roles:      roles,
		}
		result[i] = r
	}
	return page.NewPage(
		count,
		result,
	), nil
}

func (s *StaffSvc) Add(sctx core.SvcContext, param *model.AddStaff) error {
	ctx := sctx.Context()
	mgr := s.StaffRepo.Mgr(ctx, s.DB.GetDb())

	_, role := s.CheckStaffGroup(ctx, 0)
	if role != model.RoleAdmin {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"只有管理员可以访问此接口",
		)
	}

	// 默认密码： 123456
	passwordByte, _ := bcrypt.GenerateFromPassword([]byte("123456"), 0)

	bean := &model.Staff{
		Name:       param.Name,
		Password:   string(passwordByte),
		CreateTime: time.Now(),
		CreateBy:   int(sctx.UserId()),
	}
	err := mgr.CreateStaff(bean)
	if err != nil {
		if mysqlerr_helper.IsMysqlDupEntryError(err) {
			return response.NewErrorWithStatusOk(
				response.ParamBindError,
				"已经存在相同的Name，请保证Key唯一",
			)
		}
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (s *StaffSvc) Edit(sctx core.SvcContext, param *model.EditStaff) error {
	ctx := sctx.Context()
	mgr := s.StaffRepo.Mgr(ctx, s.DB.GetDb())
	_, role := s.CheckStaffGroup(ctx, 0)
	if role != model.RoleAdmin {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"只有管理员可以访问此接口",
		)
	}

	bean := &model.Staff{
		ID: param.StaffID,
	}

	updateColumns := make([]string, 0)

	if param.Name != nil && !g.IsEmpty(*param.Name) {
		bean.Name = *param.Name
		updateColumns = append(updateColumns, model.ProjectGroupColumns.Name)
	}

	err := mgr.WithSelects(model.ProjectGroupColumns.ID, updateColumns...).UpdateStaff(bean)
	if err != nil {
		if mysqlerr_helper.IsMysqlDupEntryError(err) {
			return response.NewErrorWithStatusOk(
				response.ParamBindError,
				"已经存在相同的Name，请保证Key唯一",
			)
		}
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (s *StaffSvc) Delete(sctx core.SvcContext, staffID int) error {
	ctx := sctx.Context()
	mgr := s.StaffRepo.Mgr(ctx, s.DB.GetDb())
	sgMgr := s.StaffGroupRepo.Mgr(ctx, s.DB.GetDb())
	_, role := s.CheckStaffGroup(ctx, 0)
	if role != model.RoleAdmin {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"只有管理员可以访问此接口",
		)
	}

	bean := &model.Staff{
		ID: staffID,
	}
	tx := s.DB.GetDb().Begin()
	defer tx.Rollback()

	mgr.Tx(tx)
	sgMgr.Tx(tx)

	// 删除此员工的对应的 staffGroup
	err := sgMgr.WithOptions(sgMgr.WithStaffID(staffID)).Delete(&model.StaffGroupRel{}).Error
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	err = mgr.DeleteStaff(bean)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	tx.Commit()
	return nil
}

func (s *StaffSvc) Grant(sctx core.SvcContext, param *model.GrantStaff) error {
	ctx := sctx.Context()
	mgr := s.StaffRepo.Mgr(ctx, s.DB.GetDb())
	sgMgr := s.StaffGroupRepo.Mgr(ctx, s.DB.GetDb())
	pgMgr := s.PGRepo.Mgr(ctx, s.DB.GetDb())
	_, role := s.CheckStaffGroup(ctx, 0)
	if role != model.RoleAdmin {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"只有管理员可以访问此接口",
		)
	}

	hasRecord, err := mgr.WithOptions(mgr.WithID(param.StaffID)).HasRecord()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if !hasRecord {
		return response.NewErrorWithStatusOk(
			response.ServerError,
			"未找到此员工",
		).WithErr(err)
	}

	if param.Role != model.RoleAdmin {
		hasRecord, err = pgMgr.WithOptions(pgMgr.WithID(param.ProjectGroupID)).HasRecord()
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
		if !hasRecord {
			return response.NewErrorWithStatusOk(
				response.ServerError,
				"未找到此项目组",
			).WithErr(err)
		}
	}

	// 如果用户是管理员，他不能被赋予其它任何权限
	// 如果用户不是管理员，他不能被赋予管理员，除非删除他的其它权限
	sgList, _ := sgMgr.WithOptions(sgMgr.WithStaffID(param.StaffID)).Gets()
	if param.Role == model.RoleAdmin {
		if len(sgList) > 0 {
			return response.NewErrorWithStatusOk(
				response.ServerError,
				"此用户不能被赋权，请删除他的其它权限",
			).WithErr(err)
		}
	}
	for _, e := range sgList {
		if model.Role(e.RoleType) == model.RoleAdmin {
			return response.NewErrorWithStatusOk(
				response.ServerError,
				"管理员不能被赋予其它权限，它已经是最高权限",
			).WithErr(err)
		}
	}

	bean := &model.StaffGroupRel{
		ProjectGroupID: param.ProjectGroupID,
		StaffID:        param.StaffID,
		RoleType:       int(param.Role),
	}

	err = sgMgr.CreateStaffGroupRel(bean)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (s *StaffSvc) DelGrant(sctx core.SvcContext, staffGroupRelID int) error {
	ctx := sctx.Context()
	sgMgr := s.StaffGroupRepo.Mgr(ctx, s.DB.GetDb())
	_, role := s.CheckStaffGroup(ctx, 0)
	if role != model.RoleAdmin {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"只有管理员可以访问此接口",
		)
	}

	bean := &model.StaffGroupRel{
		ID: staffGroupRelID,
	}

	err := sgMgr.DeleteStaffGroupRel(bean)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (s *StaffSvc) StaffGroup(sctx core.SvcContext) ([]model.StaffGroup, error) {
	ctx := sctx.Context()
	userId := int(sctx.UserId())
	sgMgr := s.StaffGroupRepo.Mgr(ctx, s.DB.GetDb())

	rel, err := sgMgr.WithOptions(sgMgr.WithStaffID(userId)).Gets()
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	result := make([]model.StaffGroup, len(rel))
	for i, e := range rel {
		b := model.StaffGroup{
			ProjectGroupID: e.ProjectGroupID,
			Role:           model.Role(e.RoleType),
		}
		result[i] = b
	}
	return result, nil
}
