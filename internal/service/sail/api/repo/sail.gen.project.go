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
func ProjectMgr(db *gorm.DB) *_ProjectMgr {
	if db == nil {
		panic(fmt.Errorf("ProjectMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &_ProjectMgr{_BaseMgr: &_BaseMgr{DB: db.Table("project"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

// WithContext set context to db
func (obj *_ProjectMgr) WithContext(c context.Context) *_ProjectMgr {
	if c != nil {
		obj.ctx = c
	}
	return obj
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
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	obj.DB = obj.DB.Where(options.query)
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

// WithID id获取
func (obj *_ProjectMgr) WithID(id int) Option {
	return optionFunc(func(o *options) { o.query["id"] = id })
}

// WithProjectGroupID project_group_id获取
func (obj *_ProjectMgr) WithProjectGroupID(projectGroupID int) Option {
	return optionFunc(func(o *options) { o.query["project_group_id"] = projectGroupID })
}

// WithKey key获取
func (obj *_ProjectMgr) WithKey(key string) Option {
	return optionFunc(func(o *options) { o.query["key"] = key })
}

// WithName name获取
func (obj *_ProjectMgr) WithName(name string) Option {
	return optionFunc(func(o *options) { o.query["name"] = name })
}

// WithCreateTime create_time获取
func (obj *_ProjectMgr) WithCreateTime(createTime time.Time) Option {
	return optionFunc(func(o *options) { o.query["create_time"] = createTime })
}

// WithCreateBy create_by获取
func (obj *_ProjectMgr) WithCreateBy(createBy int) Option {
	return optionFunc(func(o *options) { o.query["create_by"] = createBy })
}

// WithDeleteTime delete_time获取
func (obj *_ProjectMgr) WithDeleteTime(deleteTime int) Option {
	return optionFunc(func(o *options) { o.query["delete_time"] = deleteTime })
}

// GetFromID 通过id获取内容
func (obj *_ProjectMgr) GetFromID(id int) (result model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找
func (obj *_ProjectMgr) GetBatchFromID(ids []int) (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromProjectGroupID 通过project_group_id获取内容
func (obj *_ProjectMgr) GetFromProjectGroupID(projectGroupID int) (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`project_group_id` = ?", projectGroupID).Find(&results).Error

	return
}

// GetBatchFromProjectGroupID 批量查找
func (obj *_ProjectMgr) GetBatchFromProjectGroupID(projectGroupIDs []int) (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`project_group_id` IN (?)", projectGroupIDs).Find(&results).Error

	return
}

// GetFromKey 通过key获取内容
func (obj *_ProjectMgr) GetFromKey(key string) (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`key` = ?", key).Find(&results).Error

	return
}

// GetBatchFromKey 批量查找
func (obj *_ProjectMgr) GetBatchFromKey(keys []string) (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`key` IN (?)", keys).Find(&results).Error

	return
}

// GetFromName 通过name获取内容
func (obj *_ProjectMgr) GetFromName(name string) (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`name` = ?", name).Find(&results).Error

	return
}

// GetBatchFromName 批量查找
func (obj *_ProjectMgr) GetBatchFromName(names []string) (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`name` IN (?)", names).Find(&results).Error

	return
}

// GetFromCreateTime 通过create_time获取内容
func (obj *_ProjectMgr) GetFromCreateTime(createTime time.Time) (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`create_time` = ?", createTime).Find(&results).Error

	return
}

// GetBatchFromCreateTime 批量查找
func (obj *_ProjectMgr) GetBatchFromCreateTime(createTimes []time.Time) (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`create_time` IN (?)", createTimes).Find(&results).Error

	return
}

// GetFromCreateBy 通过create_by获取内容
func (obj *_ProjectMgr) GetFromCreateBy(createBy int) (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`create_by` = ?", createBy).Find(&results).Error

	return
}

// GetBatchFromCreateBy 批量查找
func (obj *_ProjectMgr) GetBatchFromCreateBy(createBys []int) (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`create_by` IN (?)", createBys).Find(&results).Error

	return
}

// GetFromDeleteTime 通过delete_time获取内容
func (obj *_ProjectMgr) GetFromDeleteTime(deleteTime int) (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`delete_time` = ?", deleteTime).Find(&results).Error

	return
}

// GetBatchFromDeleteTime 批量查找
func (obj *_ProjectMgr) GetBatchFromDeleteTime(deleteTimes []int) (results []*model.Project, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Project{}).Where("`delete_time` IN (?)", deleteTimes).Find(&results).Error

	return
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
