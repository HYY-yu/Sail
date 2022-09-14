package handler

import (
	"net/http"

	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/HYY-yu/sail/internal/service/sail/api/svc"
	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type ProjectHandler struct {
	projectSvc *svc.ProjectSvc
}

func NewProjectHandler(projectSvc *svc.ProjectSvc) *ProjectHandler {
	return &ProjectHandler{
		projectSvc: projectSvc,
	}
}

// List
// @Summary  项目列表
// @Tags     项目管理
// @Param    page_index    query     int                                                            false  "页号"  default(1)
// @Param    page_size     query     int                                                            false  "页长"  default(10)
// @Param    sort          query     string                                                         false  "排序字段"
// @Param    project_id    query     int                                                            false  "项目ID"
// @Param    project_name  query     string                                                         false  "项目名称"
// @Success  200           {object}  response.JsonResponse{data=page.Page{List=model.ProjectList}}  "data"
// @Router   /v1/project/list    [GET]
func (h *ProjectHandler) List(c core.Context) {
	err := c.RequestContext().Request.ParseForm()
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}
	pageRequest := page.NewPageFromRequest(c.RequestContext().Request.Form)

	data, err := h.projectSvc.List(c.SvcContext(), pageRequest)
	c.AbortWithError(err)
	c.Payload(data)
}

// Add
// @Summary  添加项目
// @Tags     项目管理
// @Param    params  body      model.AddProject                    true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/project/add    [POST]
func (h *ProjectHandler) Add(c core.Context) {
	params := &model.AddProject{}

	err := c.ShouldBindJSON(params)
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}

	validErr := g.Validator().Data(params).Run(c.SvcContext().Context())
	if validErr != nil {
		c.AbortWithError(response.NewError(
			http.StatusBadRequest,
			response.ParamBindError,
			validErr.Error(),
		))
		return
	}

	err = h.projectSvc.Add(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

// Edit
// @Summary  编辑项目
// @Tags     项目管理
// @Param    params  body      model.EditProject                   true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/project/edit    [POST]
func (h *ProjectHandler) Edit(c core.Context) {
	params := &model.EditProject{}

	err := c.ShouldBindJSON(params)
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}

	validErr := g.Validator().Data(params).Run(c.SvcContext().Context())
	if validErr != nil {
		c.AbortWithError(response.NewError(
			http.StatusBadRequest,
			response.ParamBindError,
			validErr.Error(),
		))
		return
	}

	err = h.projectSvc.Edit(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

// Del
// @Summary  删除项目
// @Tags     项目管理
// @Param    project_id  body      int                                 true  "ProjectId"
// @Success  200         {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/project/del    [POST]
func (h *ProjectHandler) Del(c core.Context) {
	type Param struct {
		ProjectID int `json:"project_id"`
	}
	params := &Param{}

	err := c.ShouldBindJSON(params)
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}

	validErr := g.Validator().Data(params).Run(c.SvcContext().Context())
	if validErr != nil {
		c.AbortWithError(response.NewError(
			http.StatusBadRequest,
			response.ParamBindError,
			validErr.Error(),
		))
		return
	}

	err = h.projectSvc.Delete(c.SvcContext(), params.ProjectID)
	c.AbortWithError(err)
	c.Payload(nil)
}
