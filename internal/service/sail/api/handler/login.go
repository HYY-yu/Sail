package handler

import (
	"net/http"

	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/HYY-yu/sail/internal/service/sail/api/svc"
	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type LoginHandler struct {
	loginSvc *svc.LoginSvc
}

func NewLoginHandler(loginSvc *svc.LoginSvc) *LoginHandler {
	return &LoginHandler{
		loginSvc: loginSvc,
	}
}

func (h *LoginHandler) LoginHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// Login
// @Summary  登陆
// @Tags     WEB
// @Param    params  body      model.LoginParams                      true  "data"
// @Success  200     {object}  response.JsonResponse{data=model.LoginResponse}  "data=ok"
// @Router   /v1/login    [POST]
func (h *LoginHandler) Login(c core.Context) {
	params := &model.LoginParams{}

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

	data, err := h.loginSvc.Login(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(data)
}

// RefreshToken
// @Summary  刷新TOKEN
// @Tags     WEB
// @Param    old_refresh_token  body      string                                 true  "老的RefreshToken"
// @Success  200       {object}  response.JsonResponse{data=model.LoginResponse}  "data=ok"
// @Router   /v1/login/refresh    [POST]
func (h *LoginHandler) RefreshToken(c core.Context) {
	type Param struct {
		OldRefreshToken string `json:"old_refresh_token"`
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

	data, err := h.loginSvc.RefreshToken(c.SvcContext(), params.OldRefreshToken)
	c.AbortWithError(err)
	c.Payload(data)
}

// LoginOut
// @Summary  登出
// @Tags     WEB
// @Success  200       {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/login/login_out    [GET]
func (h *LoginHandler) LoginOut(c core.Context) {
	err := h.loginSvc.LoginOut(c.SvcContext())
	c.AbortWithError(err)
	c.Payload(nil)
}

// ChangePassword
// @Summary  更改密码
// @Tags     WEB
// @Param    new_pass  body      string                                 true  "新密码"
// @Success  200       {object}  response.JsonResponse{data=model.LoginResponse}  "data=ok"
// @Router   /v1/login/new_pass    [POST]
func (h *LoginHandler) ChangePassword(c core.Context) {
	type Param struct {
		NewPass string `json:"new_pass"`
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

	err = h.loginSvc.ChangePassword(c.SvcContext(), params.NewPass)
	c.AbortWithError(err)
	c.Payload(nil)
}
