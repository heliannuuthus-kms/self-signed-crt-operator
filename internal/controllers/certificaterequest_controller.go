package controllers

import (
	"context"
	"fmt"
	"github.com/heliannuuthus/privateca-issuer/internal/utils"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"

	"github.com/go-logr/logr"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/clock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cmutil "github.com/cert-manager/cert-manager/pkg/api/util"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	piv1alpha1api "github.com/heliannuuthus/privateca-issuer/api/v1alpha1"
)

// CertificateRequestReconciler reconciles a SelfSignedIssuer object
type CertificateRequestReconciler struct {
	client.Client
	Log                    logr.Logger
	Scheme                 *runtime.Scheme
	Recorder               record.EventRecorder
	Clock                  clock.Clock
	CheckApprovedCondition bool
}

// +kubebuilder:rbac:groups=cert-manager.io,resources=certificaterequests,verbs=get;list;watch;update
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificaterequests/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *CertificateRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("certificaterequest", req.NamespacedName)
	cr := new(cmapi.CertificateRequest)
	if err := r.Client.Get(ctx, req.NamespacedName, cr); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		log.Error(err, "Failed to request CertificateRequest")
		return ctrl.Result{}, err
	}

	if cr.Spec.IssuerRef.Group != piv1alpha1api.GroupVersion.Group {
		log.V(4).Info("CertificateRequest does not specify an issuerRef matching our group")
		return ctrl.Result{}, nil
	}

	// Ignore CertificateRequest if it is already Ready
	if cmutil.CertificateRequestHasCondition(cr, cmapi.CertificateRequestCondition{
		Type:   cmapi.CertificateRequestConditionReady,
		Status: cmmeta.ConditionTrue,
	}) {
		log.V(4).Info("CertificateRequest is Ready. Ignoring.")
		return ctrl.Result{}, nil
	}
	// Ignore CertificateRequest if it is already Failed
	if cmutil.CertificateRequestHasCondition(cr, cmapi.CertificateRequestCondition{
		Type:   cmapi.CertificateRequestConditionReady,
		Status: cmmeta.ConditionFalse,
		Reason: cmapi.CertificateRequestReasonFailed,
	}) {
		log.V(4).Info("CertificateRequest is Failed. Ignoring.")
		return ctrl.Result{}, nil
	}
	// Ignore CertificateRequest if it already has a Denied Ready Reason
	if cmutil.CertificateRequestHasCondition(cr, cmapi.CertificateRequestCondition{
		Type:   cmapi.CertificateRequestConditionReady,
		Status: cmmeta.ConditionFalse,
		Reason: cmapi.CertificateRequestReasonDenied,
	}) {
		log.V(4).Info("CertificateRequest already has a Ready condition with Denied Reason. Ignoring.")
		return ctrl.Result{}, nil
	}

	// If CertificateRequest has been denied, mark the CertificateRequest as
	// Ready=Denied and set FailureTime if not already.
	if cmutil.CertificateRequestIsDenied(cr) {
		log.V(4).Info("CertificateRequest has been denied. Marking as failed.")

		if cr.Status.FailureTime == nil {
			nowTime := metav1.NewTime(r.Clock.Now())
			cr.Status.FailureTime = &nowTime
		}

		message := "The CertificateRequest was denied by an approval controllers"
		return ctrl.Result{}, r.setStatus(ctx, cr, cmmeta.ConditionFalse, cmapi.CertificateRequestReasonDenied, message)
	}

	if r.CheckApprovedCondition {
		// If CertificateRequest has not been approved, exit early.
		if !cmutil.CertificateRequestIsApproved(cr) {
			log.V(4).Info("certificate request has not been approved")
			return ctrl.Result{}, nil
		}
	}

	if len(cr.Status.Certificate) > 0 {
		log.V(4).Info("Certificate was already signed")
		return ctrl.Result{}, nil
	}

	issuerName := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.IssuerRef.Name,
	}

	if cr.Spec.IssuerRef.Kind == "ClusterSelfSignedIssuer" {
		issuerName.Namespace = ""
	}
	// getOrDefault(issuer, clusterIssuer)
	iss, err := utils.GetIssuer(ctx, r.Client, issuerName)

	if err != nil {
		log.Error(err, "failed to retrieve SelfSignedIssuer resource")
		_ = r.setStatus(ctx, cr, cmmeta.ConditionFalse, cmapi.CertificateRequestReasonFailed, "issuer could not be found")
		return ctrl.Result{}, err
	}

	if !isReady(iss) {
		err := fmt.Errorf("issuer %s is not ready", iss.GetName())
		_ = r.setStatus(ctx, cr, cmmeta.ConditionFalse, cmapi.CertificateRequestReasonFailed, "issuer is not ready")
		return ctrl.Result{}, err
	}
	//if iss.GetSpec().CASecretName != "" {
	//	if r.secretManager, err = secrets.NewSecretsManager(ctx, r.Client, types.NamespacedName{Name: iss.GetSpec().CASecretName, Namespace: cr.GetNamespace()}); err != nil {
	//		return ctrl.Result{}, fmt.Errorf("failed to get Secret containing Issuer credentials, secrets name: %s, reason: %v", iss.GetSpec().CASecretName, err)
	//	}
	//} else {
	//	r.secretManager = secrets.NewPKIManager()
	//}
	//
	//var caSigner signer.Signer
	//if caSigner, err = r.Builder(r.secretManager); err != nil {
	//	return ctrl.Result{}, r.setStatus(ctx, cr, cmmeta.ConditionFalse, cmapi.CertificateRequestReasonFailed, "failed to request certificate from PCA: "+err.Error())
	//}
	//pem, rootPem, err := caSigner.Sign(ctx, cr, log)
	//if err != nil {
	//	log.Error(err, "failed to request certificate from PCA")
	//	return ctrl.Result{}, r.setStatus(ctx, cr, cmmeta.ConditionFalse, cmapi.CertificateRequestReasonFailed, "failed to request certificate from PCA: "+err.Error())
	//}
	//cr.Status.Certificate = pem
	//cr.Status.CA = rootPem

	return ctrl.Result{}, r.setStatus(ctx, cr, cmmeta.ConditionTrue, cmapi.CertificateRequestReasonIssued, "certificate issued")
}

// SetupWithManager sets up the controllers with the Manager.
func (r *CertificateRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cmapi.CertificateRequest{}).
		Complete(r)
}

func isReady(issuer piv1alpha1api.GenericIssuer) bool {
	for _, condition := range issuer.GetStatus().Conditions {
		if condition.Type == piv1alpha1api.ConditionTypeReady && condition.Status == metav1.ConditionTrue {
			return true
		}
	}
	return false
}

func (r *CertificateRequestReconciler) setStatus(ctx context.Context, cr *cmapi.CertificateRequest, status cmmeta.ConditionStatus, reason, message string, args ...interface{}) error {
	completeMessage := fmt.Sprintf(message, args...)
	cmutil.SetCertificateRequestCondition(cr, "Ready", status, reason, completeMessage)

	eventType := core.EventTypeNormal
	if status == cmmeta.ConditionFalse {
		eventType = core.EventTypeWarning
	}
	r.Recorder.Event(cr, eventType, reason, completeMessage)

	return r.Client.Status().Update(ctx, cr)
}
