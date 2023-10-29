package signer

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/go-logr/logr"
	"github.com/heliannuuthus/privateca-issuer/internal/issuer/secrets"
	capi "k8s.io/api/certificates/v1beta1"
	"math/big"
	"time"
)

type webCaSigner struct {
	secretManager secrets.Manager
}

func (o *webCaSigner) Check() error {
	return nil
}

func NewCASinger(secretManager secrets.Manager) (Signer, error) {
	return &webCaSigner{secretManager: secretManager}, nil
}

func (o *webCaSigner) Sign(ctx context.Context, cr *cmapi.CertificateRequest, log logr.Logger) ([]byte, []byte, error) {
	csr, err := toX509CSR(cr.Spec.Request)
	if err != nil {
		log.V(4).Error(err, "csr is invalid format")
		return nil, nil, err
	}
	priKey, err := o.secretManager.GetPriKey()
	if err != nil {
		log.V(4).Error(err, "secrets get priKey failed")
		return nil, nil, err
	}
	key, err := toRsaPriKey(priKey)
	if err != nil {
		log.V(4).Error(err, "priKey must be rsa")
		return nil, nil, err
	}
	certPem, err := o.secretManager.GetCert()
	if err != nil {
		log.V(4).Error(err, "load root cert failed")
		return nil, nil, err
	}
	cert, err := toX509Cert(certPem)
	if err != nil {
		log.V(4).Error(err, "root cert be format to X509 failed")
		return nil, nil, err
	}
	ca := &CertificateAuthority{
		Certificate: cert,
		PrivateKey:  key,
		Backdate:    5 * time.Minute,
	}
	crtDER, err := ca.Sign(csr.Raw, PermissiveSigningPolicy{
		TTL: duration,
		Usages: []capi.KeyUsage{
			capi.UsageServerAuth,
		},
	})
	if err != nil {
		return nil, nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: crtDER,
	}), nil, nil
}

func toRsaPriKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	// extract PEM from request object
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("PEM block type must be RSA PRIVATE KEY")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func toX509Cert(pemBytes []byte) (*x509.Certificate, error) {
	// extract PEM from request object
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("PEM block type must be CERTIFICATE")
	}
	return x509.ParseCertificate(block.Bytes)
}

func toX509CSR(pemBytes []byte) (*x509.CertificateRequest, error) {
	// extract PEM from request object
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, errors.New("PEM block type must be CERTIFICATE REQUEST")
	}
	return x509.ParseCertificateRequest(block.Bytes)
}

var (
	duration = time.Hour * 24 * 365
)

var serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)

// CertificateAuthority implements a certificate authority that supports policy
// based signing. It's used by the signing controllers.
type CertificateAuthority struct {
	// RawCert is an optional field to determine if signing cert/key pairs have changed
	RawCert []byte
	// RawKey is an optional field to determine if signing cert/key pairs have changed
	RawKey []byte

	Certificate *x509.Certificate
	PrivateKey  crypto.Signer
	Backdate    time.Duration
	Now         func() time.Time
}

// Sign signs a certificate request, applying a SigningPolicy and returns a DER
// encoded x509 certificate.
func (ca *CertificateAuthority) Sign(crDER []byte, policy SigningPolicy) ([]byte, error) {
	now := time.Now()
	if ca.Now != nil {
		now = ca.Now()
	}

	nbf := now.Add(-ca.Backdate)
	if !nbf.Before(ca.Certificate.NotAfter) {
		return nil, fmt.Errorf("the signer has expired: NotAfter=%v", ca.Certificate.NotAfter)
	}

	cr, err := x509.ParseCertificateRequest(crDER)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate request: %v", err)
	}
	if err := cr.CheckSignature(); err != nil {
		return nil, fmt.Errorf("unable to verify certificate request signature: %v", err)
	}

	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("unable to generate a serial number for %s: %v", cr.Subject.CommonName, err)
	}

	tmpl := &x509.Certificate{
		SerialNumber:       serialNumber,
		Subject:            cr.Subject,
		DNSNames:           cr.DNSNames,
		IPAddresses:        cr.IPAddresses,
		EmailAddresses:     cr.EmailAddresses,
		URIs:               cr.URIs,
		PublicKeyAlgorithm: cr.PublicKeyAlgorithm,
		PublicKey:          cr.PublicKey,
		Extensions:         cr.Extensions,
		ExtraExtensions:    cr.ExtraExtensions,
		NotBefore:          nbf,
	}
	if err := policy.apply(tmpl); err != nil {
		return nil, err
	}

	if tmpl.NotAfter.After(ca.Certificate.NotAfter) {
		tmpl.NotAfter = ca.Certificate.NotAfter
	}
	if now.After(ca.Certificate.NotAfter) {
		return nil, fmt.Errorf("refusing to sign a certificate that expired in the past")
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, ca.Certificate, cr.PublicKey, ca.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign certificate: %v", err)
	}
	return der, nil
}
