package idpserver

import (
	"crypto"
	"crypto/x509"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/logger"
	"github.com/crewjam/saml/samlidp"
	"log"
	"net/url"
)

type IdpServer struct {
	logger logger.Interface
	store  samlidp.Store
	idp    saml.IdentityProvider
}

func (s *IdpServer) GetMetadata() *saml.EntityDescriptor {
	return s.idp.Metadata()
}

type Options struct {
	BaseUrl     url.URL
	Key         crypto.PrivateKey
	Certificate *x509.Certificate
	Logger      logger.Interface
	Store       samlidp.Store
}

func New(options Options) (server *IdpServer) {
	loggerImpl := options.Logger
	if loggerImpl == nil {
		loggerImpl = log.Default()
	}

	storeImpl := options.Store
	if storeImpl == nil {
		storeImpl = &samlidp.MemoryStore{}
	}

	metadataUrl := options.BaseUrl
	metadataUrl.Path += "/metadata"

	ssoUrl := options.BaseUrl
	ssoUrl.Path += "/sso"

	server = &IdpServer{
		logger: loggerImpl,
		store:  storeImpl,
		idp: saml.IdentityProvider{
			Logger:      loggerImpl,
			Key:         options.Key,
			Certificate: options.Certificate,
			MetadataURL: metadataUrl,
			SSOURL:      ssoUrl,
		},
	}

	server.idp.ServiceProviderProvider = server
	server.idp.SessionProvider = server
	server.idp.AssertionMaker = server

	return
}
