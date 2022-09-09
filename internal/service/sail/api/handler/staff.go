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

type StaffHandler struct {
	staffSvc *svc.StaffSvc
}

func NewStaffHandler(staffSvc *svc.StaffSvc) *StaffHandler {
	return &StaffHandler{
		staffSvc: staffSvc,
	}
}

// List
// @Summary  员工列表
// @Tags     员工管理
// @Param    page_index  query     int                                                          false  "页号"  default(1)
// @Param    page_size   query     int                                                          false  "页长"  default(10)
// @Param    sort        query     string                                                       false  "排序字段"
// @Param    staff_id    query     int                                                          false  "员工ID"
// @Param    staff_name  query     string                                                       false  "员工名称"
// @Success  200         {object}  response.JsonResponse{data=page.Page{List=model.StaffList}}  "data"
// @Router   /v1/staff/list    [GET]
func (h *StaffHandler) List(c core.Context) {
	err := c.RequestContext().Request.ParseForm()
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}
	pageRequest := page.NewPageFromRequest(c.RequestContext().Request.Form)

	data, err := h.staffSvc.List(c.SvcContext(), pageRequest)
	c.AbortWithError(err)
	c.Payload(data)
}

// Add
// @Summary  添加员工
// @Tags     员工管理
// @Param    params  body      model.AddStaff                      true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/staff/add    [POST]
func (h *StaffHandler) Add(c core.Context) {
	params := &model.AddStaff{}

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

	err = h.staffSvc.Add(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

// Edit
// @Summary  编辑员工
// @Tags     员工管理
// @Param    params  body      model.EditStaff                     true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/staff/edit    [POST]
func (h *StaffHandler) Edit(c core.Context) {
	params := &model.EditStaff{}

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

	err = h.staffSvc.Edit(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

// Grant
// @Summary  赋权员工
// @Tags     员工管理
// @Param    params  body      model.GrantStaff                    true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/staff/grant    [POST]
func (h *StaffHandler) Grant(c core.Context) {
	params := &model.GrantStaff{}

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

	err = h.staffSvc.Grant(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

// DelGrant
// @Summary  删除授权
// @Tags     员工管理
// @Param    staff_group_rel_id  body      int                                 true  "staff_group_rel_id"
// @Success  200       {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/staff/del_grant    [POST]
func (h *StaffHandler) DelGrant(c core.Context) {
	type Param struct {
		StaffGroupRelID int `json:"staff_group_rel_id"`
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

	err = h.staffSvc.DelGrant(c.SvcContext(), params.StaffGroupRelID)
	c.AbortWithError(err)
	c.Payload(nil)
}

// Del
// @Summary  删除员工
// @Tags     员工管理
// @Param    staff_id  body      int                                 true  "StaffId"
// @Success  200       {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/staff/del    [POST]
func (h *StaffHandler) Del(c core.Context) {
	type Param struct {
		StaffID int `json:"staff_id"`
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

	err = h.staffSvc.Delete(c.SvcContext(), params.StaffID)
	c.AbortWithError(err)
	c.Payload(nil)
}
