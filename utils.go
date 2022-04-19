package main

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func getEnvOrDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)

	if ok {
		return value
	}

	return defaultValue
}

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

	return x509.ParsePKCS8PrivateKey(block.Bytes)
}

func loadCertificate() (*x509.Certificate, error) {
	block, err := loadPemFile("./idp_cert.pem")
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(block.Bytes)
}
