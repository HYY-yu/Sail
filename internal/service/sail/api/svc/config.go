package svc

import (
	"context"
	"errors"
	"github.com/HYY-yu/sail/internal/service/sail/api/svc_interface"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/encrypt"
	"github.com/HYY-yu/seckill.pkg/pkg/mysqlerr_helper"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/gogf/gf/v2/errors/gerror"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
	"github.com/HYY-yu/sail/internal/service/sail/config"
	"github.com/HYY-yu/sail/internal/service/sail/model"
	"github.com/HYY-yu/sail/internal/service/sail/storage"
)

type ConfigSvc struct {
	BaseSvc
	DB    db.Repo
	Store storage.Repo

	publishSystem svc_interface.PublishSystem

	ConfigRepo        repo.ConfigRepo
	ConfigHistoryRepo repo.ConfigHistoryRepo
	ConfigLinkRepo    repo.ConfigLinkRepo

	ProjectRepo   repo.ProjectRepo
	NamespaceRepo repo.NamespaceRepo
	StaffRepo     repo.StaffRepo
}

func NewConfigSvc(
	db db.Repo,
	store storage.Repo,
	cr repo.ConfigRepo,
	ch repo.ConfigHistoryRepo,
	cl repo.ConfigLinkRepo,
	pr repo.ProjectRepo,
	nr repo.NamespaceRepo,
	sr repo.StaffRepo,
) *ConfigSvc {
	svc := &ConfigSvc{
		DB:                db,
		ConfigRepo:        cr,
		ConfigHistoryRepo: ch,
		ConfigLinkRepo:    cl,
		Store:             store,
		ProjectRepo:       pr,
		NamespaceRepo:     nr,
		StaffRepo:         sr,
	}
	return svc
}

func (s *ConfigSvc) SimplePublicTree(projectGroupID int) map[int][]model.ConfigNode {
	mgr := s.ConfigRepo.Mgr(context.Background(), s.DB.GetDb())
	mgr.WithPrepareStmt()

	configList, _ := mgr.WithOptions(
		mgr.WithProjectID(0),
		mgr.WithProjectGroupID(projectGroupID),
	).WithSelects(
		model.ConfigColumns.ID,
		model.ConfigColumns.Name,
		model.ConfigColumns.NamespaceID,
		model.ConfigColumns.ConfigType,
	).
		Gets()
	configNamespaceMap := make(map[int][]model.ConfigNode)
	for _, e := range configList {
		configNamespaceMap[e.NamespaceID] = append(configNamespaceMap[e.NamespaceID], model.ConfigNode{
			ConfigID: e.ID,
			Name:     e.Name,
			Type:     e.ConfigType,
			Title:    e.Name + "." + e.ConfigType,
		})
	}
	return configNamespaceMap
}

func (s *ConfigSvc) Tree(sctx core.SvcContext, projectID int, projectGroupID int) ([]model.ProjectTree, error) {
	ctx := sctx.Context()
	namespaceMgr := s.NamespaceRepo.Mgr(ctx, s.DB.GetDb())
	namespaceMgr.WithPrepareStmt()
	mgr := s.ConfigRepo.Mgr(ctx, s.DB.GetDb())
	mgr.WithPrepareStmt()

	_, role := s.CheckStaffGroup(ctx, projectGroupID)
	if role > model.RoleManager {
		return nil, response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}
	if projectID == 0 && role > model.RoleOwner {
		// projectID == 0 说明在访问公共配置，此时需要Owner权限
		return nil, response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此数据",
		)
	}

	namespaceList, err := namespaceMgr.
		WithOptions(namespaceMgr.WithProjectGroupID(projectGroupID), namespaceMgr.WithDeleteTime(0)).
		Gets()
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	configList, err := mgr.WithOptions(
		mgr.WithProjectID(projectID),
		mgr.WithProjectGroupID(projectGroupID),
	).WithSelects(
		model.ConfigColumns.ID,
		model.ConfigColumns.Name,
		model.ConfigColumns.NamespaceID,
		model.ConfigColumns.ConfigType,
	).
		Gets()
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	configNamespaceMap := make(map[int][]model.ConfigNode)
	for _, e := range configList {
		configNamespaceMap[e.NamespaceID] = append(configNamespaceMap[e.NamespaceID], model.ConfigNode{
			ConfigID: e.ID,
			Name:     e.Name,
			Type:     e.ConfigType,
			Title:    e.Name + "." + e.ConfigType,
		})
	}

	tree := make([]model.ProjectTree, len(namespaceList))
	for i, e := range namespaceList {
		title := e.Name
		if !e.RealTime {
			title += " 📣 "
		}
		if len(e.SecretKey) > 0 {
			title += " 🔐 "
		}
		// TODO 检测待发布状态

		b := model.ProjectTree{
			NamespaceID: e.ID,
			Name:        e.Name,
			RealTime:    e.RealTime,
			CanSecret:   e.SecretKey != "",
			Spread:      true,
			Title:       title,
		}

		b.Nodes = configNamespaceMap[e.ID]
		tree[i] = b
	}
	return tree, nil
}

func (s *ConfigSvc) Info(sctx core.SvcContext, configID int) (*model.ConfigInfo, error) {
	ctx := sctx.Context()
	mgr := s.ConfigRepo.Mgr(ctx, s.DB.GetDb())
	cfg, err := mgr.WithOptions(mgr.WithID(configID)).Catch()
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	_, role := s.CheckStaffGroup(ctx, cfg.ProjectGroupID)
	if err := s.roleCheck(&cfg, role); err != nil {
		return nil, err
	}

	project, namespace, err := s.getConfigProjectAndNamespace(ctx, cfg.ProjectID, cfg.NamespaceID)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	configKey := s.ConfigKey(
		cfg.IsPublic,
		cfg.ProjectGroupID,
		project.Key,
		namespace.Name,
		cfg.Name,
		model.ConfigType(cfg.ConfigType),
	)
	info := &model.ConfigInfo{
		ConfigID:     cfg.ID,
		ConfigKey:    configKey,
		Name:         cfg.Name,
		Type:         cfg.ConfigType,
		IsPublic:     cfg.IsPublic,
		IsLinkPublic: cfg.IsLinkPublic,
		IsEncrypt:    cfg.IsEncrypt,
	}

	gresp := s.Store.Get(ctx, configKey)
	if gresp.Err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	info.Content = gresp.Value
	// 如果是owner，则自动解密
	if cfg.IsEncrypt && role <= model.RoleOwner {
		decrypt, err := s.DecryptConfigContent(info.Content, namespace.SecretKey)
		if err != nil {
			return nil, response.NewError(
				http.StatusInternalServerError,
				response.ServerError,
				"文件解码失败",
			).WithErr(err)
		}
		info.Content = decrypt
	}

	return info, nil
}

func (s *ConfigSvc) Add(sctx core.SvcContext, param *model.AddConfig) error {
	ctx := sctx.Context()
	userId := sctx.UserId()

	_, role := s.CheckStaffGroup(ctx, param.ProjectGroupID)
	if role > model.RoleManager {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}

	if !param.Type.Valid() {
		return response.NewErrorWithStatusOk(
			response.ParamBindError,
			"请传正确的Type",
		)
	}

	if !param.Content.Valid(param.Type) {
		return response.NewErrorWithStatusOk(
			response.ParamBindError,
			"您上传的配置格式有误，请仔细检查。",
		)
	}

	// 只有RoleOwner可以加密\创建公共配置
	if param.IsEncrypt || param.IsPublic {
		if role > model.RoleOwner {
			return response.NewErrorWithStatusOk(
				response.AuthorizationError,
				"没有权限访问此接口",
			)
		}
	}

	tx := s.DB.GetDb().Begin()
	defer tx.Rollback()

	configID, revision, err := s.addConfig(ctx, tx, param)
	if err != nil {
		if err == model.ErrNotEncryptNamespace {
			return response.NewErrorWithStatusOk(
				response.ParamBindError,
				"此命名空间未配置加密密钥，无法加密",
			)
		}
		if mysqlerr_helper.IsMysqlDupEntryError(err) {
			return response.NewErrorWithStatusOk(
				response.ParamBindError,
				"已经存在相同的配置路径，请保证唯一",
			)
		}
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	err = s.addConfigHistory(ctx, tx, configID, revision, int(userId), model.OpTypeAdd)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	tx.Commit()
	return nil
}

func (s *ConfigSvc) History(sctx core.SvcContext, configID int) ([]model.ConfigHistoryList, error) {
	ctx := sctx.Context()
	historyMgr := s.ConfigHistoryRepo.Mgr(ctx, s.DB.GetDb())
	mgr := s.ConfigRepo.Mgr(ctx, s.DB.GetDb())

	cfg, err := mgr.WithOptions(mgr.WithID(configID)).Catch()
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	_, role := s.CheckStaffGroup(ctx, cfg.ProjectGroupID)
	if err := s.roleCheck(&cfg, role); err != nil {
		return nil, err
	}

	var ch []model.ConfigHistory
	err = historyMgr.WithOptions(historyMgr.WithConfigID(configID)).
		Order(model.ConfigHistoryColumns.Reversion + " desc").
		Limit(int(config.Get().Server.HistoryListLen)).Find(&ch).Error
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	result := make([]model.ConfigHistoryList, len(ch))
	for i, e := range ch {
		b := model.ConfigHistoryList{
			ConfigID:     e.ConfigID,
			CreateBy:     e.CreateBy,
			CreateByName: s.GetCreateByName(ctx, s.DB, s.StaffRepo, e.CreateBy),
			CreateTime:   e.CreateTime.Unix(),
			Reversion:    e.Reversion,
			OpType:       int(e.OpType),
			OpTypeStr:    model.ConfigHistoryOpType(e.OpType).String(),
		}
		result[i] = b
	}
	return result, nil
}

func (s *ConfigSvc) HistoryInfo(sctx core.SvcContext, configID int, reversion int) (string, error) {
	ctx := sctx.Context()
	mgr := s.ConfigRepo.Mgr(ctx, s.DB.GetDb())

	cfg, err := mgr.WithOptions(mgr.WithID(configID)).Catch()
	if err != nil {
		return "", response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	_, role := s.CheckStaffGroup(ctx, cfg.ProjectGroupID)
	if err := s.roleCheck(&cfg, role); err != nil {
		return "", err
	}

	project, namespace, err := s.getConfigProjectAndNamespace(ctx, cfg.ProjectID, cfg.NamespaceID)
	if err != nil {
		return "", response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	configKey := s.ConfigKey(
		cfg.IsPublic,
		cfg.ProjectGroupID,
		project.Key,
		namespace.Name,
		cfg.Name,
		model.ConfigType(cfg.ConfigType),
	)
	gresp := s.Store.GetWithReversion(ctx, configKey, reversion)
	if gresp.Err != nil {
		if err := s.etcdErrCompacted(ctx, gresp.Err, configID, reversion); err != nil {
			return "", err
		}
	}

	// 解密
	if cfg.IsEncrypt && role <= model.RoleOwner {
		decrypt, err := s.DecryptConfigContent(gresp.Value, namespace.SecretKey)
		if err != nil {
			return "", response.NewError(
				http.StatusInternalServerError,
				response.ServerError,
				"文件解码失败",
			).WithErr(err)
		}
		gresp.Value = decrypt
	}
	return gresp.Value, nil
}

func (s *ConfigSvc) Rollback(sctx core.SvcContext, param *model.RollbackConfig) error {
	ctx := sctx.Context()
	userId := int(sctx.UserId())
	mgr := s.ConfigRepo.Mgr(ctx, s.DB.GetDb())

	cfg, err := mgr.WithOptions(mgr.WithID(param.ConfigID)).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	_, role := s.CheckStaffGroup(ctx, cfg.ProjectGroupID)
	if err := s.roleCheck(&cfg, role); err != nil {
		return err
	}

	project, namespace, err := s.getConfigProjectAndNamespace(ctx, cfg.ProjectID, cfg.NamespaceID)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	configKey := s.ConfigKey(
		cfg.IsPublic,
		cfg.ProjectGroupID,
		project.Key,
		namespace.Name,
		cfg.Name,
		model.ConfigType(cfg.ConfigType),
	)
	gresp := s.Store.GetWithReversion(ctx, configKey, param.Reversion)
	if gresp.Err != nil {
		if err := s.etcdErrCompacted(ctx, gresp.Err, param.ConfigID, param.Reversion); err != nil {
			return err
		}
	}

	// 用获取的Content替换现有的版本
	// 不需要解密
	sresp := s.Store.Set(ctx, configKey, gresp.Value)
	if sresp.Err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	// 新增一条历史
	err = s.addConfigHistory(ctx, s.DB.GetDb(), cfg.ID, sresp.Revision, userId, model.OpTypeRollback)
	if sresp.Err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (s *ConfigSvc) etcdErrCompacted(ctx context.Context, err error, configID int, reversion int) error {
	if err == rpctypes.ErrCompacted {
		// 糟糕，ETCD已经把这个 reversion 压缩了，我们把这个 <= reversion 的记录删掉把
		historyMgr := s.ConfigHistoryRepo.Mgr(ctx, s.DB.GetDb())

		historyMgr.WithOptions(
			historyMgr.WithConfigID(configID),
			historyMgr.WithReversion(reversion, " <= ?"),
		).Delete(&model.ConfigHistory{})

		return response.NewErrorWithStatusOk(
			10012,
			"此版本在底层存储中被清除，请尝试其它版本！",
		)
	}
	return nil
}

// Copy
// 可以转为副本：即公共配置的更新不影响自己，此时自己可以编辑
// 可以关联公共配置：即把自己与公共配置重新关联，此时不能编辑了
// 注意：这个操作不会吧 config_link 表的记录删除，只是把flag改变，所以可以重新关联（或者叫再次关联）
// 如果一个配置创建时没关联公共配置，那么它是不能重新关联公共配置的。
func (s *ConfigSvc) Copy(sctx core.SvcContext, param *model.ConfigCopy) error {
	ctx := sctx.Context()
	userId := int(sctx.UserId())

	mgr := s.ConfigRepo.Mgr(ctx, s.DB.GetDb())
	linkMgr := s.ConfigLinkRepo.Mgr(ctx, s.DB.GetDb())
	isLink := false

	cfg, err := mgr.WithOptions(mgr.WithID(param.ConfigID)).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	_, role := s.CheckStaffGroup(ctx, cfg.ProjectGroupID)
	if err := s.roleCheck(&cfg, role); err != nil {
		return err
	}

	link, err := linkMgr.WithOptions(linkMgr.WithConfigID(cfg.ID)).Get()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if link.ID == 0 {
		// 没有公共配置关联记录
		return response.NewErrorWithStatusOk(
			response.ServerError,
			"此记录没有关联过公共配置",
		)
	}

	switch param.Op {
	case 1:
		// 转为副本
		isLink = false
	case 2:
		// 转为公共配置
		isLink = true
		// 重新关联需要用公共配置内容做一次覆盖
		publicConfig, err := mgr.WithOptions(mgr.WithID(link.PublicConfigID)).WithSelects(
			model.ConfigColumns.ID,
			model.ConfigColumns.Name,
			model.ConfigColumns.ConfigType,
			model.ConfigColumns.IsEncrypt,
		).Catch()
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
		project, namespace, err := s.getConfigProjectAndNamespace(ctx, cfg.ProjectID, cfg.NamespaceID)
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}

		publicConfigKey := s.ConfigKey(
			true,
			cfg.ProjectGroupID,
			project.Key,
			namespace.Name,
			publicConfig.Name,
			model.ConfigType(publicConfig.ConfigType),
		)
		gresp := s.Store.Get(ctx, publicConfigKey)

		configKey := s.ConfigKey(
			cfg.IsPublic,
			cfg.ProjectGroupID,
			project.Key,
			namespace.Name,
			cfg.Name,
			model.ConfigType(cfg.ConfigType),
		)
		grespConfig := s.Store.Get(ctx, configKey)
		if grespConfig.Err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
		rollbackETCD := func(s storage.Repo) {
			s.Set(ctx, configKey, grespConfig.Value)
		}

		sresp := s.Store.Set(ctx, configKey, gresp.Value)
		if sresp.Err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
		err = s.addConfigHistory(ctx, s.DB.GetDb(), cfg.ID, sresp.Revision, userId, model.OpTypeLink)
		if err != nil {
			rollbackETCD(s.Store)
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}
	err = mgr.WithOptions(mgr.WithID(param.ConfigID)).
		Update(model.ConfigColumns.IsLinkPublic, isLink).Error
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (s *ConfigSvc) Del(sctx core.SvcContext, configID int) error {
	ctx := sctx.Context()

	tx := s.DB.GetDb().Begin()
	defer tx.Rollback()

	historyMgr := s.ConfigHistoryRepo.Mgr(ctx, tx)
	linkMgr := s.ConfigLinkRepo.Mgr(ctx, tx)
	mgr := s.ConfigRepo.Mgr(ctx, tx)

	cfg, err := mgr.WithOptions(mgr.WithID(configID)).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	_, role := s.CheckStaffGroup(ctx, cfg.ProjectGroupID)
	if err := s.roleCheck(&cfg, role); err != nil {
		return err
	}

	err = historyMgr.WithOptions(historyMgr.WithConfigID(configID)).
		Delete(&model.ConfigHistory{}).Error
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	if !cfg.IsPublic {
		err = linkMgr.WithOptions(linkMgr.WithConfigID(configID)).
			Delete(&model.ConfigLink{}).Error
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	} else {
		// 公共配置还需要删除 config_link ，并将其中的config is_link_public 解绑
		cl, _ := linkMgr.WithOptions(linkMgr.WithPublicConfigID(cfg.ID)).Gets()
		for _, e := range cl {
			cfg, _ := mgr.WithOptions(mgr.WithID(e.ConfigID)).Get()
			if cfg.ID == 0 {
				continue
			}
			if !cfg.IsLinkPublic {
				continue
			}

			err = linkMgr.WithOptions(linkMgr.WithConfigID(cfg.ID), linkMgr.WithPublicConfigID(configID)).Delete(&model.ConfigLink{}).Error
			if err != nil {
				return response.NewErrorAutoMsg(
					http.StatusInternalServerError,
					response.ServerError,
				).WithErr(err)
			}
		}
	}

	project, namespace, err := s.getConfigProjectAndNamespace(ctx, cfg.ProjectID, cfg.NamespaceID)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	configKey := s.ConfigKey(cfg.IsPublic, cfg.ProjectGroupID, project.Key, namespace.Name, cfg.Name, model.ConfigType(cfg.ConfigType))
	_, err = s.Store.Del(ctx, configKey)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	err = mgr.WithOptions(mgr.WithID(cfg.ID)).Delete(&model.Config{}).Error
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	tx.Commit()
	return nil
}

// Edit
// 改公共配置，需要把Link的配置全部改掉
// 要保证这个操作的原子性
func (s *ConfigSvc) Edit(sctx core.SvcContext, param *model.EditConfig) error {
	ctx := sctx.Context()
	userId := int(sctx.UserId())

	mgr := s.ConfigRepo.Mgr(ctx, s.DB.GetDb())

	cfg, err := mgr.WithOptions(mgr.WithID(param.ConfigID)).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	_, role := s.CheckStaffGroup(ctx, cfg.ProjectGroupID)
	if err := s.roleCheck(&cfg, role); err != nil {
		return err
	}

	if !param.Content.Valid(model.ConfigType(cfg.ConfigType)) {
		return response.NewErrorWithStatusOk(
			response.ParamBindError,
			"您上传的配置格式有误，请仔细检查。",
		)
	}

	if cfg.IsLinkPublic {
		// 禁止编辑
		return response.NewErrorWithStatusOk(
			response.ParamBindError,
			"此配置关联到公共配置，无法编辑",
		)
	}

	project, namespace, err := s.getConfigProjectAndNamespace(ctx, cfg.ProjectID, cfg.NamespaceID)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	tx := s.DB.GetDb().Begin()
	defer tx.Rollback()
	paramContent := string(param.Content)

	_, err = s.editConfig(ctx, userId, tx, paramContent, &cfg, project, namespace)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	if cfg.IsPublic {
		// 更新ConfigLink涉及的配置
		configLinkMgr := s.ConfigLinkRepo.Mgr(ctx, s.DB.GetDb())
		cl, _ := configLinkMgr.WithOptions(configLinkMgr.WithPublicConfigID(cfg.ID)).Gets()
		updateKeys := make([]string, 0)
		updateValues := make([]string, 0)
		var bErr error

		// 通过 ETCD 事务保证更新的原子性
		for _, e := range cl {
			cfg, _ := mgr.WithOptions(mgr.WithID(e.ConfigID)).Get()
			if cfg.ID == 0 {
				continue
			}
			if !cfg.IsLinkPublic {
				continue
			}
			project, namespace, _ := s.getConfigProjectAndNamespace(ctx, cfg.ProjectID, cfg.NamespaceID)
			configKey := s.ConfigKey(
				cfg.IsPublic,
				cfg.ProjectGroupID,
				project.Key,
				namespace.Name,
				cfg.Name,
				model.ConfigType(cfg.ConfigType),
			)
			updateKeys = append(updateKeys, configKey)

			if cfg.IsEncrypt {
				encryptContent, err := s.EncryptConfigContent(paramContent, namespace.SecretKey)
				if err != nil {
					bErr = err
					break
				}
				paramContent = encryptContent
			}
			updateValues = append(updateValues, paramContent)
		}

		if bErr != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(bErr)
		}

		setResp := s.Store.AtomicBatchSet(ctx, updateKeys, updateValues)
		if setResp.Err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(setResp.Err)
		}
	}

	tx.Commit()
	return nil
}

func (s *ConfigSvc) editConfig(
	ctx context.Context,
	userId int,
	tx *gorm.DB,
	newContent string,
	config *model.Config,
	project *model.Project,
	namespace *model.Namespace,
) (func(s storage.Repo), error) {
	// 更新到ETCD
	configKey := s.ConfigKey(
		config.IsPublic,
		config.ProjectGroupID,
		project.Key,
		namespace.Name,
		config.Name,
		model.ConfigType(config.ConfigType),
	)
	gresp := s.Store.Get(ctx, configKey)
	if gresp.Err != nil {
		return nil, gresp.Err
	}
	if len(gresp.Value) == 0 {
		return nil, errors.New("not found key: " + configKey)
	}
	rollbackETCD := func(s storage.Repo) {
		s.Set(ctx, configKey, gresp.Value)
	}
	if config.IsEncrypt {
		encryptContent, err := s.EncryptConfigContent(newContent, namespace.SecretKey)
		if err != nil {
			return nil, errors.New("encryptContent: " + configKey)
		}
		newContent = encryptContent
	}

	sresp := s.Store.Set(ctx, configKey, newContent)
	if sresp.Err != nil {
		return nil, sresp.Err
	}

	err := s.addConfigHistory(ctx, tx, config.ID, sresp.Revision, userId, model.OpTypeEdit)
	if err != nil {
		rollbackETCD(s.Store)
		return nil, err
	}
	return rollbackETCD, nil
}

func (s *ConfigSvc) linkPublicConfig(ctx context.Context, configID int, publicConfigID int) error {
	bean := &model.ConfigLink{
		ConfigID:       configID,
		PublicConfigID: publicConfigID,
	}
	mgr := s.ConfigLinkRepo.Mgr(ctx, s.DB.GetDb())

	err := mgr.CreateConfigLink(bean)
	if err != nil {
		return gerror.Wrap(err, "")
	}
	return nil
}

func (s *ConfigSvc) unlinkPublicConfig(ctx context.Context, configID int, publicConfigID int) {
	mgr := s.ConfigLinkRepo.Mgr(ctx, s.DB.GetDb())
	mgr.WithOptions(mgr.WithConfigID(configID), mgr.WithPublicConfigID(publicConfigID)).Delete(&model.ConfigLink{})

	configMgr := s.ConfigRepo.Mgr(ctx, s.DB.GetDb())
	configMgr.WithOptions(configMgr.WithID(configID)).Update(model.ConfigColumns.IsLinkPublic, false)
}

func (s *ConfigSvc) addConfigHistory(ctx context.Context, db *gorm.DB, configID int, revision int, userId int, opType model.ConfigHistoryOpType) error {
	bean := &model.ConfigHistory{
		ConfigID:   configID,
		Reversion:  revision,
		CreateTime: time.Now(),
		CreateBy:   userId,
		OpType:     int8(opType),
	}

	hMgr := s.ConfigHistoryRepo.Mgr(ctx, db)
	err := hMgr.CreateConfigHistory(bean)
	if err != nil {
		return err
	}

	// 检查历史长度
	cfCount, err := hMgr.
		WithOptions(hMgr.WithConfigID(configID)).Count()
	if err != nil {
		return nil
	}
	if cfCount > config.Get().Server.HistoryListLen+rand.Int63n(5) {
		// 删除
		deta := cfCount - config.Get().Server.HistoryListLen
		hMgr.WithOptions(hMgr.WithConfigID(configID)).
			Order(model.ConfigHistoryColumns.Reversion).Limit(int(deta)).Delete(&model.ConfigHistory{})
	}
	return nil
}

func (s *ConfigSvc) addConfig(ctx context.Context, db *gorm.DB, param *model.AddConfig) (configID int, revision int, err error) {
	bean := &model.Config{
		Name:           param.Name,
		ProjectID:      param.ProjectID,
		ProjectGroupID: param.ProjectGroupID,
		NamespaceID:    param.NamespaceID,
		IsPublic:       param.IsPublic,
		IsLinkPublic:   param.IsLinkPublic,
		IsEncrypt:      param.IsEncrypt,
		ConfigType:     string(param.Type),
	}
	mgr := s.ConfigRepo.Mgr(ctx, db)

	project, namespace, err := s.getConfigProjectAndNamespace(ctx, bean.ProjectID, bean.NamespaceID)
	if err != nil {
		return 0, 0, gerror.Wrap(err, "addConfig")
	}

	if param.IsEncrypt {
		encryptContent, err := s.EncryptConfigContent(string(param.Content), namespace.SecretKey)
		if err != nil {
			return 0, 0, err
		}
		param.Content = model.ConfigContent(encryptContent)
	}

	err = mgr.CreateConfig(bean)
	if err != nil {
		return 0, 0, gerror.Wrap(err, "addConfig")
	}

	if param.IsLinkPublic && !param.IsPublic {
		// 取公共配置内容
		publicConfig, err := mgr.WithOptions(mgr.WithID(param.PublicConfigID)).WithSelects(
			model.ConfigColumns.ID,
			model.ConfigColumns.Name,
			model.ConfigColumns.ConfigType,
			model.ConfigColumns.IsEncrypt,
		).Catch()
		if err != nil {
			return 0, 0, err
		}

		publicConfigKey := s.ConfigKey(
			true,
			param.ProjectGroupID,
			project.Key,
			namespace.Name,
			publicConfig.Name,
			model.ConfigType(publicConfig.ConfigType),
		)
		gresp := s.Store.Get(ctx, publicConfigKey)

		param.Content = model.ConfigContent(gresp.Value)
		param.Type = model.ConfigType(publicConfig.ConfigType)

		// 不需要在 addConfig 的事务里
		err = s.linkPublicConfig(ctx, bean.ID, param.PublicConfigID)
		if err != nil {
			return 0, 0, err
		}

		mgr.WithOptions(mgr.WithID(bean.ID)).Updates(map[string]interface{}{
			model.ConfigColumns.IsEncrypt:  publicConfig.IsEncrypt,
			model.ConfigColumns.ConfigType: publicConfig.ConfigType,
		})
	}

	// 写入 ETCD
	configKey := s.ConfigKey(bean.IsPublic, bean.ProjectGroupID, project.Key, namespace.Name, bean.Name, param.Type)
	sresp := s.Store.Set(ctx, configKey, string(param.Content))
	if sresp.Err != nil {
		return 0, 0, gerror.Wrap(sresp.Err, "addConfig")
	}

	return bean.ID, sresp.Revision, nil
}

func (s *ConfigSvc) getConfigProjectAndNamespace(ctx context.Context, projectID int, namespaceID int) (*model.Project, *model.Namespace, error) {
	pMgr := s.ProjectRepo.Mgr(ctx, s.DB.GetDb())
	nMgr := s.NamespaceRepo.Mgr(ctx, s.DB.GetDb())
	pMgr.WithPrepareStmt()
	nMgr.WithPrepareStmt()

	// projectID 有可能为 0，用 Get 而不是 Catch
	project, err := pMgr.WithOptions(pMgr.WithID(projectID)).
		WithSelects(model.ProjectColumns.ID, model.ProjectColumns.Name, model.ProjectColumns.Key).Get()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "getConfigProjectAndNamespace")
	}
	namespace, err := nMgr.WithOptions(pMgr.WithID(namespaceID)).
		WithSelects(model.NamespaceColumns.ID, model.NamespaceColumns.Name, model.NamespaceColumns.SecretKey).Catch()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "getConfigProjectAndNamespace")
	}

	return &project, &namespace, nil
}

func (s *ConfigSvc) roleCheck(cfg *model.Config, role model.Role) error {
	// 普通配置只有RoleManager可以访问
	if role > model.RoleManager {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}
	// 公共配置只有RoleOwner可以访问
	if cfg.IsPublic && role > model.RoleOwner {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}
	// 加密配置只有RoleOwner可以访问
	if cfg.IsEncrypt && role > model.RoleOwner {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"您没有权限修改加密配置",
		)
	}
	return nil
}

// 实现 ConfigSystem

func (s *ConfigSvc) EncryptConfigContent(content string, namespaceKey string) (string, error) {
	if namespaceKey == "" {
		return "", model.ErrNotEncryptNamespace
	}

	goAES := encrypt.NewGoAES(namespaceKey, encrypt.AES192)
	encryptContent, err := goAES.WithModel(encrypt.ECB).WithEncoding(encrypt.NewBase64Encoding()).Encrypt(content)
	if err != nil {
		return "", err
	}
	return encryptContent, nil
}

func (s *ConfigSvc) DecryptConfigContent(content string, namespaceKey string) (string, error) {
	if namespaceKey == "" {
		return "", model.ErrNotEncryptNamespace
	}

	goAES := encrypt.NewGoAES(namespaceKey, encrypt.AES192)
	decryptContent, err := goAES.WithModel(encrypt.ECB).WithEncoding(encrypt.NewBase64Encoding()).Decrypt(content)
	if err != nil {
		return "", err
	}
	return decryptContent, nil
}

func (s *ConfigSvc) SetPublishSystem(ps svc_interface.PublishSystem) {
	s.publishSystem = ps
}

// ConfigKey
// NormalConfigKey : /conf/project_key/namespace_name/config_name.config.type
// PublicConfigKey : /conf/1-public/namespace_name/config_name.config.type
func (s *ConfigSvc) ConfigKey(isPublic bool, projectGroupID int, projectKey string, namespaceName string, configName string, configType model.ConfigType) string {
	builder := strings.Builder{}
	builder.WriteString("/conf")

	if !isPublic {
		builder.WriteByte('/')
		builder.WriteString(projectKey)
	}
	if isPublic {
		builder.WriteByte('/')
		builder.WriteString(strconv.Itoa(projectGroupID) + "-public")
	}

	builder.WriteByte('/')
	builder.WriteString(namespaceName)

	builder.WriteByte('/')
	builder.WriteString(configName)

	builder.WriteByte('.')
	builder.WriteString(string(configType))

	return builder.String()
}

func (s *ConfigSvc) GetConfig(ctx context.Context, configID int) (*model.Config, error) {
	mgr := s.ConfigRepo.Mgr(ctx, s.DB.GetDb())
	cfg, err := mgr.WithOptions(mgr.WithID(configID)).Catch()
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	return &cfg, nil
}

func (s *ConfigSvc) ConfigEdit() {

}

func (s *ConfigSvc) GetConfigProjectAndNamespace(ctx context.Context, projectID int, namespaceID int) (*model.Project, *model.Namespace, error) {
	return s.getConfigProjectAndNamespace(ctx, projectID, namespaceID)
}
