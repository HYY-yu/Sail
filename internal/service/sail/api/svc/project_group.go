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

type ProjectGroupSvc struct {
	DB db.Repo

	PGRepo repo.ProjectGroupRepo
}

func NewProjectGroupSvc(
	db db.Repo,
	pgRepo repo.ProjectGroupRepo,
) *ProjectGroupSvc {
	svc := &ProjectGroupSvc{
		DB:     db,
		PGRepo: pgRepo,
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
	Op := make([]repo.Option, 0)
	if v, ok := pr.Filter["project_group_id"]; ok && util.IsNotZero(v) {
		Op = append(Op, mgr.WithID(gconv.Int(v)))
	}
	if v, ok := pr.Filter["project_group_name"]; ok && util.IsNotZero(v) {
		Op = append(Op, mgr.WithName(
			util.WrapSqlLike(gconv.String(v)),
			" LIKE ?",
		))
	}

	data, err := mgr.WithOptions(Op...).ListProjectGroup(limit, offset, sort)
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
			CreateBy:       e.CreateBy, // TODO GetCreateByName
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
	return nil
}
