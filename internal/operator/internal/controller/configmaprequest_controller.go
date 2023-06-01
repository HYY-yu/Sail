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
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/HYY-yu/sail/internal/operator/internal/config_server"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
//+kubebuilder:rbac:groups="",resources=secret,verbs=get;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ConfigMapRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	L := log.FromContext(ctx)

	var cmr cmrv1beta1.ConfigMapRequest
	err := r.Get(ctx, req.NamespacedName, &cmr)
	if err != nil {
		// We'll ignore not-found errors, since they can't be fixed by an immediate requeue
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if req.Namespace != r.Namespace {
		return ctrl.Result{}, errors.New("namespace mismatch. ")
	}

	// Print the cmr
	L.Info("ConfigMapRequest", "cmr", cmr.Name)

	// getSecret if it has .
	var secretNamespaceKey string
	if cmr.Spec.NamespaceKeyInSecret != nil {
		secretNamespaceKey, err = r.getSecret(ctx, &cmr)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	err = r.ConfigServer.InitOrUpdate(ctx, req.NamespacedName.String(), secretNamespaceKey, &cmr.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}
	// TODO cmr status

	// TODO cmr finalizer

	return ctrl.Result{}, nil
}

func (r *ConfigMapRequestReconciler) getSecret(ctx context.Context, cmr *cmrv1beta1.ConfigMapRequest) (string, error) {
	var secret *corev1.Secret

	err := r.Get(ctx, client.ObjectKey{
		Namespace: cmr.Namespace,
		Name:      cmr.Spec.NamespaceKeyInSecret.Name,
	}, secret)
	if err != nil {
		return "", fmt.Errorf("get secret error: %w", err)
	}
	if data, ok := secret.Data["namespace_key"]; ok {
		var buffers []byte
		_, err := base64.StdEncoding.Decode(buffers, data)
		if err != nil {
			return "", fmt.Errorf("decode base64 error: %w", err)
		}
		return string(buffers), nil
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
