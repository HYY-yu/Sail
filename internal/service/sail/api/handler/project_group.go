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

type ProjectGroupHandler struct {
	projectGroupSvc *svc.ProjectGroupSvc
}

func NewProjectGroupHandler(projectGroupSvc *svc.ProjectGroupSvc) *ProjectGroupHandler {
	return &ProjectGroupHandler{
		projectGroupSvc: projectGroupSvc,
	}
}

// List
// @Summary  项目组列表
// @Tags     项目组管理
// @Param    page_index          query     int                                                                 false  "页号"  default(1)
// @Param    page_size           query     int                                                                 false  "页长"  default(10)
// @Param    sort                query     string                                                              false  "排序字段"
// @Param    project_group_id    query     int                                                                 false  "项目组ID"
// @Param    project_group_name  query     string                                                              false  "项目组名称"
// @Success  200                 {object}  response.JsonResponse{data=page.Page{List=model.ProjectGroupList}}  "data"
// @Router   /v1/project_group/list    [GET]
func (h *ProjectGroupHandler) List(c core.Context) {
	err := c.RequestContext().Request.ParseForm()
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}
	pageRequest := page.NewPageFromRequest(c.RequestContext().Request.Form)

	data, err := h.projectGroupSvc.List(c.SvcContext(), pageRequest)
	c.AbortWithError(err)
	c.Payload(data)
}

// Add
// @Summary  添加项目组
// @Tags     项目组管理
// @Param    params  body      model.AddProjectGroup               true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/project_group/add    [POST]
func (h *ProjectGroupHandler) Add(c core.Context) {
	params := &model.AddProjectGroup{}

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

	err = h.projectGroupSvc.Add(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

// Edit
// @Summary  编辑项目组
// @Tags     项目组管理
// @Param    params  body      model.EditProjectGroup              true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/project_group/edit    [POST]
func (h *ProjectGroupHandler) Edit(c core.Context) {
	params := &model.EditProjectGroup{}

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

	err = h.projectGroupSvc.Edit(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

// Del
// @Summary  删除项目组
// @Tags     项目组管理
// @Param    group_id  body      int                                 true  "GroupId"
// @Success  200       {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/project_group/del    [POST]
func (h *ProjectGroupHandler) Del(c core.Context) {
	type Param struct {
		GroupID int `json:"group_id"`
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

	err = h.projectGroupSvc.Delete(c.SvcContext(), params.GroupID)
	c.AbortWithError(err)
	c.Payload(nil)
}
