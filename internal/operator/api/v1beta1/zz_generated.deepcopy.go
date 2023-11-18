//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigMapRequest) DeepCopyInto(out *ConfigMapRequest) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigMapRequest.
func (in *ConfigMapRequest) DeepCopy() *ConfigMapRequest {
	if in == nil {
		return nil
	}
	out := new(ConfigMapRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ConfigMapRequest) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigMapRequestList) DeepCopyInto(out *ConfigMapRequestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ConfigMapRequest, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigMapRequestList.
func (in *ConfigMapRequestList) DeepCopy() *ConfigMapRequestList {
	if in == nil {
		return nil
	}
	out := new(ConfigMapRequestList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ConfigMapRequestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigMapRequestSpec) DeepCopyInto(out *ConfigMapRequestSpec) {
	*out = *in
	if in.NamespaceKeyInSecret != nil {
		in, out := &in.NamespaceKeyInSecret, &out.NamespaceKeyInSecret
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
	if in.Configs != nil {
		in, out := &in.Configs, &out.Configs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Merge != nil {
		in, out := &in.Merge, &out.Merge
		*out = new(bool)
		**out = **in
	}
	if in.MergeConfigFile != nil {
		in, out := &in.MergeConfigFile, &out.MergeConfigFile
		*out = new(string)
		**out = **in
	}
	if in.Watch != nil {
		in, out := &in.Watch, &out.Watch
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigMapRequestSpec.
func (in *ConfigMapRequestSpec) DeepCopy() *ConfigMapRequestSpec {
	if in == nil {
		return nil
	}
	out := new(ConfigMapRequestSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigMapRequestStatus) DeepCopyInto(out *ConfigMapRequestStatus) {
	*out = *in
	if in.LastUpdateTime != nil {
		in, out := &in.LastUpdateTime, &out.LastUpdateTime
		*out = (*in).DeepCopy()
	}
	if in.ManagedConfigList != nil {
		in, out := &in.ManagedConfigList, &out.ManagedConfigList
		*out = make([]ManagedConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigMapRequestStatus.
func (in *ConfigMapRequestStatus) DeepCopy() *ConfigMapRequestStatus {
	if in == nil {
		return nil
	}
	out := new(ConfigMapRequestStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedConfig) DeepCopyInto(out *ManagedConfig) {
	*out = *in
	if in.LastUpdateTime != nil {
		in, out := &in.LastUpdateTime, &out.LastUpdateTime
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedConfig.
func (in *ManagedConfig) DeepCopy() *ManagedConfig {
	if in == nil {
		return nil
	}
	out := new(ManagedConfig)
	in.DeepCopyInto(out)
	return out
}