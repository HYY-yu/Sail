// Package envtest 构造一个管理后台的集成测试环境
// 检查或新建以下资源：
//    - 一个测试命名空间（可加密、可发布）
//    - 一个测试项目
//    - 一个测试命名空间公共配置文件（加密）
//    - 一个测试的项目内配置文件（关联公共配置）
//    - 一个测试的项目内配置文件（不关联公共配置）
// 以及一些方便集成测试的接口封装，比如对这些配置文件的操作

package envtest

// 依赖默认的 Test 账号，请勿删除
const (
	TestAccountName    = "Test"
	TestAccountPass    = "Test123"
	TestProjectGroupId = 1
)

type ApiEnvTest struct {
}

func (*ApiEnvTest) Start() {

}

func (*ApiEnvTest) Stop() {

}
