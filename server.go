package main

import (
	"crypto"
	"crypto/x509"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlidp"
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
)

const (
	templatesGlob = "./templates/*.tmpl"
	metadataRoute = "/metadata"
	ssoRoute      = "/sso"
)

type IdpServer struct {
	router *gin.Engine
	idp    *saml.IdentityProvider
	store  samlidp.Store
}

func (s *IdpServer) Run() error {
	return s.router.Run()
}

type ServerOptions struct {
	BaseUrl     url.URL
	Key         crypto.PrivateKey
	Certificate *x509.Certificate
}

func NewServer(o ServerOptions) *IdpServer {
	metadataUrl := o.BaseUrl
	metadataUrl.Path += metadataRoute

	ssoUrl := o.BaseUrl
	ssoUrl.Path += ssoRoute

	idp := &saml.IdentityProvider{
		Logger:      log.Default(),
		Key:         o.Key,
		Certificate: o.Certificate,
		MetadataURL: metadataUrl,
		SSOURL:      ssoUrl,
	}

	router := gin.Default()
	router.LoadHTMLGlob(templatesGlob)

	router.GET(metadataRoute, func(c *gin.Context) {
		metadata := idp.Metadata()
		c.XML(200, metadata)
	})

	router.GET(ssoRoute, func(c *gin.Context) {
		idp.ServeSSO(c.Writer, c.Request)
	})

	router.POST(ssoRoute, func(c *gin.Context) {
		idp.ServeSSO(c.Writer, c.Request)
	})

	server := &IdpServer{
		router: router,
		idp:    idp,
		store:  &samlidp.MemoryStore{},
	}

	return server
}
