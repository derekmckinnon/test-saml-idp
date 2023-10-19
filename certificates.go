package idp

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

func GenerateDevelopmentCertificate() (*x509.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test SAML IdP"},
		},
		NotBefore: now,
		NotAfter:  now.Add(time.Hour * 24 * 365 * 3),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
