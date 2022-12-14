package svc

import (
	"context"
	"errors"
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
			"???????????????????????????",
		)
	}
	if projectID == 0 && role > model.RoleOwner {
		// projectID == 0 ??????????????????????????????????????????Owner??????
		return nil, response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"???????????????????????????",
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
			title += " ???? "
		}
		if len(e.SecretKey) > 0 {
			title += " ???? "
		}
		// TODO ?????????????????????

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
	// ?????????owner??????????????????
	if cfg.IsEncrypt && role <= model.RoleOwner {
		decrypt, err := s.DecryptConfigContent(info.Content, namespace.SecretKey)
		if err != nil {
			return nil, response.NewError(
				http.StatusInternalServerError,
				response.ServerError,
				"??????????????????",
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
			"???????????????????????????",
		)
	}

	if !param.Type.Valid() {
		return response.NewErrorWithStatusOk(
			response.ParamBindError,
			"???????????????Type",
		)
	}

	if !param.Content.Valid(param.Type) {
		return response.NewErrorWithStatusOk(
			response.ParamBindError,
			"???????????????????????????????????????????????????",
		)
	}

	// ??????RoleOwner????????????\??????????????????
	if param.IsEncrypt || param.IsPublic {
		if role > model.RoleOwner {
			return response.NewErrorWithStatusOk(
				response.AuthorizationError,
				"???????????????????????????",
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
				"???????????????????????????????????????????????????",
			)
		}
		if mysqlerr_helper.IsMysqlDupEntryError(err) {
			return response.NewErrorWithStatusOk(
				response.ParamBindError,
				"???????????????????????????????????????????????????",
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
			OpType:       e.OpType,
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

	// ??????
	if cfg.IsEncrypt && role <= model.RoleOwner {
		decrypt, err := s.DecryptConfigContent(gresp.Value, namespace.SecretKey)
		if err != nil {
			return "", response.NewError(
				http.StatusInternalServerError,
				response.ServerError,
				"??????????????????",
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

	// ????????????Content?????????????????????
	// ???????????????
	sresp := s.Store.Set(ctx, configKey, gresp.Value)
	if sresp.Err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	// ??????????????????
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
		// ?????????ETCD??????????????? reversion ??????????????????????????? <= reversion ??????????????????
		historyMgr := s.ConfigHistoryRepo.Mgr(ctx, s.DB.GetDb())

		historyMgr.WithOptions(
			historyMgr.WithConfigID(configID),
			historyMgr.WithReversion(reversion, " <= ?"),
		).Delete(&model.ConfigHistory{})

		return response.NewErrorWithStatusOk(
			10012,
			"???????????????????????????????????????????????????????????????",
		)
	}
	return nil
}

// Copy
// ???????????????????????????????????????????????????????????????????????????????????????
// ??????????????????????????????????????????????????????????????????????????????????????????
// ?????????????????????????????? config_link ??????????????????????????????flag????????????????????????????????????????????????????????????
// ???????????????????????????????????????????????????????????????????????????????????????????????????
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
		// ??????????????????????????????
		return response.NewErrorWithStatusOk(
			response.ServerError,
			"????????????????????????????????????",
		)
	}

	switch param.Op {
	case 1:
		// ????????????
		isLink = false
	case 2:
		// ??????????????????
		isLink = true
		// ??????????????????????????????????????????????????????
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
		// ??????????????????????????? config_link ??????????????????config is_link_public ??????
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
	err = s.Store.Del(ctx, configKey)
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
// ???????????????????????????Link?????????????????????
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
			"???????????????????????????????????????????????????",
		)
	}

	if cfg.IsLinkPublic {
		// ????????????
		return response.NewErrorWithStatusOk(
			response.ParamBindError,
			"?????????????????????????????????????????????",
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

	rollback, err := s.editConfig(ctx, userId, tx, string(param.Content), &cfg, project, namespace)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	if cfg.IsPublic {
		// ??????ConfigLink???????????????
		configLinkMgr := s.ConfigLinkRepo.Mgr(ctx, s.DB.GetDb())
		cl, _ := configLinkMgr.WithOptions(configLinkMgr.WithPublicConfigID(cfg.ID)).Gets()

		rollbacks := append([]func(s storage.Repo){}, rollback)
		isRoll := false
		var bErr error

		for _, e := range cl {
			cfg, _ := mgr.WithOptions(mgr.WithID(e.ConfigID)).Get()
			if cfg.ID == 0 {
				continue
			}
			if !cfg.IsLinkPublic {
				continue
			}
			project, namespace, _ := s.getConfigProjectAndNamespace(ctx, cfg.ProjectID, cfg.NamespaceID)
			rb, err := s.editConfig(ctx, userId, tx, string(param.Content), &cfg, project, namespace)
			if err != nil {
				bErr = err
				isRoll = true
				break
			}
			rollbacks = append(rollbacks, rb)
		}

		if isRoll {
			for _, f := range rollbacks {
				f(s.Store)
			}
			if bErr != nil {
				return response.NewErrorAutoMsg(
					http.StatusInternalServerError,
					response.ServerError,
				).WithErr(bErr)
			}
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
	// ?????????ETCD
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
		OpType:     int(opType),
	}

	hMgr := s.ConfigHistoryRepo.Mgr(ctx, db)
	err := hMgr.CreateConfigHistory(bean)
	if err != nil {
		return err
	}

	// ??????????????????
	cfCount, err := hMgr.
		WithOptions(hMgr.WithConfigID(configID)).Count()
	if err != nil {
		return nil
	}
	if cfCount > config.Get().Server.HistoryListLen+rand.Int63n(5) {
		// ??????
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
		// ?????????????????????
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

		// ???????????? addConfig ????????????
		err = s.linkPublicConfig(ctx, bean.ID, param.PublicConfigID)
		if err != nil {
			return 0, 0, err
		}

		mgr.WithOptions(mgr.WithID(bean.ID)).Updates(map[string]interface{}{
			model.ConfigColumns.IsEncrypt:  publicConfig.IsEncrypt,
			model.ConfigColumns.ConfigType: publicConfig.ConfigType,
		})
	}

	// ?????? ETCD
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
	// ??????????????????RoleManager????????????
	if role > model.RoleManager {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"???????????????????????????",
		)
	}
	// ??????????????????RoleOwner????????????
	if cfg.IsPublic && role > model.RoleOwner {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"???????????????????????????",
		)
	}
	// ??????????????????RoleOwner????????????
	if cfg.IsEncrypt && role > model.RoleOwner {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"?????????????????????????????????",
		)
	}
	return nil
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
