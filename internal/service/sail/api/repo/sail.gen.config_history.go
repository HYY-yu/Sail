package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/model"
)

// Code generated by gormt. DO NOT EDIT.

type _ConfigHistoryMgr struct {
	*_BaseMgr
}

// ConfigHistoryMgr open func
func ConfigHistoryMgr(ctx context.Context, db *gorm.DB) *_ConfigHistoryMgr {
	if db == nil {
		panic(fmt.Errorf("ConfigHistoryMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(ctx)
	return &_ConfigHistoryMgr{_BaseMgr: &_BaseMgr{DB: db.Table("config_history").WithContext(ctx), isRelated: globalIsRelated, ctx: ctx, cancel: cancel}}
}

func (obj *_ConfigHistoryMgr) WithSelects(idName string, selects ...string) *_ConfigHistoryMgr {
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

func (obj *_ConfigHistoryMgr) WithOptions(opts ...Option) *_ConfigHistoryMgr {
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
func (obj *_ConfigHistoryMgr) GetTableName() string {
	return "config_history"
}

// Tx 开启事务会话
func (obj *_ConfigHistoryMgr) Tx(db *gorm.DB) *_ConfigHistoryMgr {
	obj.UpdateDB(db.Table(obj.GetTableName()).WithContext(obj.ctx))
	return obj
}

// Reset 重置gorm会话
func (obj *_ConfigHistoryMgr) Reset() *_ConfigHistoryMgr {
	obj.new()
	return obj
}

// Get 获取
func (obj *_ConfigHistoryMgr) Get() (result model.ConfigHistory, err error) {
	err = obj.DB.Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_ConfigHistoryMgr) Gets() (results []model.ConfigHistory, err error) {
	err = obj.DB.Find(&results).Error

	return
}

// Take 必须获取结果（单条）
func (obj *_ConfigHistoryMgr) Catch() (results model.ConfigHistory, err error) {
	err = obj.DB.Take(&results).Error

	return
}

func (obj *_ConfigHistoryMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.Count(count)
}

func (obj *_ConfigHistoryMgr) HasRecord() (bool, error) {
	var count int64
	err := obj.DB.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// WithID id获取
func (obj *_ConfigHistoryMgr) WithID(id interface{}, cond ...string) Option {
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

// WithConfigID config_id获取
func (obj *_ConfigHistoryMgr) WithConfigID(configID interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["config_id"] = queryData{
			cond: cond[0],
			data: configID,
		}
	})
}

// WithReversion reversion获取
func (obj *_ConfigHistoryMgr) WithReversion(reversion interface{}, cond ...string) Option {
	return optionFunc(func(o *options) {
		if len(cond) == 0 {
			cond = []string{" = ? "}
		}
		o.query["reversion"] = queryData{
			cond: cond[0],
			data: reversion,
		}
	})
}

// WithCreateTime create_time获取
func (obj *_ConfigHistoryMgr) WithCreateTime(createTime interface{}, cond ...string) Option {
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
func (obj *_ConfigHistoryMgr) WithCreateBy(createBy interface{}, cond ...string) Option {
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

func (obj *_ConfigHistoryMgr) CreateConfigHistory(bean *model.ConfigHistory) (err error) {
	err = obj.DB.Create(bean).Error

	return
}

func (obj *_ConfigHistoryMgr) UpdateConfigHistory(bean *model.ConfigHistory) (err error) {
	err = obj.DB.Updates(bean).Error

	return
}

func (obj *_ConfigHistoryMgr) DeleteConfigHistory(bean *model.ConfigHistory) (err error) {
	err = obj.DB.Delete(bean).Error

	return
}
