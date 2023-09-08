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
	"context"
	"github.com/HYY-yu/sail/internal/operator/internal/config_server"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
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

var ctx context.Context
var cancel context.CancelFunc

const namespace = "default"
const (
	etcd_endpoint = "127.0.0.1:2379"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "`Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	// 要进行此集成测试，需要准备好 API 测试环境
	// Start() 检查测试环境是否正确
	var err error
	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	apiEnvTest = new(envtestAPI.ApiEnvTest)
	err = apiEnvTest.Start()
	Expect(err).NotTo(HaveOccurred())

	t := true
	testEnv = &envtest.Environment{
		//CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		//ErrorIfCRDPathMissing: true,
		UseExistingCluster: &t,
	}

	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = cmrv1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             scheme.Scheme,
		MetricsBindAddress: "0",
	})
	Expect(err).ToNot(HaveOccurred())

	configServer := config_server.NewConfigServer(
		k8sManager.GetLogger(),
		k8sManager.GetConfig(),
		config_server.MetaConfig{
			Namespace:     namespace,
			ETCDEndpoints: etcd_endpoint,
		},
	)

	err = k8sManager.Add(configServer.(manager.Runnable))
	Expect(err).ToNot(HaveOccurred())

	err = (&ConfigMapRequestReconciler{
		Client:       k8sManager.GetClient(),
		Scheme:       k8sManager.GetScheme(),
		Namespace:    namespace,
		ConfigServer: configServer,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

	k8sClient = k8sManager.GetClient()
	Expect(k8sClient).ToNot(BeNil())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	cancel()
	err := apiEnvTest.Stop()
	Expect(err).NotTo(HaveOccurred())

	err = testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
