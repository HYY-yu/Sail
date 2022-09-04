package handler

import (
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type NamespaceHandler struct {
}

func NewNamespaceHandler() *NamespaceHandler {
	return &NamespaceHandler{}
}

// List
// @Summary  命名空间列表
// @Tags     命名空间管理
// @Param    page_index      query     int                                                              false  "页号"  default(1)
// @Param    page_size       query     int                                                              false  "页长"  default(10)
// @Param    sort            query     string                                                           false  "排序字段"
// @Param    namespace_id    query     int                                                              false  "命名空间ID"
// @Param    namespace_name  query     string                                                           false  "命名空间名称"
// @Success  200             {object}  response.JsonResponse{data=page.Page{List=model.NamespaceList}}  "data"
// @Router   /v1/namespace/list    [GET]
func (h *NamespaceHandler) List(c core.Context) {
	_ = response.JsonResponse{}
	_ = page.Page{}
	_ = model.NamespaceList{}
}

// Add
// @Summary  添加命名空间
// @Tags     命名空间管理
// @Param    params  body      model.AddNamespace                  true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/namespace/add    [POST]
func (h *NamespaceHandler) Add(c core.Context) {

}

// Edit
// @Summary  编辑命名空间
// @Tags     命名空间管理
// @Param    params  body      model.EditNamespace                 true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/namespace/edit    [POST]
func (h *NamespaceHandler) Edit(c core.Context) {

}

// Del
// @Summary  删除命名空间
// @Tags     命名空间管理
// @Param    namespace_id  body      true                                "NamespaceId"
// @Success  200           {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/namespace/del    [POST]
func (h *NamespaceHandler) Del(c core.Context) {

}
