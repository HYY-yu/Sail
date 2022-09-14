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

type NamespaceHandler struct {
	namespaceSvc *svc.NamespaceSvc
}

func NewNamespaceHandler(namespaceSvc *svc.NamespaceSvc) *NamespaceHandler {
	return &NamespaceHandler{
		namespaceSvc: namespaceSvc,
	}
}

// List
// @Summary  命名空间列表
// @Tags     命名空间管理
// @Param    page_index      query     int                                                              false  "页号"  default(1)
// @Param    page_size       query     int                                                              false  "页长"  default(10)
// @Param    sort            query     string                                                           false  "排序字段"
// @Param    project_group_id    query     int                                                              true  "PGID"
// @Param    namespace_id    query     int                                                              false  "命名空间ID"
// @Param    namespace_name  query     string                                                           false  "命名空间名称"
// @Success  200             {object}  response.JsonResponse{data=page.Page{List=model.NamespaceList}}  "data"
// @Router   /v1/namespace/list    [GET]
func (h *NamespaceHandler) List(c core.Context) {
	err := c.RequestContext().Request.ParseForm()
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}
	pageRequest := page.NewPageFromRequest(c.RequestContext().Request.Form)

	data, err := h.namespaceSvc.List(c.SvcContext(), pageRequest)
	c.AbortWithError(err)
	c.Payload(data)
}

// Add
// @Summary  添加命名空间
// @Tags     命名空间管理
// @Param    params  body      model.AddNamespace                  true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/namespace/add    [POST]
func (h *NamespaceHandler) Add(c core.Context) {
	params := &model.AddNamespace{}

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

	err = h.namespaceSvc.Add(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

// Edit
// @Summary  编辑命名空间
// @Tags     命名空间管理
// @Param    params  body      model.EditNamespace                 true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/namespace/edit    [POST]
func (h *NamespaceHandler) Edit(c core.Context) {
	params := &model.EditNamespace{}

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

	err = h.namespaceSvc.Edit(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

// Del
// @Summary  删除命名空间
// @Tags     命名空间管理
// @Param    namespace_id  body      int                                 true  "NamespaceId"
// @Success  200           {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/namespace/del    [POST]
func (h *NamespaceHandler) Del(c core.Context) {
	type Param struct {
		NamespaceID int `json:"namespace_id"`
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

	err = h.namespaceSvc.Delete(c.SvcContext(), params.NamespaceID)
	c.AbortWithError(err)
	c.Payload(nil)
}
