package handler

import (
	"github.com/HYY-yu/seckill.pkg/core"
)

type ProjectGroupHandler struct {
}

func NewProjectGroupHandler() *ProjectGroupHandler {
	return &ProjectGroupHandler{}
}

// List
// @Summary  项目组列表
// @Tags    项目组管理
// @Param       page_index      query     int    false    "页号"     default(1)
// @Param       page_size      query     int    false    "页长"     default(10)
// @Param       sort      query     string    false    "排序字段"
// @Param       group_id     query     int   false   "项目ID"
// @Param       group_name  query  string    false "项目名称"
// @Success     200      {object}
// @Router      /v1/projectGroup/list    [GET]
func (h *ProjectGroupHandler) List(c core.Context) {
}
