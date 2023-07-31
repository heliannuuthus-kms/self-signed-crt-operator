package secret

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Manager interface {
	GetCert() ([]byte, error)
	GetPriKey() ([]byte, error)
}

type SecretsManager struct {
	secret v1.Secret
}

func NewSecretsManager(ctx context.Context, client client.Client, secretName types.NamespacedName) (*SecretsManager, error) {
	var secret v1.Secret
	if err := client.Get(ctx, secretName, &secret); err != nil {
		return nil, err
	}
	return &SecretsManager{
		secret: secret,
	}, nil
}

func (m *SecretsManager) GetCert() ([]byte, error) {
	return m.secret.Data["rootCert"], nil
}

func (m *SecretsManager) GetPriKey() ([]byte, error) {
	return m.secret.Data["priKey"], nil
}

type PKIManager struct {
}

func NewPKIManager() *PKIManager {
	return &PKIManager{}
}

func (m *PKIManager) GetCert() ([]byte, error) {
	return nil, nil
}

func (m *PKIManager) GetPriKey() ([]byte, error) {
	return nil, nil
}
