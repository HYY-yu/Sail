package repo

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

// Code generated by gormt. DO NOT EDIT.

type _StaffGroupRelMgr struct {
	*_BaseMgr
}

// StaffGroupRelMgr open func
func StaffGroupRelMgr(db *gorm.DB) *_StaffGroupRelMgr {
	if db == nil {
		panic(fmt.Errorf("StaffGroupRelMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &_StaffGroupRelMgr{_BaseMgr: &_BaseMgr{DB: db.Table("staff_group_rel"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

// WithContext set context to db
func (obj *_StaffGroupRelMgr) WithContext(c context.Context) *_StaffGroupRelMgr {
	if c != nil {
		obj.ctx = c
	}
	return obj
}

func (obj *_StaffGroupRelMgr) WithSelects(idName string, selects ...string) *_StaffGroupRelMgr {
	if len(selects) > 0 {
		if len(idName) > 0 {
			selects = append(selects, idName)
		}
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

func (obj *_StaffGroupRelMgr) WithOmit(omit ...string) *_StaffGroupRelMgr {
	if len(omit) > 0 {
		obj.DB = obj.DB.Omit(omit...)
	}
	return obj
}

func (obj *_StaffGroupRelMgr) WithOptions(opts ...Option) *_StaffGroupRelMgr {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	obj.DB = obj.DB.Where(options.query)
	return obj
}

// GetTableName get sql table name.获取数据库名字
func (obj *_StaffGroupRelMgr) GetTableName() string {
	return "staff_group_rel"
}

// Reset 重置gorm会话
func (obj *_StaffGroupRelMgr) Reset() *_StaffGroupRelMgr {
	obj.new()
	return obj
}

// Get 获取
func (obj *_StaffGroupRelMgr) Get() (result StaffGroupRel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(StaffGroupRel{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_StaffGroupRelMgr) Gets() (results []*StaffGroupRel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(StaffGroupRel{}).Find(&results).Error

	return
}

func (obj *_StaffGroupRelMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(StaffGroupRel{}).Count(count)
}

// WithID id获取
func (obj *_StaffGroupRelMgr) WithID(id int) Option {
	return optionFunc(func(o *options) { o.query["id"] = id })
}

// WithStaffGroupID staff_group_id获取
func (obj *_StaffGroupRelMgr) WithStaffGroupID(staffGroupID int) Option {
	return optionFunc(func(o *options) { o.query["staff_group_id"] = staffGroupID })
}

// WithStaffID staff_id获取
func (obj *_StaffGroupRelMgr) WithStaffID(staffID int) Option {
	return optionFunc(func(o *options) { o.query["staff_id"] = staffID })
}

// GetFromID 通过id获取内容
func (obj *_StaffGroupRelMgr) GetFromID(id int) (result StaffGroupRel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(StaffGroupRel{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找
func (obj *_StaffGroupRelMgr) GetBatchFromID(ids []int) (results []*StaffGroupRel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(StaffGroupRel{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromStaffGroupID 通过staff_group_id获取内容
func (obj *_StaffGroupRelMgr) GetFromStaffGroupID(staffGroupID int) (results []*StaffGroupRel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(StaffGroupRel{}).Where("`staff_group_id` = ?", staffGroupID).Find(&results).Error

	return
}

// GetBatchFromStaffGroupID 批量查找
func (obj *_StaffGroupRelMgr) GetBatchFromStaffGroupID(staffGroupIDs []int) (results []*StaffGroupRel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(StaffGroupRel{}).Where("`staff_group_id` IN (?)", staffGroupIDs).Find(&results).Error

	return
}

// GetFromStaffID 通过staff_id获取内容
func (obj *_StaffGroupRelMgr) GetFromStaffID(staffID int) (results []*StaffGroupRel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(StaffGroupRel{}).Where("`staff_id` = ?", staffID).Find(&results).Error

	return
}

// GetBatchFromStaffID 批量查找
func (obj *_StaffGroupRelMgr) GetBatchFromStaffID(staffIDs []int) (results []*StaffGroupRel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(StaffGroupRel{}).Where("`staff_id` IN (?)", staffIDs).Find(&results).Error

	return
}

func (obj *_StaffGroupRelMgr) CreateStaffGroupRel(bean *StaffGroupRel) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(StaffGroupRel{}).Create(bean).Error

	return
}

func (obj *_StaffGroupRelMgr) UpdateStaffGroupRel(bean *StaffGroupRel) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(bean).Updates(bean).Error

	return
}

func (obj *_StaffGroupRelMgr) DeleteStaffGroupRel(bean *StaffGroupRel) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(StaffGroupRel{}).Delete(bean).Error

	return
}
