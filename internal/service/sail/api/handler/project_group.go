package handler

import (
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type ProjectGroupHandler struct {
}

func NewProjectGroupHandler() *ProjectGroupHandler {
	return &ProjectGroupHandler{}
}

// List
// @Summary  项目组列表
// @Tags     项目组管理
// @Param    page_index  query     int                                                                 false  "页号"  default(1)
// @Param    page_size   query     int                                                                 false  "页长"  default(10)
// @Param    sort        query     string                                                              false  "排序字段"
// @Param    project_group_id    query     int                                                                 false  "项目组ID"
// @Param    project_group_name  query     string                                                              false  "项目组名称"
// @Success  200         {object}  response.JsonResponse{data=page.Page{List=model.ProjectGroupList}}  "data"
// @Router   /v1/project_group/list    [GET]
func (h *ProjectGroupHandler) List(c core.Context) {
	_ = response.JsonResponse{}
	_ = page.Page{}
	_ = model.ProjectGroupList{}
}

// Add
// @Summary  添加项目组
// @Tags     项目组管理
// @Param    params  body      model.AddProjectGroup               true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/project_group/add    [POST]
func (h *ProjectGroupHandler) Add(c core.Context) {

}

// Edit
// @Summary  编辑项目组
// @Tags     项目组管理
// @Param    params  body      model.EditProjectGroup              true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/project_group/edit    [POST]
func (h *ProjectGroupHandler) Edit(c core.Context) {

}

// Del
// @Summary  删除项目组
// @Tags     项目组管理
// @Param    group_id  body      true                                "GroupId"
// @Success  200       {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/project_group/del    [POST]
func (h *ProjectGroupHandler) Del(c core.Context) {

}
