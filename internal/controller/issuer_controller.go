package controller

import (
	"context"
	"errors"
	piv1alpha1api "github.com/heliannuuthus/privateca-issuer/api/v1alpha1"
	"github.com/heliannuuthus/privateca-issuer/internal/issuer/signer"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	defaultHealthCheckInterval = time.Minute
)

var (
	errGetAuthSecret        = errors.New("failed to get Secret containing Issuer credentials")
	errHealthCheckerBuilder = errors.New("failed to build the healthchecker")
	errHealthCheckerCheck   = errors.New("healthcheck failed")
)

type IssuerReconciler struct {
	client.Client
	Kind                     string
	Scheme                   *runtime.Scheme
	ClusterResourceNamespace string
	HealthCheckerBuilder     signer.HealthCheckerBuilder
}

//+kubebuilder:rbac:groups=privateca-issuer.gitlabee.chehejia.com,resources=issuers,verbs=get;list;watch
//+kubebuilder:rbac:groups=privateca-issuer.gitlabee.chehejia.com,resources=issuers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=privateca-issuer.gitlabee.chehejia.com,resources=issuers/finalizers,verbs=update

func (r *IssuerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

func (r *IssuerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&piv1alpha1api.Issuer{}).
		Complete(r)
}
