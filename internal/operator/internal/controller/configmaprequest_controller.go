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
	"errors"
	"fmt"
	"github.com/HYY-yu/sail/internal/operator/internal/config_server"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cmrv1beta1 "github.com/HYY-yu/sail/internal/operator/api/v1beta1"
)

// ConfigMapRequestReconciler reconciles a ConfigMapRequest object
type ConfigMapRequestReconciler struct {
	Namespace string
	client.Client
	Scheme *runtime.Scheme

	ConfigServer config_server.ConfigServer
}

//+kubebuilder:rbac:groups=cmr.sail.hyy-yu.space,resources=configmaprequests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cmr.sail.hyy-yu.space,resources=configmaprequests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cmr.sail.hyy-yu.space,resources=configmaprequests/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ConfigMapRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	L := log.FromContext(ctx)

	cmr := &cmrv1beta1.ConfigMapRequest{}
	err := r.Get(ctx, req.NamespacedName, cmr)
	if err != nil {
		// We'll ignore not-found errors, since they can't be fixed by an immediate requeue
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if req.Namespace != r.Namespace {
		return ctrl.Result{}, errors.New("namespace mismatch. ")
	}

	// Print the cmr
	L.Info("ConfigMapRequest", "cmr", cmr.Name)

	// cmr finalizer
	cmrFinalizerName := config_server.BaseHost + "finalizer"
	if cmr.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(cmr, cmrFinalizerName) {
			controllerutil.AddFinalizer(cmr, cmrFinalizerName)
			if err := r.Update(ctx, cmr); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(cmr, cmrFinalizerName) {
			L.V(1).Info("Finalizer configMapRequest And ConfigMap .")
			if err := r.ConfigServer.Delete(ctx, req.NamespacedName.String(), &cmr.Spec); err != nil {
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(cmr, cmrFinalizerName)
			if err := r.Update(ctx, cmr); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// getSecret if it has .
	var secretNamespaceKey string
	if cmr.Spec.NamespaceKeyInSecret != nil {
		secretNamespaceKey, err = r.getSecret(ctx, cmr)
		if err != nil {
			L.Error(err, "get Secret fail. ")
			return ctrl.Result{}, err
		}
	}

	err = r.ConfigServer.InitOrUpdate(ctx, req.NamespacedName.String(), secretNamespaceKey, &cmr.Spec)
	if err != nil {
		L.Error(err, "InitOrUpdate fail. ")
		return ctrl.Result{}, err
	}
	// cmr status
	watching, managedConfig, err := r.ConfigServer.Get(ctx, req.NamespacedName.String(), &cmr.Spec)
	if err != nil {
		L.Error(err, "get ConfigServer fail. ")
		return ctrl.Result{}, err
	}

	cmr.Status.TotalConfig = len(managedConfig)
	cmr.Status.Watching = watching
	cmr.Status.LastUpdateTime = &metav1.Time{}

	for k, v := range managedConfig {
		t := metav1.NewTime(*v.LastUpdateTime)
		mc := cmrv1beta1.ManagedConfig{
			ConfigFileName: k.String(),
			LastUpdateTime: &t,
		}

		found := false
		for i, e := range cmr.Status.ManagedConfigList {
			if e.ConfigFileName == k.String() {
				cmr.Status.ManagedConfigList[i].LastUpdateTime = mc.LastUpdateTime
				found = true
				break
			}
		}
		if !found {
			cmr.Status.ManagedConfigList = append(cmr.Status.ManagedConfigList, mc)
		}

		if cmr.Status.LastUpdateTime.Before(mc.LastUpdateTime) {
			cmr.Status.LastUpdateTime = mc.LastUpdateTime
		}
	}

	if err := r.Status().Update(ctx, cmr); err != nil {
		L.Error(err, "unable to update CMR status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ConfigMapRequestReconciler) getSecret(ctx context.Context, cmr *cmrv1beta1.ConfigMapRequest) (string, error) {
	secret := &corev1.Secret{}

	err := r.Get(ctx, client.ObjectKey{
		Namespace: cmr.Namespace,
		Name:      cmr.Spec.NamespaceKeyInSecret.Name,
	}, secret)
	if err != nil {
		return "", fmt.Errorf("get secret error: %w", err)
	}
	if data, ok := secret.Data["namespace_key"]; ok {
		return string(data), nil
	}
	return "", fmt.Errorf("you secret must has key: 'namespace_key' in secret: %s", cmr.Spec.NamespaceKeyInSecret.Name)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cmrv1beta1.ConfigMapRequest{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
