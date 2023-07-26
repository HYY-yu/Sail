// Package envtest 构造一个管理后台的集成测试环境
// 集成测试的接口封装，比如对这些配置文件的操作、测试环境的构建

package envtest

type ApiEnvTest struct {
}

func (*ApiEnvTest) Start() {
	// 1. 登录

	// 2. 调用接口 CreateData
}

func (*ApiEnvTest) Stop() {
	// 调用 CleanData 接口清除测试数据
}
