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
	TestAPIURL     = "http://127.0.0.1:8108/sail/v1"
	TestCheckHeath = "/system/health"

	TestCreateTestData = "/create_test_data"
	TestCleanTestData  = "/clean_test_data"

	TestAccountName    = "Test"
	TestAccountPass    = "Test123"
	TestProjectGroupId = 1
)
