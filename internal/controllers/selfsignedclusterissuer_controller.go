package controllers

import (
	"context"
	piv1alpha1api "github.com/heliannuuthus/privateca-issuer/api/v1alpha1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterSelfSignedIssuerReconciler reconciles a ClusterSelfSignedIssuer object
type ClusterSelfSignedIssuerReconciler struct {
	client.Client
	Log               logr.Logger
	Scheme            *runtime.Scheme
	GenericController *GenericIssuerReconciler
}

// +kubebuilder:rbac:groups=privateca-issuer.github.com,resources=clustercaissuers,verbs=get;list;watch
// +kubebuilder:rbac:groups=privateca-issuer.github.com,resources=clustercaissuers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=privateca-issuer.github.com,resources=clustercaissuers/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *ClusterSelfSignedIssuerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("clustercaissuer", req.NamespacedName)
	iss := new(piv1alpha1api.ClusterSelfSignedIssuer)
	if err := r.Client.Get(ctx, req.NamespacedName, iss); err != nil {
		log.Error(err, "Failed to request ClusterSelfSignedIssuer")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return r.GenericController.Reconcile(ctx, req, iss)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterSelfSignedIssuerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&piv1alpha1api.ClusterSelfSignedIssuer{}).
		Complete(r)
}
