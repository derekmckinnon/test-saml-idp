package idp

import (
	"crypto"
	"crypto/x509"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlidp"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/url"
)

const (
	templatesGlob = "./templates/*.tmpl"
	metadataRoute = "/metadata"
	ssoRoute      = "/sso"
)

type Server struct {
	router *gin.Engine
	idp    *saml.IdentityProvider
	Store  *Store
}

func (s *Server) LoadUsers(users []User) error {
	for _, user := range users {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		err := s.Store.AddUser(&samlidp.User{
			Name:           user.Username,
			Email:          user.Email,
			HashedPassword: hashedPassword,
			GivenName:      user.FirstName,
			Surname:        user.LastName,
		})

		if err != nil {
			return err
		}

		log.Printf("Initialized User: %s\n", user.Username)
	}

	return nil
}

func (s *Server) LoadServices(services []Service) error {
	for _, service := range services {
		acs := saml.IndexedEndpoint{
			Binding:  saml.HTTPPostBinding,
			Location: service.AssertionConsumerService,
		}

		descriptor := saml.SPSSODescriptor{
			AssertionConsumerServices: []saml.IndexedEndpoint{acs},
		}

		err := s.Store.AddServiceProvider(&samlidp.Service{
			Name: service.EntityId,
			Metadata: saml.EntityDescriptor{
				EntityID:         service.EntityId,
				SPSSODescriptors: []saml.SPSSODescriptor{descriptor},
			},
		})

		if err != nil {
			return err
		}

		log.Printf("Initialized service provider: %s\n", service.EntityId)
	}

	return nil
}

func (s *Server) Run() error {
	return s.router.Run()
}

type ServerOptions struct {
	BaseUrl     url.URL
	Key         crypto.PrivateKey
	Certificate *x509.Certificate
}

func New(o ServerOptions) *Server {
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

	store := &Store{}

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

	router.GET("/health", func(c *gin.Context) {
		c.String(200, "Healthy")
	})

	router.GET("/", func(c *gin.Context) {
		handleErr := func(err error) {
			c.JSON(500, gin.H{
				"error": err,
			})
		}

		users, err := store.GetUsers()
		if err != nil {
			handleErr(err)
			return
		}

		services, err := store.GetServiceProviders()
		if err != nil {
			handleErr(err)
			return
		}

		sessions, err := store.GetSessions()
		if err != nil {
			handleErr(err)
			return
		}

		c.JSON(200, gin.H{
			"users":    users,
			"services": services,
			"sessions": sessions,
		})
	})

	server := &Server{
		router: router,
		idp:    idp,
		Store:  store,
	}

	idp.ServiceProviderProvider = server
	idp.SessionProvider = server

	return server
}
