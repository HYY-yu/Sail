package handler

import (
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type ProjectHandler struct {
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{}
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
	_ = response.JsonResponse{}
	_ = page.Page{}
	_ = model.ProjectGroupList{}
}

// Add
// @Summary  添加项目
// @Tags     项目管理
// @Param    params  body      model.AddProject                    true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/project/add    [POST]
func (h *ProjectHandler) Add(c core.Context) {

}

// Edit
// @Summary  编辑项目
// @Tags     项目管理
// @Param    params  body      model.EditProject                   true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/project/edit    [POST]
func (h *ProjectHandler) Edit(c core.Context) {

}

// Del
// @Summary  删除项目
// @Tags     项目管理
// @Param    project_id  body      true                                "ProjectId"
// @Success  200         {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/project/del    [POST]
func (h *ProjectHandler) Del(c core.Context) {

}
