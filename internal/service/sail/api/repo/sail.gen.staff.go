package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

// Code generated by gormt. DO NOT EDIT.

type _StaffMgr struct {
	*_BaseMgr
}

// StaffMgr open func
func StaffMgr(ctx context.Context, db *gorm.DB) *_StaffMgr {
	if db == nil {
		panic(fmt.Errorf("StaffMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(ctx)
	return &_StaffMgr{_BaseMgr: &_BaseMgr{DB: db.Table("staff").WithContext(ctx), isRelated: globalIsRelated, ctx: ctx, cancel: cancel}}
}

func (obj *_StaffMgr) WithSelects(idName string, selects ...string) *_StaffMgr {
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

func (obj *_StaffMgr) WithOptions(opts ...Option) *_StaffMgr {
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

// GetTableName get sql table name.获取数据库名字
func (obj *_StaffMgr) GetTableName() string {
	return "staff"
}

// Reset 重置gorm会话
func (obj *_StaffMgr) Reset() *_StaffMgr {
	obj.new()
	return obj
}

// Get 获取
func (obj *_StaffMgr) Get() (result model.Staff, err error) {
	err = obj.DB.Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_StaffMgr) Gets() (results []model.Staff, err error) {
	err = obj.DB.Find(&results).Error

	return
}

func (obj *_StaffMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.Count(count)
}

func (obj *_StaffMgr) HasRecord() (bool, error) {
	var count int64
	err := obj.DB.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// WithID id获取
func (obj *_StaffMgr) WithID(id interface{}, cond ...string) Option {
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

// WithName name获取
func (obj *_StaffMgr) WithName(name interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["name"] = queryData{
			cond: cond[0],
			data: name,
		}
	})
}

// WithPassword password获取
func (obj *_StaffMgr) WithPassword(password interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["password"] = queryData{
			cond: cond[0],
			data: password,
		}
	})
}

// WithRefreshToken refresh_token获取
func (obj *_StaffMgr) WithRefreshToken(refreshToken interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["refresh_token"] = queryData{
			cond: cond[0],
			data: refreshToken,
		}
	})
}

// WithCreateTime create_time获取
func (obj *_StaffMgr) WithCreateTime(createTime interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["create_time"] = queryData{
			cond: cond[0],
			data: createTime,
		}
	})
}

// WithCreateBy create_by获取
func (obj *_StaffMgr) WithCreateBy(createBy interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["create_by"] = queryData{
			cond: cond[0],
			data: createBy,
		}
	})
}

func (obj *_StaffMgr) CreateStaff(bean *model.Staff) (err error) {
	err = obj.DB.Create(bean).Error

	return
}

func (obj *_StaffMgr) UpdateStaff(bean *model.Staff) (err error) {
	err = obj.DB.Updates(bean).Error

	return
}

func (obj *_StaffMgr) DeleteStaff(bean *model.Staff) (err error) {
	err = obj.DB.Delete(bean).Error

	return
}
