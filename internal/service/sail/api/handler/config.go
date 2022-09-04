package handler

import (
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type ConfigHandler struct {
}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

// Tree
// @Summary  配置树
// @Tags     配置管理
// @Param    project_id  query     int                                            false  "配置ID"
// @Success  200         {object}  response.JsonResponse{data=model.ProjectTree}  "data"
// @Router   /v1/config/tree    [GET]
func (h *ConfigHandler) Tree(c core.Context) {
	_ = response.JsonResponse{}
	_ = page.Page{}
	_ = model.ProjectTree{}
}

// Info
// @Summary  配置详情
// @Tags     配置管理
// @Param    config_id  query     int                                           false  "配置ID"
// @Success  200        {object}  response.JsonResponse{data=model.ConfigInfo}  "data"
// @Router   /v1/config/info    [GET]
func (h *ConfigHandler) Info(c core.Context) {
	_ = response.JsonResponse{}
	_ = page.Page{}
	_ = model.ProjectTree{}
}

// History
// @Summary  配置历史
// @Tags     配置管理
// @Param    config_id  query     int                                                    false  "配置ID"
// @Success  200        {object}  response.JsonResponse{data=[]model.ConfigHistoryList}  "data"
// @Router   /v1/config/history    [GET]
func (h *ConfigHandler) History(c core.Context) {
	_ = response.JsonResponse{}
	_ = page.Page{}
	_ = model.ProjectTree{}
}

// Rollback
// @Summary  回滚配置
// @Tags     配置管理
// @Param    params  body      model.RollbackConfig                true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/config/rollback    [POST]
func (h *ConfigHandler) Rollback(c core.Context) {

}

// Add
// @Summary  添加配置
// @Tags     配置管理
// @Param    params  body      model.AddConfig                     true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/config/add    [POST]
func (h *ConfigHandler) Add(c core.Context) {

}

// Edit
// @Summary  编辑配置
// @Tags     配置管理
// @Param    params  body      model.EditConfig                    true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/config/edit    [POST]
func (h *ConfigHandler) Edit(c core.Context) {

}

// Copy
// @Summary  副本配置
// @Tags     配置管理
// @Param    params  body      model.ConfigCopy                    true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/config/copy    [POST]
func (h *ConfigHandler) Copy(c core.Context) {

}

// Del
// @Summary  删除配置
// @Tags     配置管理
// @Param    config_id  body      int                                 true  "ConfigId"
// @Success  200        {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/config/del    [POST]
func (h *ConfigHandler) Del(c core.Context) {

}
