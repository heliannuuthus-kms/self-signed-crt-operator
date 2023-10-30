package signer

import (
	"context"
	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/go-logr/logr"
	piv1alpha1api "github.com/heliannuuthus/privateca-issuer/api/v1alpha1"
)

type HealthChecker interface {
	Check() error
}

type HealthCheckerBuilder func(piv1alpha1api.GenericIssuer, map[string][]byte) (HealthChecker, error)

type Signer interface {
	Sign(ctx context.Context, log logr.Logger, cr *cmapi.CertificateRequest) ([]byte, []byte, error)
}
