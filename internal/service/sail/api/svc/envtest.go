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

func (s *TestDataSvc) GetTestData(sctx core.SvcContext) (map[string]interface{}, error) {
	ctx := sctx.Context()

	_, role := s.CheckStaffGroup(ctx, model_envtest.TestProjectGroupId)
	if role > model.RoleManager {
		return nil, response.NewErrorWithStatusOk(
			response.AuthorizationError,
			"没有权限访问此接口",
		)
	}

	result := make(map[string]interface{})
	// ProjectId
	testPjId, err := s.findTestProjectId(sctx)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	result["project_id"] = testPjId
	// namespaceId
	testNamespaceId, err := s.findTestNamespaceId(sctx)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	result["namespace_id"] = testNamespaceId
	// publicConfigId
	testPublicConfigId, err := s.findTestPublicConfigId(sctx, testNamespaceId)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	result["public_config_id"] = testPublicConfigId
	// config_link_id
	testConfigLinkId, err := s.findTestLinkPublicConfigId(sctx, testPjId, testNamespaceId)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	result["config_link_id"] = testConfigLinkId
	// config_id
	testConfigId, err := s.findTestConfigId(sctx, testPjId, testNamespaceId)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	result["config_id"] = testConfigId

	return result, nil
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

func (s *TestDataSvc) findTestProjectId(sctx core.SvcContext) (int, error) {
	ctx := sctx.Context()
	projectMgr := s.ProjectSvc.ProjectRepo.Mgr(ctx, s.DB.GetDb())

	resultPj, err := projectMgr.WithSelects(
		model.ProjectColumns.ID,
	).WithOptions(
		projectMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		projectMgr.WithName(model_envtest.TestProjectName),
	).Catch()
	return resultPj.ID, err
}

func (s *TestDataSvc) findTestNamespaceId(sctx core.SvcContext) (int, error) {
	ctx := sctx.Context()
	namespaceMgr := s.NamespaceSvc.NamespaceRepo.Mgr(ctx, s.DB.GetDb())

	resultNs, err := namespaceMgr.WithSelects(
		model.NamespaceColumns.ID,
	).WithOptions(
		namespaceMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		namespaceMgr.WithName(model_envtest.TestNamespaceName),
	).Catch()

	return resultNs.ID, err
}

func (s *TestDataSvc) findTestLinkPublicConfigId(sctx core.SvcContext, testProjectId, testNamespaceId int) (int, error) {
	ctx := sctx.Context()
	configMgr := s.ConfigSvc.ConfigRepo.Mgr(ctx, s.DB.GetDb())

	resultCfgNotLink, err := configMgr.WithSelects(
		model.ConfigColumns.ID,
	).WithOptions(
		configMgr.WithProjectID(testProjectId),
		configMgr.WithNamespaceID(testNamespaceId),
		configMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		configMgr.WithName(model_envtest.TestProjectConfigLinkPublic),
		configMgr.WithIsLinkPublic(true),
	).Get()
	return resultCfgNotLink.ID, err
}

func (s *TestDataSvc) findTestConfigId(sctx core.SvcContext, testProjectId, testNamespaceId int) (int, error) {
	ctx := sctx.Context()
	configMgr := s.ConfigSvc.ConfigRepo.Mgr(ctx, s.DB.GetDb())

	resultCfgNotLink, err := configMgr.WithSelects(
		model.ConfigColumns.ID,
	).WithOptions(
		configMgr.WithProjectID(testProjectId),
		configMgr.WithNamespaceID(testNamespaceId),
		configMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		configMgr.WithName(model_envtest.TestProjectConfigName),
		configMgr.WithIsPublic(false),
	).Get()
	return resultCfgNotLink.ID, err
}

func (s *TestDataSvc) findTestPublicConfigId(sctx core.SvcContext, testNamespaceId int) (int, error) {
	ctx := sctx.Context()
	configMgr := s.ConfigSvc.ConfigRepo.Mgr(ctx, s.DB.GetDb())

	resultPublicCfg, err := configMgr.WithSelects(
		model.ConfigColumns.ID,
	).WithOptions(
		configMgr.WithProjectID(0),
		configMgr.WithNamespaceID(testNamespaceId),
		configMgr.WithProjectGroupID(model_envtest.TestProjectGroupId),
		configMgr.WithName(model_envtest.TestPublicConfigName),
		configMgr.WithIsPublic(true),
	).Get()
	return resultPublicCfg.ID, err
}

func (s *TestDataSvc) cleanConfig(sctx core.SvcContext) error {
	resultPj, err := s.findTestProjectId(sctx)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	resultNs, err := s.findTestNamespaceId(sctx)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	// 不关联公共配置
	resultCfgNotLinkID, err := s.findTestConfigId(sctx, resultPj, resultNs)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if resultCfgNotLinkID > 0 {
		err := s.ConfigSvc.Del(sctx, resultCfgNotLinkID)
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}
	// 关联公共配置
	resultLinkCfgID, err := s.findTestLinkPublicConfigId(sctx, resultPj, resultNs)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if resultLinkCfgID > 0 {
		err := s.ConfigSvc.Del(sctx, resultLinkCfgID)
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
	resultNsID, err := s.findTestNamespaceId(sctx)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	resultPublicCfgId, err := s.findTestPublicConfigId(sctx, resultNsID)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	if resultPublicCfgId > 0 {
		err = s.ConfigSvc.Del(sctx, resultPublicCfgId)
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

	resultPjId, err := s.findTestProjectId(sctx)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if resultPjId > 0 {
		err = projectMgr.DeleteProject(&model.Project{ID: resultPjId})
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

	resultNsId, err := s.findTestNamespaceId(sctx)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	if resultNsId > 0 {
		err = namespaceMgr.DeleteNamespace(&model.Namespace{ID: resultNsId})
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

	resultPjId, err := s.findTestProjectId(sctx)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	resultNsId, err := s.findTestNamespaceId(sctx)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	// 不关联公共配置
	has, err := configMgr.WithOptions(
		configMgr.WithProjectID(resultPjId),
		configMgr.WithNamespaceID(resultNsId),
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
			ProjectID:      resultPjId,
			NamespaceID:    resultNsId,
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
	resultPublicCfgId, err := s.findTestPublicConfigId(sctx, resultNsId)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	has, err = configMgr.WithOptions(
		configMgr.WithProjectID(resultPjId),
		configMgr.WithNamespaceID(resultNsId),
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
			ProjectID:      resultPjId,
			NamespaceID:    resultNsId,
			Name:           model_envtest.TestProjectConfigLinkPublic,
			IsLinkPublic:   true,
			Type:           model_envtest.TestConfigType,
			Content:        model_envtest.TestProjectConfigContent,
			PublicConfigID: resultPublicCfgId,
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

	resultNsId, err := s.findTestNamespaceId(sctx)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	has, err := publicConfigMgr.WithOptions(
		publicConfigMgr.WithProjectID(0),
		publicConfigMgr.WithNamespaceID(resultNsId),
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
			NamespaceID:    resultNsId,
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
