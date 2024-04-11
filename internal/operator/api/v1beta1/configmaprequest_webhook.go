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

package v1beta1

import (
	"errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var configmaprequestlog = logf.Log.WithName("configmaprequest-resource")

func (r *ConfigMapRequest) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-cmr-sail-hyy-yu-space-v1beta1-configmaprequest,mutating=true,failurePolicy=fail,sideEffects=None,groups=cmr.sail.hyy-yu.space,resources=configmaprequests,verbs=create;update,versions=v1beta1,name=mconfigmaprequest.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &ConfigMapRequest{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *ConfigMapRequest) Default() {
	configmaprequestlog.Info("default", "name", r.Name)

	if r.Spec.Watch == nil {
		r.Spec.Watch = new(bool)
		*r.Spec.Watch = true
	}

	if r.Spec.Merge != nil {
		if *r.Spec.Merge {
			// default MergeConfig is: config.toml

			if r.Spec.MergeConfigFile == nil {
				r.Spec.MergeConfigFile = new(string)
				*r.Spec.MergeConfigFile = "config.toml"
			}

			defaultConfig := *r.Spec.MergeConfigFile
			r.Spec.Configs = append(r.Spec.Configs, defaultConfig)
		}
	}
}

//+kubebuilder:webhook:path=/validate-cmr-sail-hyy-yu-space-v1beta1-configmaprequest,mutating=false,failurePolicy=fail,sideEffects=None,groups=cmr.sail.hyy-yu.space,resources=configmaprequests,verbs=create;update;delete,versions=v1beta1,name=vconfigmaprequest.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &ConfigMapRequest{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ConfigMapRequest) ValidateCreate() (warnings admission.Warnings, err error) {
	configmaprequestlog.Info("validate create", "name", r.Name)

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ConfigMapRequest) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	configmaprequestlog.Info("validate update", "name", r.Name)

	// 不允许修改，因为这涉及到 ConfigMapRequest 下的所有 ConfigMap
	if r.Spec.ProjectKey != old.(*ConfigMapRequest).Spec.ProjectKey {
		return nil, errors.New("you can't update the project key, please create a new one. ")
	}
	if r.Spec.Namespace != old.(*ConfigMapRequest).Spec.Namespace {
		return nil, errors.New("you can't update the namespace, please create a new one. ")
	}

	// 不能修改 merged 后的配置，对已创建的 ConfigMap 合并或者拆分，可能会影响到使用它的 Pod。
	if r.Spec.Merge != old.(*ConfigMapRequest).Spec.Merge {
		return nil, errors.New("merge config can't update, because it's hard to split. ")
	}
	if r.Spec.MergeConfigFile != old.(*ConfigMapRequest).Spec.MergeConfigFile {
		return nil, errors.New("merge format can't update, because it's hard to split. ")
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ConfigMapRequest) ValidateDelete() (admission.Warnings, error) {
	configmaprequestlog.Info("validate delete", "name", r.Name)

	configmaprequestlog.Info("Delete the CMR will delete it corresponding ConfigMap, Please careful!")
	return nil, nil
}
