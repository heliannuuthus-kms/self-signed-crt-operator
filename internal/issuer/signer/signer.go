package signer

import (
	piv1alpha1api "github.com/heliannuuthus/privateca-issuer/api/v1alpha1"
)

type HealthChecker interface {
	Check() error
}

type HealthCheckerBuilder func(*piv1alpha1api.IssuerSpec, map[string][]byte) (HealthChecker, error)

type Signer interface {
	Sign([]byte) ([]byte, error)
}

type Builder func(*piv1alpha1api.IssuerSpec, map[string][]byte) (Signer, error)
