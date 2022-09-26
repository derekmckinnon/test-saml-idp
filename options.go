package idp

import (
	"crypto"
	"crypto/x509"
)

type ServerOptions struct {
	Config      *Config
	Key         crypto.PrivateKey
	Certificate *x509.Certificate
}
