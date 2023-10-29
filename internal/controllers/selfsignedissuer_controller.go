package controllers

import (
	"context"
	piv1alpha1api "github.com/heliannuuthus/privateca-issuer/api/v1alpha1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SelfSignedIssuerReconciler reconciles a SelfSignedIssuer object
type SelfSignedIssuerReconciler struct {
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

func (r *SelfSignedIssuerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("caissuer", req.NamespacedName)
	iss := new(piv1alpha1api.SelfSignedIssuer)
	if err := r.Client.Get(ctx, req.NamespacedName, iss); err != nil {
		log.Error(err, "Failed to request SelfSignedIssuer")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return r.GenericController.Reconcile(ctx, req, iss)
}

// SetupWithManager sets up the controller with the Manager.
func (r *SelfSignedIssuerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&piv1alpha1api.SelfSignedIssuer{}).
		Complete(r)
}
