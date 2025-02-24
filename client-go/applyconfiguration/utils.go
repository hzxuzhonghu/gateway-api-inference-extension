/*
Copyright 2024 The Kubernetes Authors.

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
// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
	v1alpha1 "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha1"
	v1alpha2 "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"
	apiv1alpha1 "sigs.k8s.io/gateway-api-inference-extension/client-go/applyconfiguration/api/v1alpha1"
	apiv1alpha2 "sigs.k8s.io/gateway-api-inference-extension/client-go/applyconfiguration/api/v1alpha2"
	internal "sigs.k8s.io/gateway-api-inference-extension/client-go/applyconfiguration/internal"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=inference.networking.x-k8s.io, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("EndpointPickerConfig"):
		return &apiv1alpha1.EndpointPickerConfigApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Extension"):
		return &apiv1alpha1.ExtensionApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ExtensionConnection"):
		return &apiv1alpha1.ExtensionConnectionApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ExtensionReference"):
		return &apiv1alpha1.ExtensionReferenceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("InferenceModel"):
		return &apiv1alpha1.InferenceModelApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("InferenceModelSpec"):
		return &apiv1alpha1.InferenceModelSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("InferenceModelStatus"):
		return &apiv1alpha1.InferenceModelStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("InferencePool"):
		return &apiv1alpha1.InferencePoolApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("InferencePoolSpec"):
		return &apiv1alpha1.InferencePoolSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("InferencePoolStatus"):
		return &apiv1alpha1.InferencePoolStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PoolObjectReference"):
		return &apiv1alpha1.PoolObjectReferenceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("TargetModel"):
		return &apiv1alpha1.TargetModelApplyConfiguration{}

		// Group=inference.networking.x-k8s.io, Version=v1alpha2
	case v1alpha2.SchemeGroupVersion.WithKind("EndpointPickerConfig"):
		return &apiv1alpha2.EndpointPickerConfigApplyConfiguration{}
	case v1alpha2.SchemeGroupVersion.WithKind("Extension"):
		return &apiv1alpha2.ExtensionApplyConfiguration{}
	case v1alpha2.SchemeGroupVersion.WithKind("ExtensionConnection"):
		return &apiv1alpha2.ExtensionConnectionApplyConfiguration{}
	case v1alpha2.SchemeGroupVersion.WithKind("ExtensionReference"):
		return &apiv1alpha2.ExtensionReferenceApplyConfiguration{}
	case v1alpha2.SchemeGroupVersion.WithKind("InferenceModel"):
		return &apiv1alpha2.InferenceModelApplyConfiguration{}
	case v1alpha2.SchemeGroupVersion.WithKind("InferenceModelSpec"):
		return &apiv1alpha2.InferenceModelSpecApplyConfiguration{}
	case v1alpha2.SchemeGroupVersion.WithKind("InferenceModelStatus"):
		return &apiv1alpha2.InferenceModelStatusApplyConfiguration{}
	case v1alpha2.SchemeGroupVersion.WithKind("InferencePool"):
		return &apiv1alpha2.InferencePoolApplyConfiguration{}
	case v1alpha2.SchemeGroupVersion.WithKind("InferencePoolSpec"):
		return &apiv1alpha2.InferencePoolSpecApplyConfiguration{}
	case v1alpha2.SchemeGroupVersion.WithKind("InferencePoolStatus"):
		return &apiv1alpha2.InferencePoolStatusApplyConfiguration{}
	case v1alpha2.SchemeGroupVersion.WithKind("PoolObjectReference"):
		return &apiv1alpha2.PoolObjectReferenceApplyConfiguration{}
	case v1alpha2.SchemeGroupVersion.WithKind("TargetModel"):
		return &apiv1alpha2.TargetModelApplyConfiguration{}

	}
	return nil
}

func NewTypeConverter(scheme *runtime.Scheme) *testing.TypeConverter {
	return &testing.TypeConverter{Scheme: scheme, TypeResolver: internal.Parser()}
}
