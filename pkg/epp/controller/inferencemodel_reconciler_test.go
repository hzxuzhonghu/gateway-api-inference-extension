/*
Copyright 2025 The Kubernetes Authors.

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
	"sync"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"
	"sigs.k8s.io/gateway-api-inference-extension/pkg/epp/datastore"
	logutil "sigs.k8s.io/gateway-api-inference-extension/pkg/epp/util/logging"
)

var (
	infModel1 = &v1alpha2.InferenceModel{
		Spec: v1alpha2.InferenceModelSpec{
			ModelName: "fake model1",
			PoolRef:   v1alpha2.PoolObjectReference{Name: "test-pool"},
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-service",
		},
	}
	infModel1Modified = &v1alpha2.InferenceModel{
		Spec: v1alpha2.InferenceModelSpec{
			ModelName: "fake model1",
			PoolRef:   v1alpha2.PoolObjectReference{Name: "test-poolio"},
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-service",
		},
	}
	infModel2 = &v1alpha2.InferenceModel{
		Spec: v1alpha2.InferenceModelSpec{
			ModelName: "fake model",
			PoolRef:   v1alpha2.PoolObjectReference{Name: "test-pool"},
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-service-2",
		},
	}
)

func TestUpdateDatastore_InferenceModelReconciler(t *testing.T) {
	logger := logutil.NewTestLogger()

	tests := []struct {
		name                string
		datastore           datastore.Datastore
		incomingService     *v1alpha2.InferenceModel
		wantInferenceModels *sync.Map
	}{
		{
			name: "No Services registered; valid, new service incoming.",
			datastore: datastore.NewFakeDatastore(nil, nil, &v1alpha2.InferencePool{
				Spec: v1alpha2.InferencePoolSpec{
					Selector: map[v1alpha2.LabelKey]v1alpha2.LabelValue{"app": "vllm"},
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "test-pool",
					ResourceVersion: "Old and boring",
				},
			}),

			incomingService:     infModel1,
			wantInferenceModels: populateServiceMap(infModel1),
		},
		{
			name: "Removing existing service.",
			datastore: datastore.NewFakeDatastore(nil, populateServiceMap(infModel1), &v1alpha2.InferencePool{
				Spec: v1alpha2.InferencePoolSpec{
					Selector: map[v1alpha2.LabelKey]v1alpha2.LabelValue{"app": "vllm"},
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "test-pool",
					ResourceVersion: "Old and boring",
				},
			}),
			incomingService:     infModel1Modified,
			wantInferenceModels: populateServiceMap(),
		},
		{
			name: "Unrelated service, do nothing.",
			datastore: datastore.NewFakeDatastore(nil, populateServiceMap(infModel1), &v1alpha2.InferencePool{
				Spec: v1alpha2.InferencePoolSpec{
					Selector: map[v1alpha2.LabelKey]v1alpha2.LabelValue{"app": "vllm"},
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "test-pool",
					ResourceVersion: "Old and boring",
				},
			}),
			incomingService: &v1alpha2.InferenceModel{
				Spec: v1alpha2.InferenceModelSpec{
					ModelName: "fake model",
					PoolRef:   v1alpha2.PoolObjectReference{Name: "test-poolio"},
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "unrelated-service",
				},
			},
			wantInferenceModels: populateServiceMap(infModel1),
		},
		{
			name: "Add to existing",
			datastore: datastore.NewFakeDatastore(nil, populateServiceMap(infModel1), &v1alpha2.InferencePool{
				Spec: v1alpha2.InferencePoolSpec{
					Selector: map[v1alpha2.LabelKey]v1alpha2.LabelValue{"app": "vllm"},
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "test-pool",
					ResourceVersion: "Old and boring",
				},
			}),
			incomingService:     infModel2,
			wantInferenceModels: populateServiceMap(infModel1, infModel2),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool, err := test.datastore.PoolGet()
			if err != nil {
				t.Fatalf("failed to get pool: %v", err)
			}
			reconciler := &InferenceModelReconciler{
				Datastore:          test.datastore,
				PoolNamespacedName: types.NamespacedName{Name: pool.Name},
			}
			reconciler.updateDatastore(logger, test.incomingService)

			test.wantInferenceModels.Range(func(k, v any) bool {
				_, exist := test.datastore.ModelGet(k.(string))
				if !exist {
					t.Fatalf("failed to get model %s", k)
				}
				return true
			})
		})
	}
}

func TestReconcile_ResourceNotFound(t *testing.T) {
	// Set up the scheme.
	scheme := runtime.NewScheme()
	_ = v1alpha2.AddToScheme(scheme)

	// Create a fake client with no InferenceModel objects.
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

	// Create a minimal datastore.
	datastore := datastore.NewFakeDatastore(nil, nil, &v1alpha2.InferencePool{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pool"},
	})

	// Create the reconciler.
	reconciler := &InferenceModelReconciler{
		Client:             fakeClient,
		Scheme:             scheme,
		Record:             record.NewFakeRecorder(10),
		Datastore:          datastore,
		PoolNamespacedName: types.NamespacedName{Name: "test-pool"},
	}

	// Create a request for a non-existent resource.
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: "non-existent-model", Namespace: "default"}}

	// Call Reconcile.
	result, err := reconciler.Reconcile(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error when resource is not found, got %v", err)
	}

	// Check that no requeue is requested.
	if result.Requeue || result.RequeueAfter != 0 {
		t.Errorf("expected no requeue, got %+v", result)
	}
}

func TestReconcile_ModelMarkedForDeletion(t *testing.T) {
	// Set up the scheme.
	scheme := runtime.NewScheme()
	_ = v1alpha2.AddToScheme(scheme)

	// Create an InferenceModel object.
	now := metav1.Now()
	existingModel := &v1alpha2.InferenceModel{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "existing-model",
			Namespace:         "default",
			DeletionTimestamp: &now,
			Finalizers:        []string{"finalizer"},
		},
		Spec: v1alpha2.InferenceModelSpec{
			ModelName: "fake-model",
			PoolRef:   v1alpha2.PoolObjectReference{Name: "test-pool"},
		},
	}

	// Create a fake client with the existing model.
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existingModel).Build()

	// Create a minimal datastore.
	datastore := datastore.NewFakeDatastore(nil, nil, &v1alpha2.InferencePool{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pool"},
	})

	// Create the reconciler.
	reconciler := &InferenceModelReconciler{
		Client:             fakeClient,
		Scheme:             scheme,
		Record:             record.NewFakeRecorder(10),
		Datastore:          datastore,
		PoolNamespacedName: types.NamespacedName{Name: "test-pool", Namespace: "default"},
	}

	// Create a request for the existing resource.
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: "existing-model", Namespace: "default"}}

	// Call Reconcile.
	result, err := reconciler.Reconcile(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error when resource exists, got %v", err)
	}

	// Check that no requeue is requested.
	if result.Requeue || result.RequeueAfter != 0 {
		t.Errorf("expected no requeue, got %+v", result)
	}

	// Verify that the datastore was not updated.
	if _, exist := datastore.ModelGet(existingModel.Spec.ModelName); exist {
		t.Errorf("expected datastore to not contain model %q", existingModel.Spec.ModelName)
	}
}

func TestReconcile_ResourceExists(t *testing.T) {
	// Set up the scheme.
	scheme := runtime.NewScheme()
	_ = v1alpha2.AddToScheme(scheme)

	// Create an InferenceModel object.
	existingModel := &v1alpha2.InferenceModel{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "existing-model",
			Namespace: "default",
		},
		Spec: v1alpha2.InferenceModelSpec{
			ModelName: "fake-model",
			PoolRef:   v1alpha2.PoolObjectReference{Name: "test-pool"},
		},
	}

	// Create a fake client with the existing model.
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existingModel).Build()

	// Create a minimal datastore.
	datastore := datastore.NewFakeDatastore(nil, nil, &v1alpha2.InferencePool{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pool"},
	})

	// Create the reconciler.
	reconciler := &InferenceModelReconciler{
		Client:             fakeClient,
		Scheme:             scheme,
		Record:             record.NewFakeRecorder(10),
		Datastore:          datastore,
		PoolNamespacedName: types.NamespacedName{Name: "test-pool", Namespace: "default"},
	}

	// Create a request for the existing resource.
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: "existing-model", Namespace: "default"}}

	// Call Reconcile.
	result, err := reconciler.Reconcile(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error when resource exists, got %v", err)
	}

	// Check that no requeue is requested.
	if result.Requeue || result.RequeueAfter != 0 {
		t.Errorf("expected no requeue, got %+v", result)
	}

	// Verify that the datastore was updated.
	if _, exist := datastore.ModelGet(existingModel.Spec.ModelName); !exist {
		t.Errorf("expected datastore to contain model %q", existingModel.Spec.ModelName)
	}
}

func populateServiceMap(services ...*v1alpha2.InferenceModel) *sync.Map {
	returnVal := &sync.Map{}

	for _, service := range services {
		returnVal.Store(service.Spec.ModelName, service)
	}
	return returnVal
}
