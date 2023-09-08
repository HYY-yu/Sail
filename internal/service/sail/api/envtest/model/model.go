package model

// 以下数据在测试阶段会自动生成，测试结束自动删除
const (
	TestNamespaceName           = "suite-test"
	TestProjectName             = "TEST"
	TestConfigType              = "yaml"
	TestPublicConfigName        = "test-public"
	TestPublicConfigContent     = "key: value"
	TestProjectConfigName       = "test-project"
	TestProjectConfigContent    = "keyProject: value"
	TestProjectConfigLinkPublic = "test-link"
)

// 依赖默认的 Test 账号，请勿删除
// 依赖默认的项目组
// 没有以下数据，测试将无法进行。
const (
	TestAPIURL     = "http://127.0.0.1:8108"
	TestCheckHeath = "/sail/system/health"
	TestLogin      = "/sail/v1/login"

	TestCreateTestData = "/sail/v1/env_test/create"
	TestCleanTestData  = "/sail/v1/env_test/clean"
	TestGetTestData    = "/sail/v1/env_test/get"

	TestMetaConfig   = "/sail/v1/config/meta"
	TestUpdateConfig = "/sail/v1/config/edit"

	TestAccountName    = "Test"
	TestAccountPass    = "Test123"
	TestProjectGroupId = 1
)
