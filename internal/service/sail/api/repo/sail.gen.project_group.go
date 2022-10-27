package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

// Code generated by gormt. DO NOT EDIT.

// 非线程安全

type _ProjectGroupMgr struct {
	*_BaseMgr
}

// ProjectGroupMgr open func
func ProjectGroupMgr(ctx context.Context, db *gorm.DB) *_ProjectGroupMgr {
	if db == nil {
		panic(fmt.Errorf("ProjectGroupMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(ctx)
	return &_ProjectGroupMgr{_BaseMgr: &_BaseMgr{DB: db.Table("project_group").WithContext(ctx), isRelated: globalIsRelated, ctx: ctx, cancel: cancel}}
}

func (obj *_ProjectGroupMgr) WithSelects(idName string, selects ...string) *_ProjectGroupMgr {
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

func (obj *_ProjectGroupMgr) WithOptions(opts ...Option) *_ProjectGroupMgr {
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
func (obj *_ProjectGroupMgr) GetTableName() string {
	return "project_group"
}

// Tx 开启事务会话
func (obj *_ProjectGroupMgr) Tx(tx *gorm.DB) *_ProjectGroupMgr {
	obj.DB = tx.Table(obj.GetTableName()).WithContext(obj.ctx)
	return obj
}

// WithPrepareStmt 开启语句 PrepareStmt 功能
// 接下来执行的SQL将会是PrepareStmt的
func (obj *_ProjectGroupMgr) WithPrepareStmt() {
	obj.DB = obj.DB.Session(&gorm.Session{Context: obj.ctx, PrepareStmt: true})
}

// Reset 重置gorm会话
func (obj *_ProjectGroupMgr) Reset() *_ProjectGroupMgr {
	obj.DB = obj.DB.Session(&gorm.Session{NewDB: true, Context: obj.ctx}).Table(obj.GetTableName())
	return obj
}

// Get 获取
func (obj *_ProjectGroupMgr) Get() (result model.ProjectGroup, err error) {
	err = obj.DB.Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_ProjectGroupMgr) Gets() (results []model.ProjectGroup, err error) {
	err = obj.DB.Find(&results).Error

	return
}

// Catch 必须获取结果（单条）
func (obj *_ProjectGroupMgr) Catch() (results model.ProjectGroup, err error) {
	err = obj.DB.Take(&results).Error

	return
}

func (obj *_ProjectGroupMgr) Count() (count int64, err error) {
	err = obj.DB.Count(&count).Error

	return
}

func (obj *_ProjectGroupMgr) HasRecord() (bool, error) {
	count, err := obj.Count()
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// WithID id获取
func (obj *_ProjectGroupMgr) WithID(id interface{}, cond ...string) Option {
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
func (obj *_ProjectGroupMgr) WithName(name interface{}, cond ...string) Option {
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
func (obj *_ProjectGroupMgr) WithCreateTime(createTime interface{}, cond ...string) Option {
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
func (obj *_ProjectGroupMgr) WithCreateBy(createBy interface{}, cond ...string) Option {
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
func (obj *_ProjectGroupMgr) WithDeleteTime(deleteTime interface{}, cond ...string) Option {
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

func (obj *_ProjectGroupMgr) CreateProjectGroup(bean *model.ProjectGroup) (err error) {
	err = obj.DB.Create(bean).Error

	return
}

func (obj *_ProjectGroupMgr) UpdateProjectGroup(bean *model.ProjectGroup) (err error) {
	err = obj.DB.Updates(bean).Error

	return
}

func (obj *_ProjectGroupMgr) DeleteProjectGroup(bean *model.ProjectGroup) (err error) {
	err = obj.DB.Delete(bean).Error

	return
}
