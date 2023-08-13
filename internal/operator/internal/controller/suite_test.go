/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	envtestAPI "github.com/HYY-yu/sail/internal/service/sail/api/envtest"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	cmrv1beta1 "github.com/HYY-yu/sail/internal/operator/api/v1beta1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var apiEnvTest *envtestAPI.ApiEnvTest

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "`Controller Suite")

}

// 测试用例-更新测试：
// 1. 首先我们创建一个测试 CMR (不关联公共配置)
// 2. 检查是否正确创建了 ConfigMap
// 3. 调用 API 更新这个配置
// 4. 检查 ConfigMap 是否更新成功

// 测试用例-关联配置更新测试
// 1. 首先创建一个测试 CMR （关联公共配置）
// 2. 检查是否正确创建 ConfigMap
// 3. 调用 API 更新公共配置
// 4. 检查 ConfigMap 是否更新成功

// 测试用例-Merge 配置测试
// 1. 创建测试 CMR （两个配置）
// 2. 一个 CMR Merge True ，检查是否只生成了一个 ConfigMap
// 3. 一个 CMR Merge False, 检查是否创建了两个 ConfigMap

// 测试用例-CMR更新
// 1. 创建测试 CMR (不关联公共配置)
// 2. 关闭 Watch
// 3. 调用 API 更新配置
// 4. 检查 ConfigMap 是否不再更新
// 5. 打开 Watch
// 6. 调用 API 更新配置
// 7. ConfigMap 更新了

// 测试用例-CMR 删除
// 1. 创建测试 CMR （不关联公共配置）
// 2. 检查是否正确创建了 ConfigMap
// 3. 删除 CMR
// 4. ConfigMap 被删除

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	// 要进行此集成测试，需要准备好 API 测试环境
	// Start() 检查测试环境是否正确
	var err error

	By("bootstrapping test environment")
	apiEnvTest = new(envtestAPI.ApiEnvTest)
	err = apiEnvTest.Start()
	Expect(err).NotTo(HaveOccurred())

	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = cmrv1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := apiEnvTest.Stop()
	Expect(err).NotTo(HaveOccurred())

	err = testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
