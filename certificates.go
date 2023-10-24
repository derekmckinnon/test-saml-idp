package idp

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"time"
)

func GenerateDevelopmentCertificateAndKey() (*x509.Certificate, *rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
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

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, nil, err
	}

	return cert, key, nil
}

func LoadCertificatePem(path string) (*x509.Certificate, error) {
	block, err := loadPem(path)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(block.Bytes)
}

func LoadPrivateKeyPem(path string) (*rsa.PrivateKey, error) {
	block, err := loadPem(path)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("pem file does not contain an RSA private key")
	}

	return rsaKey, err
}

func loadPem(path string) (*pem.Block, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)

	return block, nil
}
