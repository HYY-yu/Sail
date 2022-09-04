package handler

import (
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type StaffHandler struct {
}

func NewStaffHandler() *StaffHandler {
	return &StaffHandler{}
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
	_ = response.JsonResponse{}
	_ = page.Page{}
	_ = model.StaffList{}
}

// Add
// @Summary  添加员工
// @Tags     员工管理
// @Param    params  body      model.AddStaff                      true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/staff/add    [POST]
func (h *StaffHandler) Add(c core.Context) {

}

// Edit
// @Summary  编辑员工
// @Tags     员工管理
// @Param    params  body      model.EditStaff                     true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/staff/edit    [POST]
func (h *StaffHandler) Edit(c core.Context) {

}

// Grant
// @Summary  赋权员工
// @Tags     员工管理
// @Param    params  body      model.GrantStaff                    true  "data"
// @Success  200     {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/staff/grant    [POST]
func (h *StaffHandler) Grant(c core.Context) {

}

// Del
// @Summary  删除员工
// @Tags     员工管理
// @Param    staff_id  body      int                                 true  "StaffId"
// @Success  200       {object}  response.JsonResponse{data=string}  "data=ok"
// @Router   /v1/staff/del    [POST]
func (h *StaffHandler) Del(c core.Context) {

}
