package handler

import (
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type PublishHandler struct {
}

func NewPublishHandler() *PublishHandler {
	return &PublishHandler{}
}

// List
// @Summary  发布列表
// @Tags     发布管理
// @Param    page_index  query     int                                                                  false  "页号"  default(1)
// @Param    page_size   query     int                                                                  false  "页长"  default(10)
// @Param    sort        query     string                                                               false  "排序字段"
// @Param    project_id  query     int                                                                  true   "项目ID"
// @Success  200         {object}  response.JsonResponse{data=page.Page{List=model.PublishConfigList}}  "data"
// @Router   /v1/publish/list    [GET]
func (h *PublishHandler) List(c core.Context) {
	_ = response.JsonResponse{}
	_ = page.Page{}
	_ = model.ProjectGroupList{}
}

// Add
// @Summary  添加发布
// @Tags     发布管理
// @Param    params  body      model.AddPublish                    true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/publish/add    [POST]
func (h *PublishHandler) Add(c core.Context) {

}

// Rollback
// @Summary  添加发布
// @Tags     发布管理
// @Param    params  body      model.RollbackPublish               true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/publish/rollback    [POST]
func (h *PublishHandler) Rollback(c core.Context) {

}
