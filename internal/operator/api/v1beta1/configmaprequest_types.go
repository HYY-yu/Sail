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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConfigMapRequestSpec defines the desired state of ConfigMapRequest
type ConfigMapRequestSpec struct {

	// ProjectKey 配置中心的项目 Key，用于规定从哪个项目取得配置
	ProjectKey string `json:"project_key"`

	// Namespace 项目的命名空间，指定获取哪一个命名空间
	Namespace string `json:"namespace"`

	// NamespaceKey 项目的命名空间秘钥
	// 为了安全起见，NamespaceKey不能直接提供，需要提供一个 Secret 其中保存了NamespaceKey。
	// 格式：namespaceKey: secret-namespace-key-name
	// +optional
	NamespaceKeyInSecret *v1.LocalObjectReference `json:"namespace_key_in_secret,omitempty"`

	// Configs 可选的配置文件列表，不传则获取对应命名空间所有的配置
	// +optional
	Configs []string `json:"configs,omitempty"`

	// Merge 可选配置，默认为 False，当 True 时，会把所有配置文件聚合成一个配置
	// 配置文件的名称将成为新 Key，如：mysql.toml: {data: }
	// 注意：如果 Merge 是 True，则无法再更改为 False，反之亦然。
	// +optional
	Merge *bool `json:"merge,omitempty"`

	// +kubebuilder:validation:Enum=toml;yaml;json;ini;properties;custom

	// MergeFormat 可选配置，默认聚合到一个格式为 config.toml 的 ConfigMap，可以在这里配置其它格式：(xx.yaml\cfg.json等)
	// 支持 json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "tfvars", "dotenv", "env", "ini"
	// +optional
	MergeConfigFile *string `json:"merge_config_file,omitempty"`

	// Watch 可选配置，默认为 True，代表当配置发生变化时，自动更新。
	// +optional
	Watch *bool `json:"watch,omitempty"`
}

// ConfigMapRequestStatus defines the observed state of ConfigMapRequest
type ConfigMapRequestStatus struct {
	// TotalConfig 这个 CMR 下管理的配置文件总数
	// +optional
	TotalConfig int `json:"total_config,omitempty"`

	// ActiveConfig 已经下载并生成 ConfigMap 的数量
	// +optional
	ActiveConfig int `json:"active_config,omitempty"`

	// LastUpdateTime 最近一次的配置更新的时间（如果有多个配置更新，取离当前时间最近的配置更新时间）
	LastUpdateTime *metav1.Time `json:"last_update_time,omitempty"`

	// ManagedConfigList 配置文件的管理列表
	ManagedConfigList []ManagedConfig `json:"managed_config_list"`
}

// ManagedConfig 配置文件实体
type ManagedConfig struct {
	// ConfigFileName 配置文件名
	ConfigFileName string `json:"config_file_name"`

	// LastUpdateTime 这个配置最后更新时间
	LastUpdateTime *metav1.Time `json:"last_update_time,omitempty"`

	// Watched 是否在监听它的配置变更
	Watched bool `json:"watched"`

	// IsDeleted 指示此配置是否已经被删除
	// 为 true 的配置
	IsDeleted bool `json:"is_deleted"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ConfigMapRequest is the Schema for the configmaprequests API
// +kubebuilder:resource:shortName=cmr
// +kubebuilder:printcolumn:name="Active",type="integer",JSONPath=".status.active_config"
// +kubebuilder:printcolumn:name="Total",type="integer",JSONPath=".status.total_config"
// +kubebuilder:printcolumn:name="LastUpdateTIme",type="date",JSONPath=".status.last_update_time"
type ConfigMapRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigMapRequestSpec   `json:"spec,omitempty"`
	Status ConfigMapRequestStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConfigMapRequestList contains a list of ConfigMapRequest
type ConfigMapRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigMapRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConfigMapRequest{}, &ConfigMapRequestList{})
}
