/*
Copyright 2021 The Kubernetes Authors.

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

package controllers

import (
	"context"
	piv1alpha1api "github.com/heliannuuthus/privateca-issuer/api/v1alpha1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CAIssuerReconciler reconciles a CAIssuer object
type CAIssuerReconciler struct {
	client.Client
	Log               logr.Logger
	Scheme            *runtime.Scheme
	GenericController *GenericIssuerReconciler
}

// +kubebuilder:rbac:groups=privateca-issuer.github.com,resources=caissuers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=privateca-issuer.github.com,resources=caissuers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=privateca-issuer.github.com,resources=caissuers/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *CAIssuerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("caissuer", req.NamespacedName)
	iss := new(piv1alpha1api.CAIssuer)
	if err := r.Client.Get(ctx, req.NamespacedName, iss); err != nil {
		log.Error(err, "Failed to request CAIssuer")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return r.GenericController.Reconcile(ctx, req, iss)
}

// SetupWithManager sets up the controller with the Manager.
func (r *CAIssuerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&piv1alpha1api.CAIssuer{}).
		Complete(r)
}
