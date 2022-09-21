package svc

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/encrypt"
	"github.com/HYY-yu/seckill.pkg/pkg/mysqlerr_helper"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/gogf/gf/v2/errors/gerror"
	"gorm.io/gorm"

	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
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
}

func NewConfigSvc(
	db db.Repo,
	store storage.Repo,
	cr repo.ConfigRepo,
	ch repo.ConfigHistoryRepo,
	cl repo.ConfigLinkRepo,
	pr repo.ProjectRepo,
	nr repo.NamespaceRepo,
) *ConfigSvc {
	svc := &ConfigSvc{
		DB:                db,
		ConfigRepo:        cr,
		ConfigHistoryRepo: ch,
		ConfigLinkRepo:    cl,
		Store:             store,
		ProjectRepo:       pr,
		NamespaceRepo:     nr,
	}
	return svc
}

func (s *ConfigSvc) Tree(sctx core.SvcContext, projectID int, projectGroupID int) ([]model.ProjectTree, error) {
	ctx := sctx.Context()
	namespaceMgr := s.NamespaceRepo.Mgr(ctx, s.DB.GetDb())
	namespaceMgr.UpdateDB(namespaceMgr.WithPrepareStmt())
	mgr := s.ConfigRepo.Mgr(ctx, s.DB.GetDb())
	mgr.UpdateDB(mgr.WithPrepareStmt())

	_, role := s.CheckStaffGroup(ctx, projectGroupID)
	if role > model.RoleManager {
		return nil, response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}

	namespaceList, err := namespaceMgr.WithOptions(namespaceMgr.WithProjectGroupID(projectGroupID)).Gets()
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	configList, err := mgr.WithOptions(mgr.WithProjectID(projectID)).WithSelects(
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
		})
	}

	tree := make([]model.ProjectTree, len(namespaceList))
	for i, e := range namespaceList {
		b := model.ProjectTree{
			NamespaceID: e.ID,
			Name:        e.Name,
			RealTime:    e.RealTime,
		}

		b.Nodes = configNamespaceMap[e.ID]
		tree[i] = b
	}
	return tree, nil
}

func (s *ConfigSvc) Info(sctx core.SvcContext, configID int, projectGroupID int) (*model.ConfigInfo, error) {
	ctx := sctx.Context()
	namespaceMgr := s.NamespaceRepo.Mgr(ctx, s.DB.GetDb())
	mgr := s.ConfigRepo.Mgr(ctx, s.DB.GetDb())

	_, role := s.CheckStaffGroup(ctx, projectGroupID)
	if role > model.RoleManager {
		return nil, response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}

	// 如果是owner，则自动解密

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

	err = s.addConfigHistory(ctx, tx, configID, revision, int(userId))
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	tx.Commit()
	return nil
}

func (s *ConfigSvc) LinkPublicConfig(ctx context.Context, configID int, publicConfigID int) error {
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

func (s *ConfigSvc) UnlinkPublicConfig(ctx context.Context, configID int, publicConfigID int) {
	mgr := s.ConfigLinkRepo.Mgr(ctx, s.DB.GetDb())
	mgr.WithOptions(mgr.WithConfigID(configID), mgr.WithPublicConfigID(publicConfigID)).Delete(&model.ConfigLink{})

	configMgr := s.ConfigRepo.Mgr(ctx, s.DB.GetDb())
	configMgr.WithOptions(configMgr.WithID(configID)).Update(model.ConfigColumns.IsLinkPublic, false)
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

func (s *ConfigSvc) addConfigHistory(ctx context.Context, db *gorm.DB, configID int, revision int, userId int) error {
	bean := &model.ConfigHistory{
		ConfigID:   configID,
		Reversion:  revision,
		CreateTime: time.Now(),
		CreateBy:   userId,
	}

	hMgr := s.ConfigHistoryRepo.Mgr(ctx, db)
	return hMgr.CreateConfigHistory(bean)
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
		encryptContent, err := s.EncryptConfigContent(param.Content, namespace.SecretKey)
		if err != nil {
			return 0, 0, err
		}
		param.Content = encryptContent
	}
	if param.IsLinkPublic && !param.IsPublic {
		// 取公共配置内容
		publicConfig, err := mgr.WithOptions(mgr.WithID(param.PublicConfigID)).WithSelects(
			model.ConfigColumns.ID,
			model.ConfigColumns.Name,
			model.ConfigColumns.ConfigType,
		).Catch()
		if err != nil {
			return 0, 0, err
		}

		publicConfigKey := s.ConfigKey(param.IsPublic, project.Key, namespace.Name, publicConfig.Name, model.ConfigType(publicConfig.ConfigType))
		gresp := s.Store.Get(ctx, publicConfigKey)

		param.Content = gresp.Value

		// 不需要在 addConfig 的事务里
		err = s.LinkPublicConfig(ctx, configID, param.PublicConfigID)
		if err != nil {
			return 0, 0, err
		}
	}

	err = mgr.CreateConfig(bean)
	if err != nil {
		return 0, 0, gerror.Wrap(err, "addConfig")
	}

	// 写入 ETCD
	configKey := s.ConfigKey(bean.IsPublic, project.Key, namespace.Name, bean.Name, param.Type)
	sresp := s.Store.Set(ctx, configKey, param.Content)
	if sresp.Err != nil {
		return 0, 0, gerror.Wrap(sresp.Err, "addConfig")
	}

	return bean.ID, sresp.Revision, nil
}

func (s *ConfigSvc) getConfigProjectAndNamespace(ctx context.Context, projectID int, namespaceID int) (*model.Project, *model.Namespace, error) {
	pMgr := s.ProjectRepo.Mgr(ctx, s.DB.GetDb())
	nMgr := s.NamespaceRepo.Mgr(ctx, s.DB.GetDb())

	project, err := pMgr.WithOptions(pMgr.WithID(projectID)).
		WithSelects(model.ProjectColumns.ID, model.ProjectColumns.Name, model.ProjectColumns.Key).Catch()
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

// ConfigKey
// Key : /conf/project_key/namespace_name/config_name.config.type
func (s *ConfigSvc) ConfigKey(isPublic bool, projectKey string, namespaceName string, configName string, configType model.ConfigType) string {
	builder := strings.Builder{}
	builder.WriteString("/conf")

	builder.WriteByte('/')
	builder.WriteString(projectKey)

	builder.WriteByte('/')
	builder.WriteString(namespaceName)

	if isPublic {
		builder.WriteByte('/')
		builder.WriteString("public")
	}

	builder.WriteByte('/')
	builder.WriteString(configName)

	builder.WriteByte('.')
	builder.WriteString(string(configType))

	return builder.String()
}
