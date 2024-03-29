package repo

import (
	"context"
	"fmt"
	"github.com/HYY-yu/sail/internal/service/sail/model"
	"gorm.io/gorm"
)

// Code generated by gormt. DO NOT EDIT.

// 非线程安全

type _StaffGroupRelMgr struct {
	*_BaseMgr
}

// StaffGroupRelMgr open func
func StaffGroupRelMgr(ctx context.Context, db *gorm.DB) *_StaffGroupRelMgr {
	if db == nil {
		panic(fmt.Errorf("StaffGroupRelMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(ctx)
	return &_StaffGroupRelMgr{_BaseMgr: &_BaseMgr{DB: db.Table("staff_group_rel").WithContext(ctx), isRelated: globalIsRelated, ctx: ctx, cancel: cancel}}
}

func (obj *_StaffGroupRelMgr) WithSelects(idName string, selects ...string) *_StaffGroupRelMgr {
	if len(idName) > 0 {
		selects = append(selects, idName)
	}
	if len(selects) > 0 {
		// 对Select进行去重
		selectMap := make(map[string]int, len(selects))
		for _, e := range selects {
			if _, ok := selectMap[e]; !ok {
				selectMap[e] = 1
			}
		}

		newSelects := make([]string, 0, len(selects))
		for k := range selectMap {
			if len(k) > 0 {
				newSelects = append(newSelects, k)
			}
		}
		obj.DB = obj.DB.Select(newSelects)
	}
	return obj
}

func (obj *_StaffGroupRelMgr) WithOptions(opts ...Option) *_StaffGroupRelMgr {
	obj.Reset()

	options := options{
		query: make(map[string]queryData, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	for k, v := range options.query {
		if v.data == nil {
			obj.DB = obj.DB.Where(k + " " + v.cond)
		} else {
			obj.DB = obj.DB.Where(k+" "+v.cond, v.data)
		}
	}
	return obj
}

// GetTableName get sql table name.获取表名字
func (obj *_StaffGroupRelMgr) GetTableName() string {
	return "staff_group_rel"
}

// Tx 开启事务会话
func (obj *_StaffGroupRelMgr) Tx(tx *gorm.DB) *_StaffGroupRelMgr {
	obj.DB = tx.Table(obj.GetTableName()).WithContext(obj.ctx)
	return obj
}

// WithPrepareStmt 开启语句 PrepareStmt 功能
// 接下来执行的SQL将会是PrepareStmt的
func (obj *_StaffGroupRelMgr) WithPrepareStmt() {
	obj.DB = obj.DB.Session(&gorm.Session{Context: obj.ctx, PrepareStmt: true})
}

// Reset 重置gorm会话
func (obj *_StaffGroupRelMgr) Reset() *_StaffGroupRelMgr {
	obj.DB = obj.DB.Session(&gorm.Session{NewDB: true, Context: obj.ctx}).Table(obj.GetTableName())
	return obj
}

// Get 获取
func (obj *_StaffGroupRelMgr) Get() (result model.StaffGroupRel, err error) {
	err = obj.DB.Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_StaffGroupRelMgr) Gets() (results []model.StaffGroupRel, err error) {
	err = obj.DB.Find(&results).Error

	return
}

// Catch 必须获取结果（单条）
func (obj *_StaffGroupRelMgr) Catch() (results model.StaffGroupRel, err error) {
	err = obj.DB.Take(&results).Error

	return
}

func (obj *_StaffGroupRelMgr) Count() (count int64, err error) {
	err = obj.DB.Count(&count).Error

	return
}

func (obj *_StaffGroupRelMgr) HasRecord() (bool, error) {
	count, err := obj.Count()
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// WithID id获取
func (obj *_StaffGroupRelMgr) WithID(id interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["id"] = queryData{
			cond: cond[0],
			data: id,
		}
	})
}

// WithProjectGroupID project_group_id获取
func (obj *_StaffGroupRelMgr) WithProjectGroupID(projectGroupID interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["project_group_id"] = queryData{
			cond: cond[0],
			data: projectGroupID,
		}
	})
}

// WithStaffID staff_id获取
func (obj *_StaffGroupRelMgr) WithStaffID(staffID interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["staff_id"] = queryData{
			cond: cond[0],
			data: staffID,
		}
	})
}

// WithRoleType role_type获取 权限角色
func (obj *_StaffGroupRelMgr) WithRoleType(roleType interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["role_type"] = queryData{
			cond: cond[0],
			data: roleType,
		}
	})
}

func (obj *_StaffGroupRelMgr) CreateStaffGroupRel(bean *model.StaffGroupRel) (err error) {
	err = obj.DB.Create(bean).Error

	return
}

func (obj *_StaffGroupRelMgr) UpdateStaffGroupRel(bean *model.StaffGroupRel) (err error) {
	err = obj.DB.Updates(bean).Error

	return
}

func (obj *_StaffGroupRelMgr) DeleteStaffGroupRel(bean *model.StaffGroupRel) (err error) {
	err = obj.DB.Delete(bean).Error

	return
}

// FetchIndexByProjectGroupID  获取多个内容
func (obj *_StaffGroupRelMgr) FetchIndexByProjectGroupID(projectGroupID int) (results []*model.StaffGroupRel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.StaffGroupRel{}).Where("`project_group_id` = ?", projectGroupID).Find(&results).Error

	return
}

// FetchIndexByStaffID  获取多个内容
func (obj *_StaffGroupRelMgr) FetchIndexByStaffID(staffID int) (results []*model.StaffGroupRel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.StaffGroupRel{}).Where("`staff_id` = ?", staffID).Find(&results).Error

	return
}
