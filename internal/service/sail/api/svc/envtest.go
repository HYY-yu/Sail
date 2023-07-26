package svc

import (
	"fmt"
	model_envtest "github.com/HYY-yu/sail/internal/service/sail/api/envtest/model"
	"github.com/HYY-yu/sail/internal/service/sail/model"
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"net/http"
)

// 检查或新建以下资源：
//   - 一个测试命名空间（可加密、可发布）
//   - 一个测试项目
//   - 一个测试命名空间公共配置文件（加密）
//   - 一个测试的项目内配置文件（关联公共配置）
//   - 一个测试的项目内配置文件（不关联公共配置）

type ResourceType int

const (
	ResourceNamespace ResourceType = iota + 1
	ResourceProject
	ResourcePublicConfig
	ResourceConfig
)

var resourceOrders = []ResourceType{
	ResourceNamespace,
	ResourceProject,
	ResourcePublicConfig,
	ResourceConfig,
}

type TestDataSvc struct {
	BaseSvc
	DB db.Repo

	NamespaceSvc *NamespaceSvc
	ProjectSvc   *ProjectSvc
	ConfigSvc    *ConfigSvc
}

func NewTestDataSvc(
	db db.Repo,
	ns *NamespaceSvc,
	ps *ProjectSvc,
	cs *ConfigSvc,
) *TestDataSvc {
	svc := &TestDataSvc{
		DB:           db,
		NamespaceSvc: ns,
		ProjectSvc:   ps,
		ConfigSvc:    cs,
	}
	return svc
}

func (s *TestDataSvc) CreateTestData(sctx core.SvcContext) error {
	ctx := sctx.Context()

	_, role := s.CheckStaffGroup(ctx, model_envtest.TestProjectGroupId)
	if role > model.RoleManager {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}

	for _, e := range resourceOrders {
		if err := s.createOrClean(e, sctx, false); err != nil {
			return err
		}
	}
	return nil
}

func (s *TestDataSvc) CleanTestData(sctx core.SvcContext) error {
	ctx := sctx.Context()

	_, role := s.CheckStaffGroup(ctx, model_envtest.TestProjectGroupId)
	if role > model.RoleManager {
		return response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}

	// 逆序删除
	for i := len(resourceOrders) - 1; i >= 0; i-- {
		if err := s.createOrClean(resourceOrders[i], sctx, true); err != nil {
			return err
		}
	}
	return nil
}

// createOrClean 一种外观模式
func (s *TestDataSvc) createOrClean(r ResourceType, sctx core.SvcContext, isClean bool) error {
	switch r {
	case ResourceNamespace:
		if isClean {
			return s.cleanNamespace(sctx)
		}
		// 检查命名空间，没有则新建
		return s.checkNamespaceExistOrNew(sctx)
	case ResourceConfig:
		if isClean {
			return s.cleanConfig(sctx)
		}
		return s.checkConfigExistOrNew(sctx)
	case ResourceProject:
		if isClean {
			return s.cleanProject(sctx)
		}
		// 检查项目
		return s.checkProjectExistOrNew(sctx)
	case ResourcePublicConfig:
		if isClean {
			return s.cleanPublicConfig(sctx)
		}
		return s.checkPublicConfigExistOrNew(sctx)
	default:
		return fmt.Errorf("not have resource type %v", r)
	}
}

func (s *TestDataSvc) cleanConfig(sctx core.SvcContext) error {
	ctx := sctx.Context()
	configMgr := s.ConfigSvc.ConfigRepo.Mgr(ctx, s.DB.GetDb())
	namespaceMgr := s.NamespaceSvc.NamespaceRepo.Mgr(ctx, s.DB.GetDb())
	projectMgr := s.ProjectSvc.ProjectRepo.Mgr(ctx, s.DB.GetDb())

	resultPj, err := projectMgr.WithSelects(
		model.ProjectColumns.ID,
	).WithOptions(
		projectMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		projectMgr.WithName(model_envtest.TestProjectName),
	).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	resultNs, err := namespaceMgr.WithSelects(
		model.NamespaceColumns.ID,
	).WithOptions(
		namespaceMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		namespaceMgr.WithName(model_envtest.TestNamespaceName),
	).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	// 不关联公共配置
	resultCfgNotLink, err := configMgr.WithSelects(
		model.ConfigColumns.ID,
	).WithOptions(
		configMgr.WithProjectID(resultPj.ID),
		configMgr.WithNamespaceID(resultNs.ID),
		configMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		configMgr.WithName(model_envtest.TestProjectConfigName),
		configMgr.WithIsPublic(false),
	).Get()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if resultCfgNotLink.ID > 0 {
		err := s.ConfigSvc.Del(sctx, resultCfgNotLink.ID)
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}
	// 关联公共配置
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	resultLinkCfg, err := configMgr.WithSelects(
		model.ConfigColumns.ID,
	).WithOptions(
		configMgr.WithProjectID(resultPj.ID),
		configMgr.WithNamespaceID(resultNs.ID),
		configMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		configMgr.WithIsLinkPublic(true),
		configMgr.WithName(model_envtest.TestProjectConfigLinkPublic),
	).Get()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if resultLinkCfg.ID > 0 {
		err := s.ConfigSvc.Del(sctx, resultLinkCfg.ID)
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}
	return nil
}

func (s *TestDataSvc) cleanPublicConfig(sctx core.SvcContext) error {
	ctx := sctx.Context()
	configMgr := s.ConfigSvc.ConfigRepo.Mgr(ctx, s.DB.GetDb())
	namespaceMgr := s.NamespaceSvc.NamespaceRepo.Mgr(ctx, s.DB.GetDb())

	resultNs, err := namespaceMgr.WithSelects(
		model.NamespaceColumns.ID,
	).WithOptions(
		namespaceMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		namespaceMgr.WithName(model_envtest.TestNamespaceName),
	).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	resultPublicCfg, err := configMgr.WithSelects(
		model.ConfigColumns.ID,
	).WithOptions(
		configMgr.WithProjectID(0),
		configMgr.WithNamespaceID(resultNs.ID),
		configMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		configMgr.WithName(model_envtest.TestPublicConfigName),
		configMgr.WithIsPublic(true),
	).Get()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	if resultPublicCfg.ID > 0 {
		err = s.ConfigSvc.Del(sctx, resultPublicCfg.ID)
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}
	return nil
}

func (s *TestDataSvc) cleanProject(sctx core.SvcContext) error {
	ctx := sctx.Context()
	projectMgr := s.ProjectSvc.ProjectRepo.Mgr(ctx, s.DB.GetDb())

	resultPj, err := projectMgr.WithSelects(
		model.NamespaceColumns.ID,
	).WithOptions(
		projectMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		projectMgr.WithName(model_envtest.TestProjectName),
	).Get()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if resultPj.ID > 0 {
		err = projectMgr.DeleteProject(&resultPj)
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}
	return nil
}

func (s *TestDataSvc) cleanNamespace(sctx core.SvcContext) error {
	ctx := sctx.Context()
	namespaceMgr := s.NamespaceSvc.NamespaceRepo.Mgr(ctx, s.DB.GetDb())

	resultNs, err := namespaceMgr.WithSelects(
		model.NamespaceColumns.ID,
	).WithOptions(
		namespaceMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		namespaceMgr.WithName(model_envtest.TestNamespaceName),
	).Get()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	if resultNs.ID > 0 {
		err = namespaceMgr.DeleteNamespace(&resultNs)
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}
	return nil
}

func (s *TestDataSvc) checkConfigExistOrNew(sctx core.SvcContext) error {
	ctx := sctx.Context()
	configMgr := s.ConfigSvc.ConfigRepo.Mgr(ctx, s.DB.GetDb())
	namespaceMgr := s.NamespaceSvc.NamespaceRepo.Mgr(ctx, s.DB.GetDb())
	projectMgr := s.ProjectSvc.ProjectRepo.Mgr(ctx, s.DB.GetDb())

	resultPj, err := projectMgr.WithSelects(
		model.ProjectColumns.ID,
	).WithOptions(
		projectMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		projectMgr.WithName(model_envtest.TestProjectName),
	).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	resultNs, err := namespaceMgr.WithSelects(
		model.NamespaceColumns.ID,
	).WithOptions(
		namespaceMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		namespaceMgr.WithName(model_envtest.TestNamespaceName),
	).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	// 不关联公共配置
	has, err := configMgr.WithOptions(
		configMgr.WithProjectID(resultPj.ID),
		configMgr.WithNamespaceID(resultNs.ID),
		configMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		configMgr.WithName(model_envtest.TestProjectConfigName),
		configMgr.WithIsPublic(false),
	).HasRecord()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if !has {
		err := s.ConfigSvc.Add(sctx, &model.AddConfig{
			ProjectGroupID: model_envtest.TestProjectGroupId,
			ProjectID:      resultPj.ID,
			NamespaceID:    resultNs.ID,
			Name:           model_envtest.TestProjectConfigName,
			IsPublic:       false,
			Type:           model_envtest.TestConfigType,
			Content:        model_envtest.TestProjectConfigContent,
		})
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}
	// 关联公共配置
	resultPublicCfg, err := configMgr.WithSelects(
		model.ConfigColumns.ID,
	).WithOptions(
		configMgr.WithProjectID(0),
		configMgr.WithNamespaceID(resultNs.ID),
		configMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		configMgr.WithName(model_envtest.TestPublicConfigName),
		configMgr.WithIsPublic(true),
	).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	has, err = configMgr.WithOptions(
		configMgr.WithProjectID(resultPj.ID),
		configMgr.WithNamespaceID(resultNs.ID),
		configMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		configMgr.WithIsLinkPublic(true),
		configMgr.WithName(model_envtest.TestProjectConfigLinkPublic),
	).HasRecord()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if !has {
		err := s.ConfigSvc.Add(sctx, &model.AddConfig{
			ProjectGroupID: model_envtest.TestProjectGroupId,
			ProjectID:      resultPj.ID,
			NamespaceID:    resultNs.ID,
			Name:           model_envtest.TestProjectConfigName,
			IsLinkPublic:   true,
			Type:           model_envtest.TestConfigType,
			Content:        model_envtest.TestProjectConfigContent,
			PublicConfigID: resultPublicCfg.ID,
		})
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}
	return nil
}

func (s *TestDataSvc) checkPublicConfigExistOrNew(sctx core.SvcContext) error {
	ctx := sctx.Context()
	publicConfigMgr := s.ConfigSvc.ConfigRepo.Mgr(ctx, s.DB.GetDb())
	namespaceMgr := s.NamespaceSvc.NamespaceRepo.Mgr(ctx, s.DB.GetDb())

	resultNs, err := namespaceMgr.WithSelects(
		model.NamespaceColumns.ID,
	).WithOptions(
		namespaceMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		namespaceMgr.WithName(model_envtest.TestNamespaceName),
	).Catch()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	has, err := publicConfigMgr.WithOptions(
		publicConfigMgr.WithProjectID(0),
		publicConfigMgr.WithNamespaceID(resultNs.ID),
		publicConfigMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		publicConfigMgr.WithName(model_envtest.TestPublicConfigName),
		publicConfigMgr.WithIsPublic(true),
	).HasRecord()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if !has {
		err := s.ConfigSvc.Add(sctx, &model.AddConfig{
			ProjectGroupID: model_envtest.TestProjectGroupId,
			ProjectID:      0,
			NamespaceID:    resultNs.ID,
			Name:           model_envtest.TestPublicConfigName,
			IsEncrypt:      true,
			IsPublic:       true,
			Type:           model_envtest.TestConfigType,
			Content:        model_envtest.TestPublicConfigContent,
		})
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}
	return nil
}

func (s *TestDataSvc) checkProjectExistOrNew(sctx core.SvcContext) error {
	ctx := sctx.Context()
	projectMgr := s.ProjectSvc.ProjectRepo.Mgr(ctx, s.DB.GetDb())

	has, err := projectMgr.WithOptions(
		projectMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		projectMgr.WithName(model_envtest.TestProjectName),
	).HasRecord()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if !has {
		err := s.ProjectSvc.Add(sctx, &model.AddProject{
			ProjectGroupID: model_envtest.TestProjectGroupId,
			Name:           model_envtest.TestProjectName,
		})
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}
	return nil
}

func (s *TestDataSvc) checkNamespaceExistOrNew(sctx core.SvcContext) error {
	ctx := sctx.Context()
	namespaceMgr := s.NamespaceSvc.NamespaceRepo.Mgr(ctx, s.DB.GetDb())

	has, err := namespaceMgr.WithOptions(
		namespaceMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		namespaceMgr.WithName(model_envtest.TestNamespaceName),
	).HasRecord()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if !has {
		err := s.NamespaceSvc.Add(sctx, &model.AddNamespace{
			ProjectGroupID: model_envtest.TestProjectGroupId,
			Name:           model_envtest.TestNamespaceName,
			RealTime:       true,
			Secret:         true,
		})
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}
	return nil
}
