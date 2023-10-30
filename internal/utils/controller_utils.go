package utils

import (
	"context"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/clock"
	"sigs.k8s.io/controller-runtime/pkg/client"

	piv1alpha1api "github.com/heliannuuthus/privateca-issuer/api/v1alpha1"
)

var realtimeClock clock.Clock = clock.RealClock{}

// GetIssuer returns either an ClusterSelfSignedIssuer or SelfSignedIssuer by its name
func GetIssuer(ctx context.Context, client client.Client, name types.NamespacedName) (piv1alpha1api.GenericIssuer, error) {
	iss := new(piv1alpha1api.SelfSignedIssuer)
	err := client.Get(ctx, name, iss)
	if err != nil {
		ciss := new(piv1alpha1api.ClusterSelfSignedIssuer)
		cname := types.NamespacedName{
			Name: name.Name,
		}
		err = client.Get(ctx, cname, ciss)
		if err != nil {
			return nil, err
		}
		return ciss, nil
	}
	return iss, nil
}

// SetIssuerCondition sets the ready state of an issuer and updates it in the cluster
func SetIssuerCondition(log logr.Logger, issuer piv1alpha1api.GenericIssuer, conditionType string,
	status metav1.ConditionStatus, reason, message string) {
	newCondition := metav1.Condition{
		Type:    conditionType,
		Status:  status,
		Reason:  reason,
		Message: message,
	}

	now := metav1.NewTime(realtimeClock.Now())
	newCondition.LastTransitionTime = now

	for idx, cond := range issuer.GetStatus().Conditions {
		if cond.Type != conditionType {
			continue
		}

		if cond.Status == status {
			newCondition.LastTransitionTime = cond.LastTransitionTime
		}

		issuer.GetStatus().Conditions[idx] = newCondition
		return
	}

	issuer.GetStatus().Conditions = append(issuer.GetStatus().Conditions, newCondition)
}
