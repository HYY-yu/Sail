package svc

import (
	"encoding/json"
	"net/http"
	"strconv"
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

type NamespaceSvc struct {
	BaseSvc
	DB db.Repo

	NamespaceRepo    repo.NamespaceRepo
	StaffRepo        repo.StaffRepo
	ProjectGroupRepo repo.ProjectGroupRepo
}

func NewNamespaceSvc(
	db db.Repo,
	namespaceRepo repo.NamespaceRepo,
	pgRepo repo.ProjectGroupRepo,
	staffRepo repo.StaffRepo,
) *NamespaceSvc {
	svc := &NamespaceSvc{
		DB:               db,
		NamespaceRepo:    namespaceRepo,
		ProjectGroupRepo: pgRepo,
		StaffRepo:        staffRepo,
	}
	return svc
}

func (s *NamespaceSvc) List(sctx core.SvcContext, pr *page.PageRequest) (*page.Page, error) {
	ctx := sctx.Context()
	mgr := s.NamespaceRepo.Mgr(ctx, s.DB.GetDb())
	pgMgr := s.ProjectGroupRepo.Mgr(ctx, s.DB.GetDb())

	projectGroupIdInter, ok := pr.Filter["project_group_id"]
	if !ok {
		return nil, response.NewErrorWithStatusOk(
			response.ParamBindError,
			"必须提供 project_group_id",
		)
	}
	projectGroupId := gconv.Int(projectGroupIdInter)

	_, role := s.CheckStaffGroup(ctx, projectGroupId)
	if role > model.RoleOwner {
		return page.NewPage(
			0,
			[]model.ProjectList{},
		), nil
	}

	limit, offset := pr.GetLimitAndOffset()
	pr.AddAllowSortField(model.NamespaceColumns.CreateTime)
	sort, _ := pr.Sort()

	op := make([]repo.Option, 0)

	if v, ok := pr.Filter["namespace_name"]; ok && util.IsNotZero(v) {
		op = append(op, mgr.WithName(
			util.WrapSqlLike(gconv.String(v)),
			" LIKE ?",
		))
	}
	op = append(op, mgr.WithDeleteTime(0))
	op = append(op, mgr.WithProjectGroupID(projectGroupId))

	data, err := mgr.WithOptions(op...).ListNamespace(limit, offset, sort)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	count, _ := mgr.Count()
	pg, _ := pgMgr.WithOptions(pgMgr.WithID(projectGroupId)).WithSelects(
		model.ProjectGroupColumns.ID,
		model.ProjectGroupColumns.Name,
	).Get()

	result := make([]model.NamespaceList, len(data))
	for i, e := range data {
		r := model.NamespaceList{
			NamespaceID:      e.ID,
			ProjectGroupID:   e.ProjectGroupID,
			ProjectGroupName: pg.Name,
			Name:             e.Name,
			SecretKey:        e.SecretKey,
			RealTime:         e.RealTime,
			CreateTime:       e.CreateTime.Unix(),
			CreateBy:         e.CreateBy,
			CreateByName:     s.GetCreateByName(ctx, s.DB, s.StaffRepo, e.CreateBy),
		}
		result[i] = r
	}
	return page.NewPage(
		count,
		result,
	), nil
}

func (s *NamespaceSvc) Add(sctx core.SvcContext, param *model.AddNamespace) error {
	ctx := sctx.Context()
	userId := sctx.UserId()
	mgr := s.NamespaceRepo.Mgr(ctx, s.DB.GetDb())

	_, role := s.CheckStaffGroup(ctx, param.ProjectGroupID)
	if role > model.RoleOwner {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}

	bean := &model.Namespace{
		ProjectGroupID: param.ProjectGroupID,
		Name:           param.Name,
		RealTime:       param.RealTime,
		CreateTime:     time.Now(),
		CreateBy:       int(userId),
	}
	if param.Secret {
		// 生成 secret_key
		jsonBean, _ := json.Marshal(bean)
		bean.SecretKey = encrypt.SHA1WithEncoding(string(jsonBean), encrypt.NewBase32Human())
	}

	err := mgr.CreateNamespace(bean)
	if err != nil {
		if mysqlerr_helper.IsMysqlDupEntryError(err) {
			return response.NewErrorWithStatusOk(
				response.ParamBindError,
				"已经存在相同的Name，请保证唯一",
			)
		}
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (s *NamespaceSvc) Edit(sctx core.SvcContext, param *model.EditNamespace) error {
	ctx := sctx.Context()
	mgr := s.NamespaceRepo.Mgr(ctx, s.DB.GetDb())

	namespace, err := mgr.WithOptions(mgr.WithID(param.NamespaceId)).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	_, role := s.CheckStaffGroup(ctx, namespace.ProjectGroupID)
	if role > model.RoleOwner {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}

	bean := &model.Namespace{
		ID: namespace.ID,
	}
	updateColumns := make([]string, 0)

	if param.Name != nil && !g.IsEmpty(*param.Name) && *param.Name != namespace.Name {
		bean.Name = *param.Name
		updateColumns = append(updateColumns, model.NamespaceColumns.Name)
	}
	if param.RealTime != nil && *param.RealTime != namespace.RealTime {
		bean.RealTime = *param.RealTime
		updateColumns = append(updateColumns, model.NamespaceColumns.RealTime)
	}

	err = mgr.WithSelects(model.NamespaceColumns.ID, updateColumns...).UpdateNamespace(bean)
	if err != nil {
		if mysqlerr_helper.IsMysqlDupEntryError(err) {
			return response.NewErrorWithStatusOk(
				response.ParamBindError,
				"已经存在相同的Name，请保证唯一",
			)
		}
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (s *NamespaceSvc) Delete(sctx core.SvcContext, namespaceID int) error {
	ctx := sctx.Context()
	mgr := s.NamespaceRepo.Mgr(ctx, s.DB.GetDb())

	namespace, err := mgr.WithOptions(mgr.WithID(namespaceID)).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	_, role := s.CheckStaffGroup(ctx, namespace.ProjectGroupID)
	if role > model.RoleOwner {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}

	deletetime := int(time.Now().Unix())
	bean := &model.Namespace{
		ID:         namespace.ID,
		Name:       namespace.Name + strconv.Itoa(deletetime),
		DeleteTime: deletetime,
	}
	err = mgr.UpdateNamespace(bean)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}
