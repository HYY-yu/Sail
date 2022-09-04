package repo

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

// Code generated by gormt. DO NOT EDIT.

type _ProjectMgr struct {
	*_BaseMgr
}

// ProjectMgr open func
func ProjectMgr(ctx context.Context, db *gorm.DB) *_ProjectMgr {
	if db == nil {
		panic(fmt.Errorf("ProjectMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(ctx)
	return &_ProjectMgr{_BaseMgr: &_BaseMgr{DB: db.Table("project"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

func (obj *_ProjectMgr) WithSelects(idName string, selects ...string) *_ProjectMgr {
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

func (obj *_ProjectMgr) WithOmit(omit ...string) *_ProjectMgr {
	if len(omit) > 0 {
		obj.DB = obj.DB.Omit(omit...)
	}
	return obj
}

func (obj *_ProjectMgr) WithOptions(opts ...Option) *_ProjectMgr {
	options := options{
		query: make(map[string]queryData, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	for k, v := range options.query {
		obj.DB = obj.DB.Where(k+" "+v.cond, v.data)
	}
	return obj
}

// GetTableName get sql table name.获取数据库名字
func (obj *_ProjectMgr) GetTableName() string {
	return "project"
}

// Reset 重置gorm会话
func (obj *_ProjectMgr) Reset() *_ProjectMgr {
	obj.new()
	return obj
}

// Get 获取
func (obj *_ProjectMgr) Get() (result model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_ProjectMgr) Gets() (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Find(&results).Error

	return
}

func (obj *_ProjectMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(model.Project{}).Count(count)
}

func (obj *_ProjectMgr) HasRecord() (bool, error) {
	var count int64
	err := obj.DB.WithContext(obj.ctx).Model(model.Project{}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// WithID id获取
func (obj *_ProjectMgr) WithID(id int, cond ...string) Option {
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
func (obj *_ProjectMgr) WithProjectGroupID(projectGroupID int, cond ...string) Option {
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

// WithKey key获取
func (obj *_ProjectMgr) WithKey(key string, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["key"] = queryData{
			cond: cond[0],
			data: key,
		}
	})
}

// WithName name获取
func (obj *_ProjectMgr) WithName(name string, cond ...string) Option {
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

// WithCreateTime create_time获取
func (obj *_ProjectMgr) WithCreateTime(createTime time.Time, cond ...string) Option {
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
func (obj *_ProjectMgr) WithCreateBy(createBy int, cond ...string) Option {
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

// WithDeleteTime delete_time获取
func (obj *_ProjectMgr) WithDeleteTime(deleteTime int, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["delete_time"] = queryData{
			cond: cond[0],
			data: deleteTime,
		}
	})
}

func (obj *_ProjectMgr) CreateProject(bean *model.Project) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Create(bean).Error

	return
}

func (obj *_ProjectMgr) UpdateProject(bean *model.Project) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(bean).Updates(bean).Error

	return
}

func (obj *_ProjectMgr) DeleteProject(bean *model.Project) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Delete(bean).Error

	return
}
