package handler

import (
	"github.com/HYY-yu/sail/internal/service/sail/api/svc"
	"github.com/HYY-yu/seckill.pkg/core"
)

type EnvTestHandler struct {
	envTestSvc *svc.TestDataSvc
}

func NewEnvTestHandler(envTestSvc *svc.TestDataSvc) *EnvTestHandler {
	return &EnvTestHandler{
		envTestSvc: envTestSvc,
	}
}

// CreateData
// @Summary  创建测试数据
// @Tags     集成测试管理
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/env_test/create    [POST]
func (h *EnvTestHandler) CreateData(c core.Context) {
	err := h.envTestSvc.CreateTestData(c.SvcContext())
	c.AbortWithError(err)
	c.Payload(nil)
}

// CleanData
// @Summary  清除测试数据
// @Tags     集成测试管理
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/env_test/clean    [POST]
func (h *EnvTestHandler) CleanData(c core.Context) {
	err := h.envTestSvc.CleanTestData(c.SvcContext())
	c.AbortWithError(err)
	c.Payload(nil)
}
