package controller

import (
	"fmt"
	cmrv1beta1 "github.com/HYY-yu/sail/internal/operator/api/v1beta1"
	"github.com/HYY-yu/sail/internal/service/sail/api/envtest/model"
	"github.com/gogf/gf/v2/encoding/gjson"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	pkgClient "sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

const timeout = time.Second * 5
const interval = time.Second * 1

var _ = Describe("CMR Controller Testing... ", func() {
	key := types.NamespacedName{
		Namespace: namespace,
	}
	var secret *v1.Secret
	var cmr *cmrv1beta1.ConfigMapRequest
	Context("Test context.", func() {
		BeforeEach(func() {
			configYaml, err := apiEnvTest.MetaConfigYaml()
			Expect(err).NotTo(HaveOccurred())
			Expect(configYaml).NotTo(BeEmpty())

			configYamlPart := strings.Split(configYaml, "---")
			secretConfigYaml := new(gjson.Json)
			cmrConfigYaml := new(gjson.Json)
			for i, part := range configYamlPart {
				partStr := strings.TrimSpace(part)
				if i == 0 {
					secretConfigYaml, err = gjson.LoadYaml(partStr)
					Expect(err).NotTo(HaveOccurred())
				}
				if i == 1 {
					cmrConfigYaml, err = gjson.LoadYaml(partStr)
					Expect(err).NotTo(HaveOccurred())
				}
			}

			secret = &v1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: secretConfigYaml.Get("apiVersion").String(),
					Kind:       secretConfigYaml.Get("kind").String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      secretConfigYaml.Get("metadata.name").String(),
				},
				Type: v1.SecretTypeOpaque,
				Data: map[string][]byte{
					"namespace_key": []byte(secretConfigYaml.Get("data.namespace_key").String()),
				},
			}
			Expect(k8sClient.Create(ctx, secret)).NotTo(HaveOccurred())

			watch := true
			merge := cmrConfigYaml.Get("spec.merge").Bool()
			mergeConfigFile := cmrConfigYaml.Get("spec.merge_config_file").String()
			cmr = &cmrv1beta1.ConfigMapRequest{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "",
					Kind:       cmrConfigYaml.Get("kind").String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      cmrConfigYaml.Get("metadata.name").String(),
				},
				Spec: cmrv1beta1.ConfigMapRequestSpec{
					ProjectKey: cmrConfigYaml.Get("spec.project_key").String(),
					Namespace:  cmrConfigYaml.Get("spec.namespace").String(),
					NamespaceKeyInSecret: &v1.LocalObjectReference{
						Name: secretConfigYaml.Get("metadata.name").String(),
					},
					Merge:           &merge,
					MergeConfigFile: &mergeConfigFile,
					Configs:         cmrConfigYaml.Get("spec.configs").Strings(),
					Watch:           &watch,
				},
			}
			key.Name = cmr.Spec.ProjectKey + "-" + cmr.Spec.Namespace
		})

		AfterEach(func() {

			Expect(k8sClient.Delete(ctx, secret)).NotTo(HaveOccurred())
			Expect(k8sClient.Delete(ctx, cmr)).NotTo(HaveOccurred())
			time.Sleep(time.Second * 5) // wait cmr_controller reconcile.
		})

		// 测试用例-更新测试：
		// 1. 首先我们创建一个测试 CMR
		// 2. 检查是否正确创建了 ConfigMap
		// 3. 调用 API 更新配置
		// 4. 检查 ConfigMap 是否更新成功
		// 5. 测试结束后，所有的资源（Secret \ CMR \ ConfigMap） 都应该自动删除了
		It("create cmr and check update... ", func() {
			By("create cmr")
			Expect(k8sClient.Create(ctx, cmr)).NotTo(HaveOccurred())
			time.Sleep(time.Second * 2) // wait cmr_controller reconcile.

			By("check configmap")
			Eventually(func() bool {
				fetched := &v1.ConfigMap{}
				_ = k8sClient.Get(ctx, key, fetched)

				fmt.Println("fetched.Data1: ", fetched)
				return len(fetched.Data) == 2
			}, timeout, interval).Should(BeTrue())

			By("Update config. ")
			newValue := "key: hello"
			err := apiEnvTest.UpdateTestConfig(false, newValue)
			Expect(err).NotTo(HaveOccurred())

			By("check configmap")
			// ConfigMap 应该被更新成新的值。
			Eventually(func() bool {
				fetched := &v1.ConfigMap{}
				_ = k8sClient.Get(ctx, key, fetched)

				data := fetched.Data[model.TestProjectConfigName+"."+model.TestConfigType]
				fmt.Println("fetched.Data2: ", fetched)
				return data == newValue
			}, timeout, interval).Should(BeTrue())

			By("Update public config, and check config-link-public is updated.")
			newValuePublic := "key: public-hello"
			err = apiEnvTest.UpdateTestConfig(true, newValuePublic)
			Expect(err).NotTo(HaveOccurred())

			By("check configmap")
			Eventually(func() bool {
				fetched := &v1.ConfigMap{}
				_ = k8sClient.Get(ctx, key, fetched)

				data := fetched.Data[model.TestProjectConfigLinkPublic+"."+model.TestConfigType]
				fmt.Println("fetched.Data3: ", fetched)
				return data == newValuePublic
			}, timeout, interval).Should(BeTrue())

			// 复原
			err = apiEnvTest.UpdateTestConfig(true, model.TestPublicConfigContent)
			Expect(err).NotTo(HaveOccurred())
			err = apiEnvTest.UpdateTestConfig(false, model.TestProjectConfigContent)
			Expect(err).NotTo(HaveOccurred())
		})

		// 测试用例-Merge 配置测试
		// 1. 创建测试 CMR
		// 2.  CMR Merge True ，检查是否 ConfigMap 里只有一个配置
		It("test merge flag", func() {
			By("create cmr")
			merge := true
			mergeConfigFile := "merge_config.yaml"
			cmr.Spec.Merge = &merge
			cmr.Spec.MergeConfigFile = &mergeConfigFile

			Expect(k8sClient.Create(ctx, cmr)).NotTo(HaveOccurred())
			time.Sleep(time.Second * 2) // wait cmr_controller reconcile.

			By("check configmap")
			Eventually(func() bool {
				fetched := &v1.ConfigMap{}
				_ = k8sClient.Get(ctx, key, fetched)

				fmt.Println("fetched.DataMergeFlag: ", fetched)
				return len(fetched.Data) == 1
			}, timeout, interval).Should(BeTrue())
		})

		// 测试用例-CMR更新
		// 1. 创建测试 CMR (不关联公共配置)
		// 2. 关闭 Watch
		// 3. 调用 API 更新配置
		// 4. 检查 ConfigMap 是否不再更新
		// 5. 打开 Watch
		// 6. 调用 API 更新配置
		// 7. ConfigMap 又可以更新
		It("cmr update checker", func() {
			By("create cmr")
			merge := false
			cmr.Spec.Merge = &merge
			watch := false
			cmr.Spec.Watch = &watch

			Expect(k8sClient.Create(ctx, cmr)).NotTo(HaveOccurred())
			time.Sleep(time.Second * 2) // wait cmr_controller reconcile.

			By("Update public config, and check config-link-public is updated.")
			newValuePublic := "key: public-hello"
			err := apiEnvTest.UpdateTestConfig(false, newValuePublic)
			Expect(err).NotTo(HaveOccurred())

			By("check configmap")
			Consistently(func() bool {
				fetched := &v1.ConfigMap{}
				_ = k8sClient.Get(ctx, key, fetched)

				data := fetched.Data[model.TestProjectConfigName+"."+model.TestConfigType]
				fmt.Println("fetched.DataWatchFlag1: ", fetched)
				return data == model.TestProjectConfigContent // 一直不更新 是对的
			}, timeout, interval).Should(BeTrue())

			// 打开 watch
			Expect(k8sClient.Get(ctx, pkgClient.ObjectKeyFromObject(cmr), cmr)).NotTo(HaveOccurred())
			watch = true
			cmr.Spec.Watch = &watch
			Expect(k8sClient.Update(ctx, cmr)).NotTo(HaveOccurred())
			time.Sleep(time.Second * 2) // wait cmr_controller reconcile.

			By("Update public config, and check config-link-public is updated.")
			newValue := "key: public-hello_22"
			err = apiEnvTest.UpdateTestConfig(false, newValue)
			Expect(err).NotTo(HaveOccurred())

			By("check configmap")
			Eventually(func() bool {
				fetched := &v1.ConfigMap{}
				_ = k8sClient.Get(ctx, key, fetched)

				data := fetched.Data[model.TestProjectConfigName+"."+model.TestConfigType]
				fmt.Println("fetched.DataWatchFlag2: ", fetched)
				return data == newValue
			}, timeout, interval).Should(BeTrue())
		})
	})
})
