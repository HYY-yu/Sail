package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

// Code generated by gormt. DO NOT EDIT.

// 非线程安全

type _ConfigMgr struct {
	*_BaseMgr
}

// ConfigMgr open func
func ConfigMgr(ctx context.Context, db *gorm.DB) *_ConfigMgr {
	if db == nil {
		panic(fmt.Errorf("ConfigMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(ctx)
	return &_ConfigMgr{_BaseMgr: &_BaseMgr{DB: db.Table("config").WithContext(ctx), isRelated: globalIsRelated, ctx: ctx, cancel: cancel}}
}

func (obj *_ConfigMgr) WithSelects(idName string, selects ...string) *_ConfigMgr {
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

func (obj *_ConfigMgr) WithOptions(opts ...Option) *_ConfigMgr {
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
func (obj *_ConfigMgr) GetTableName() string {
	return "config"
}

// Tx 开启事务会话
func (obj *_ConfigMgr) Tx(tx *gorm.DB) *_ConfigMgr {
	obj.DB = tx.Table(obj.GetTableName()).WithContext(obj.ctx)
	return obj
}

// WithPrepareStmt 开启语句 PrepareStmt 功能
// 接下来执行的SQL将会是PrepareStmt的
func (obj *_ConfigMgr) WithPrepareStmt() {
	obj.DB = obj.DB.Session(&gorm.Session{Context: obj.ctx, PrepareStmt: true})
}

// Reset 重置gorm会话
func (obj *_ConfigMgr) Reset() *_ConfigMgr {
	obj.DB = obj.DB.Session(&gorm.Session{NewDB: true, Context: obj.ctx}).Table(obj.GetTableName())
	return obj
}

// Get 获取
func (obj *_ConfigMgr) Get() (result model.Config, err error) {
	err = obj.DB.Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_ConfigMgr) Gets() (results []model.Config, err error) {
	err = obj.DB.Find(&results).Error

	return
}

// Catch 必须获取结果（单条）
func (obj *_ConfigMgr) Catch() (results model.Config, err error) {
	err = obj.DB.Take(&results).Error

	return
}

func (obj *_ConfigMgr) Count() (count int64, err error) {
	err = obj.DB.Count(&count).Error

	return
}

func (obj *_ConfigMgr) HasRecord() (bool, error) {
	count, err := obj.Count()
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// WithID id获取
func (obj *_ConfigMgr) WithID(id interface{}, cond ...string) Option {
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
func (obj *_ConfigMgr) WithName(name interface{}, cond ...string) Option {
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

// WithProjectID project_id获取
func (obj *_ConfigMgr) WithProjectID(projectID interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["project_id"] = queryData{
			cond: cond[0],
			data: projectID,
		}
	})
}

// WithProjectGroupID project_group_id获取 公共配置只有project_group_id
func (obj *_ConfigMgr) WithProjectGroupID(projectGroupID interface{}, cond ...string) Option {
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

// WithNamespaceID namespace_id获取
func (obj *_ConfigMgr) WithNamespaceID(namespaceID interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["namespace_id"] = queryData{
			cond: cond[0],
			data: namespaceID,
		}
	})
}

// WithIsPublic is_public获取
func (obj *_ConfigMgr) WithIsPublic(isPublic interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["is_public"] = queryData{
			cond: cond[0],
			data: isPublic,
		}
	})
}

// WithIsLinkPublic is_link_public获取
func (obj *_ConfigMgr) WithIsLinkPublic(isLinkPublic interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["is_link_public"] = queryData{
			cond: cond[0],
			data: isLinkPublic,
		}
	})
}

// WithIsEncrypt is_encrypt获取
func (obj *_ConfigMgr) WithIsEncrypt(isEncrypt interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["is_encrypt"] = queryData{
			cond: cond[0],
			data: isEncrypt,
		}
	})
}

// WithConfigType config_type获取
func (obj *_ConfigMgr) WithConfigType(configType interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["config_type"] = queryData{
			cond: cond[0],
			data: configType,
		}
	})
}

func (obj *_ConfigMgr) CreateConfig(bean *model.Config) (err error) {
	err = obj.DB.Create(bean).Error

	return
}

func (obj *_ConfigMgr) UpdateConfig(bean *model.Config) (err error) {
	err = obj.DB.Updates(bean).Error

	return
}

func (obj *_ConfigMgr) DeleteConfig(bean *model.Config) (err error) {
	err = obj.DB.Delete(bean).Error

	return
}
