package main

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func loadPemFile(path string) (*pem.Block, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)

	return block, nil
}

func loadPrivateKey() (crypto.PrivateKey, error) {
	block, err := loadPemFile("./idp_key.pem")
	if err != nil {
		return nil, err
	}

	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	return key, nil
}

func loadCertificate() (*x509.Certificate, error) {
	block, err := loadPemFile("./idp_cert.pem")
	if err != nil {
		return nil, err
	}

	certificate, _ := x509.ParseCertificate(block.Bytes)

	return certificate, nil
}
