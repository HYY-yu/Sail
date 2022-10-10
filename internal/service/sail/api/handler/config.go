package handler

import (
	"net/http"

	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/HYY-yu/sail/internal/service/sail/api/svc"
	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type ConfigHandler struct {
	configSvc *svc.ConfigSvc
}

func NewConfigHandler(configSvc *svc.ConfigSvc) *ConfigHandler {
	return &ConfigHandler{
		configSvc: configSvc,
	}
}

// MetaConfig
// @Summary  获取元配置
// @Tags     配置管理
// @Param    project_id  query     int                                            false  "配置ID"
// @Param    project_group_id  query     int                                            false  "配置组ID"
// @Param    namespace_id  query     int                                            false  "配置组ID"
// @Success  200         {object}  response.JsonResponse{data=string}  "data"
// @Router   /v1/config/meta    [GET]
func (h *ConfigHandler) MetaConfig(c core.Context) {
	type Param struct {
		Temp           string `form:"temp"`
		ProjectID      int    `form:"project_id"`
		ProjectGroupID int    `form:"project_group_id"`
		NamespaceID    int    `form:"namespace_id"`
	}

	params := &Param{}

	err := c.ShouldBindForm(params)
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

	data, err := h.configSvc.GetTemplate(c.SvcContext(), params.Temp, params.ProjectID, params.ProjectGroupID, params.NamespaceID)
	c.AbortWithError(err)
	c.Payload(data)
}

// Tree
// @Summary  配置树
// @Tags     配置管理
// @Param    project_id  query     int                                            false  "配置ID"
// @Param    project_group_id  query     int                                            false  "配置组ID"
// @Success  200         {object}  response.JsonResponse{data=[]model.ProjectTree}  "data"
// @Router   /v1/config/tree    [GET]
func (h *ConfigHandler) Tree(c core.Context) {
	type Param struct {
		ProjectID      int `form:"project_id"`
		ProjectGroupID int `form:"project_group_id"`
	}

	params := &Param{}

	err := c.ShouldBindForm(params)
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

	data, err := h.configSvc.Tree(c.SvcContext(), params.ProjectID, params.ProjectGroupID)
	c.AbortWithError(err)
	c.Payload(data)
}

// Info
// @Summary  配置详情
// @Tags     配置管理
// @Param    config_id  query     int                                           false  "配置ID"
// @Success  200        {object}  response.JsonResponse{data=model.ConfigInfo}  "data"
// @Router   /v1/config/info    [GET]
func (h *ConfigHandler) Info(c core.Context) {
	type Param struct {
		ConfigID int `form:"config_id"`
	}
	params := &Param{}

	err := c.ShouldBindForm(params)
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

	data, err := h.configSvc.Info(c.SvcContext(), params.ConfigID)
	c.AbortWithError(err)
	c.Payload(data)
}

// History
// @Summary  配置历史
// @Tags     配置管理
// @Param    config_id  query     int                                                    false  "配置ID"
// @Success  200        {object}  response.JsonResponse{data=[]model.ConfigHistoryList}  "data"
// @Router   /v1/config/history    [GET]
func (h *ConfigHandler) History(c core.Context) {
	type Param struct {
		ConfigID int `form:"config_id"`
	}
	params := &Param{}

	err := c.ShouldBindForm(params)
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

	data, err := h.configSvc.History(c.SvcContext(), params.ConfigID)
	c.AbortWithError(err)
	c.Payload(data)
}

// HistoryInfo
// @Summary  配置历史详情
// @Tags     配置管理
// @Param    config_id  query     int                                                    false  "配置ID"
// @Param    reversion  query     int                                                    false  "reversion"
// @Success  200        {object}  response.JsonResponse{data=string}  "data"
// @Router   /v1/config/history_info    [GET]
func (h *ConfigHandler) HistoryInfo(c core.Context) {
	type Param struct {
		ConfigID  int `form:"config_id"`
		Reversion int `form:"reversion"`
	}
	params := &Param{}

	err := c.ShouldBindForm(params)
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

	data, err := h.configSvc.HistoryInfo(c.SvcContext(), params.ConfigID, params.Reversion)
	c.AbortWithError(err)
	c.Payload(data)
}

// Rollback
// @Summary  回滚配置
// @Tags     配置管理
// @Param    params  body      model.RollbackConfig                true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/config/rollback    [POST]
func (h *ConfigHandler) Rollback(c core.Context) {
	params := &model.RollbackConfig{}

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

	err = h.configSvc.Rollback(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

// Add
// @Summary  添加配置
// @Tags     配置管理
// @Param    params  body      model.AddConfig                     true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/config/add    [POST]
func (h *ConfigHandler) Add(c core.Context) {
	params := &model.AddConfig{}

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

	err = h.configSvc.Add(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

// Edit
// @Summary  编辑配置
// @Tags     配置管理
// @Param    params  body      model.EditConfig                    true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/config/edit    [POST]
func (h *ConfigHandler) Edit(c core.Context) {
	params := &model.EditConfig{}

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

	err = h.configSvc.Edit(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

// Copy
// @Summary  副本配置
// @Tags     配置管理
// @Param    params  body      model.ConfigCopy                    true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/config/copy    [POST]
func (h *ConfigHandler) Copy(c core.Context) {
	params := &model.ConfigCopy{}

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

	err = h.configSvc.Copy(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

// Del
// @Summary  删除配置
// @Tags     配置管理
// @Param    config_id  body      int                                 true  "ConfigId"
// @Success  200        {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/config/del    [POST]
func (h *ConfigHandler) Del(c core.Context) {
	type Param struct {
		ConfigID int `json:"config_id"`
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

	err = h.configSvc.Del(c.SvcContext(), params.ConfigID)
	c.AbortWithError(err)
	c.Payload(nil)
}
