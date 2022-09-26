package idp

import (
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
	healthRoute   = "/health"
)

type Server struct {
	config *Config
	idp    *saml.IdentityProvider
	router *gin.Engine
	Store  *Store
}

func New(options ServerOptions) *Server {
	config := options.Config

	host, err := url.Parse(config.Host)
	if err != nil {
		log.Fatalf("cannot parse host URL: %v", err)
	}

	idp := buildIdp(*host, options)

	router := buildRouter(*host, idp)

	server := &Server{
		config: config,
		idp:    idp,
		router: router,
		Store:  &Store{},
	}

	idp.ServiceProviderProvider = server
	idp.SessionProvider = server

	return server
}

func buildIdp(host url.URL, options ServerOptions) *saml.IdentityProvider {
	metadataUrl := host
	metadataUrl.Path += metadataRoute

	ssoUrl := host
	ssoUrl.Path += ssoRoute

	idp := &saml.IdentityProvider{
		Logger:      log.Default(),
		Certificate: options.Certificate,
		Key:         options.Key,
		MetadataURL: metadataUrl,
		SSOURL:      ssoUrl,
	}

	return idp
}

func buildRouter(host url.URL, idp *saml.IdentityProvider) *gin.Engine {
	basePath := getBasePath(host)

	router := gin.Default()
	router.LoadHTMLGlob(templatesGlob)

	router.GET(basePath+metadataRoute, func(c *gin.Context) {
		metadata := idp.Metadata()
		c.XML(200, metadata)
	})

	router.GET(basePath+ssoRoute, func(c *gin.Context) {
		idp.ServeSSO(c.Writer, c.Request)
	})

	router.POST(basePath+ssoRoute, func(c *gin.Context) {
		idp.ServeSSO(c.Writer, c.Request)
	})

	router.GET(basePath+healthRoute, func(c *gin.Context) {
		c.String(200, "Healthy")
	})

	return router
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
