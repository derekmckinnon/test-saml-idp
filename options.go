package idp

import (
	"crypto"
	"crypto/x509"
	"net/url"
	"strings"
)

type ServerOptions struct {
	BaseUrl     url.URL
	Key         crypto.PrivateKey
	Certificate *x509.Certificate
}

func (o *ServerOptions) getBasePath() string {
	basePath := o.BaseUrl.Path

	if basePath == "" || basePath == "/" {
		return "/"
	}

	return strings.TrimSuffix(basePath, "/")
}
