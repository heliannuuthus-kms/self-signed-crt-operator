package controllers

import (
	"context"
	piv1alpha1api "github.com/heliannuuthus/privateca-issuer/api/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ClusterIssuerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=privateca-issuer.gitlabee.chehejia.com,resources=clusterissuers,verbs=get;list;watch
// +kubebuilder:rbac:groups=privateca-issuer.gitlabee.chehejia.com,resources=clusterissuers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=privateca-issuer.gitlabee.chehejia.com,resources=clusterissuers/finalizers,verbs=update

func (r *ClusterIssuerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

func (r *ClusterIssuerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&piv1alpha1api.ClusterIssuer{}).
		Complete(r)
}
