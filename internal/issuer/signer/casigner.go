package signer

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	piv1alpha1api "github.com/heliannuuthus/privateca-issuer/api/v1alpha1"
	capi "k8s.io/api/certificates/v1beta1"
	"time"
)

type webCaSigner struct {
}

func (o *webCaSigner) Check() error {
	return nil
}

func ExampleHealthCheckerFromIssuerAndSecretData(*piv1alpha1api.IssuerSpec, map[string][]byte) (HealthChecker, error) {
	return &webCaSigner{}, nil
}

func ExampleSignerFromIssuerAndSecretData(data *piv1alpha1api.IssuerSpec, secret map[string][]byte) (Signer, error) {
	return &webCaSigner{}, nil
}

func (o *webCaSigner) Sign(csrBytes []byte) ([]byte, error) {
	csr, err := parseCSR(csrBytes)
	if err != nil {
		return nil, err
	}
	key, err := parseKey(keyPEM)
	if err != nil {
		return nil, err
	}
	cert, err := parseCert(certPEM)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: crtDER,
	}), nil
}

func parseKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	// extract PEM from request object
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("PEM block type must be RSA PRIVATE KEY")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func parseCert(pemBytes []byte) (*x509.Certificate, error) {
	// extract PEM from request object
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("PEM block type must be CERTIFICATE")
	}
	return x509.ParseCertificate(block.Bytes)
}

func parseCSR(pemBytes []byte) (*x509.CertificateRequest, error) {
	// extract PEM from request object
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, errors.New("PEM block type must be CERTIFICATE REQUEST")
	}
	return x509.ParseCertificateRequest(block.Bytes)
}

var (
	keyPEM = []byte(`
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDQFkmkG0qwlPFw
9o8vjaBOisor0VUqccc2GbSCwVjTmhJ6SBP4z4nZUh4m1uEvuFK4SE9SSbBONIbR
02W6oLlYCjIC2Uw4TMFSqR5kzSJITC5tJxb32ujbIOCXTbRtWvRgC2f45OeKiSRZ
oVVU21B6OOvvT+3I6nJJu/FBZpUxv/DExvqoPPHXWoa+yPTq7vHwKpsAj/Kmpuni
qhtm8EbKapIa3uyJKBdXmXClhvJAge5/GEA58welwLLkuI2ic9CIgWSxV8WmVQBt
7IZpMs69kcddKFQf1YCqFMMBhzJxtSg72Hs4cP+ykHjNNMUPTA+LeM9B+5gMAzPm
ZyEj+v6fAgMBAAECggEBAJwhLO3yAE+P4bylcvf2JuLnphvMfD9VkWhZTySQl+pk
/xo6/KlCZyblQ3RW5C1e+soEj2epnJyBMus612h5cbfKJo4WpubTSHaSKBjwBZoD
dw41NzmPSgool/tOtWMbzKJHzKJmdghvMBQERjdeOvsJvJUZ/ssyhcAnQTSWGLly
9W0+CuLbhdFsVc53KKmmNJ/74rcBo+dEFTgxpLmSeEtYfz3TtVXZ4RmVfXpoo24w
Y8oFV/WlR7Qr8I/T7VY35XqnbsfmfgT4WRWgrTcbXxZb+7Ce+LKHcze9GHXWO+me
m3FYQNu3Ch3fZ3n8iWIGHBWiw0aEvBR7pJRTYCbbXekCgYEA7DOZECPrxmvYCTai
3Wb3VfV8h+DJA0gmnM1dMFHnUSgxi3UXW6m3lfJjEC3kq2JAbN8dioYo7RCxP+MQ
qr4Lza2nBUA1k5cE6L/YHyiLyTKwnnTk0hlFm68yC5nWc850SDY2T1nuDer4M/0x
8Spu/VfIJtjZ3Ke+JybnNAI8YdMCgYEA4YdpBowNHKEzFJPk2nrmHDG8tJeYBfQn
8268pvgyhUQXmYVIhTf33qe6UJ7V1XGO99kq8ck7cGMSuauR+7wpb2+TjpuH3GFm
n885GMgwv42PcHJM3JP9j/1DXzwXQRfYw1OAuexZkTJGkp2lpu0I2pocU2tiziSY
YjeG27yrpIUCgYEAn1+BW76hC9Ugg7b11WXwZXOqfxRRDYHVa9+1jTD2X3A7Xdm3
1QWC9g4CgZw1ut4kklFJYXp8itjEgFL5n/tzg2g0VfqpK9iuW012yi9VgoBNY92D
t6+NpCpmHiXC6YjYNRE/O/N2CLYOmyWwWQVEtnRQfMW82oHkcA5z2kfX7jkCgYBl
ox/Kyo0SLPeXO3t0ltRjOmr/vB3P+ROUGoC8grhJ5MD5994R44I6fr5xnNNjaNT0
j5NR+c1mvc9vi4mzuD24McF/EEqvH9ofBUWHDJkjioltNKW89pjcLlgRcEROmo+e
n2Aw6foHfG/fnVpNGx/VXISNd6TEoCtof/uvxZxY/QKBgF33fvoAoGmsYLSzCG26
wydPT4pnau3g5+qcfCCMcDWr0W4RFp4ooSv/iZbz0QPvaYTOAMbJlRRaXs1SRpv2
PysXqvUQXhX2gOglfLBakhxciDYyZEiTAoSu+DmGgMRcHfPTl3X2fhZoIXd9aKj8
lgp+R9ojTy8o+cqFKZJTC6Yz
-----END PRIVATE KEY-----
`)
	certPEM = []byte(`
-----BEGIN CERTIFICATE-----
MIIDujCCAqKgAwIBAgIEBbXq6TANBgkqhkiG9w0BAQsFADB2MQswCQYDVQQGEwJD
TjELMAkGA1UEAwwCY24xCzAJBgNVBAgMAmNuMQswCQYDVQQHDAJjbjELMAkGA1UE
CgwCY20xCzAJBgNVBAsMAmNuMSYwJAYJKoZIhvcNAQkBFhdoZWxpYW5udXV0aHVz
QGdtYWlsLmNvbTAeFw0yMzA3MTkwOTQyMjlaFw0yNTA3MTgwOTQyMjlaMHYxCzAJ
BgNVBAYTAkNOMQswCQYDVQQDDAJjbjELMAkGA1UECAwCY24xCzAJBgNVBAcMAmNu
MQswCQYDVQQKDAJjbTELMAkGA1UECwwCY24xJjAkBgkqhkiG9w0BCQEWF2hlbGlh
bm51dXRodXNAZ21haWwuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEA0BZJpBtKsJTxcPaPL42gTorKK9FVKnHHNhm0gsFY05oSekgT+M+J2VIeJtbh
L7hSuEhPUkmwTjSG0dNluqC5WAoyAtlMOEzBUqkeZM0iSEwubScW99ro2yDgl020
bVr0YAtn+OTniokkWaFVVNtQejjr70/tyOpySbvxQWaVMb/wxMb6qDzx11qGvsj0
6u7x8CqbAI/ypqbp4qobZvBGymqSGt7siSgXV5lwpYbyQIHufxhAOfMHpcCy5LiN
onPQiIFksVfFplUAbeyGaTLOvZHHXShUH9WAqhTDAYcycbUoO9h7OHD/spB4zTTF
D0wPi3jPQfuYDAMz5mchI/r+nwIDAQABo1AwTjAdBgNVHQ4EFgQUQ05CgXBgkBai
maMqUV1uVjEMFkUwHwYDVR0jBBgwFoAUQ05CgXBgkBaimaMqUV1uVjEMFkUwDAYD
VR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEABPj0XtnNGdEQ/kR8a6nLnZpX
bKDDHDsDOEBJOxxkZAJnkL0hZZHR/zR8DTYOP82xIOwIt9Gqzn3gsag6EGjB2FBv
0m+7Hu4dt3AhfvfBC6bf0+9IcWy0BpVPdpp6zfP2CszBKmb+VTy9cE5s3+A+Ukgg
jpL3Jttu8oYZeE1jzvMBPL+kXCJ7lxcHehhkzLhsBqhbH4eWsaH4DmfN+QcKT13b
OpE0EWAq4LWnCzf/dBhobR7poEF1wUNb7ORL7bBLXW8vfrcAXNqU5L2vhyEqQ8Rx
E2FZUFkR4B6Pnf/b3YSo8Egp4oKP8f0xX8BSnXZjr+ROjVLryXVqgOL6NNvwdg==
-----END CERTIFICATE-----
`)
	duration = time.Hour * 24 * 365
)
