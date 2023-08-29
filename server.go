package idp

import (
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlidp"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/url"
)

const (
	templatesGlob        = "templates/*.html"
	metadataRoute        = "/metadata"
	ssoRoute             = "/sso"
	healthRoute          = "/health"
	defaultSessionMaxAge = 60 // 1 hour
)

type Server struct {
	config *Config
	idp    *saml.IdentityProvider
	router *gin.Engine
	Store  *Store
}

func New(options ServerOptions) *Server {
	config := options.Config
	if config.SessionMaxAge == 0 {
		config.SessionMaxAge = defaultSessionMaxAge
	}

	host, err := url.Parse(config.Host)
	if err != nil {
		log.Fatalf("cannot parse host URL: %v", err)
	}

	idp := buildIdp(*host, options)

	store := &Store{}

	router := buildRouter(*host, idp, store)

	server := &Server{
		config: config,
		idp:    idp,
		router: router,
		Store:  store,
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

func buildRouter(host url.URL, idp *saml.IdentityProvider, store *Store) *gin.Engine {
	basePath := getBasePath(host)

	loggerConfig := gin.LoggerConfig{
		SkipPaths: []string{
			healthRoute,
			basePath + healthRoute,
		},
	}

	router := gin.New()
	router.Use(gin.LoggerWithConfig(loggerConfig), gin.Recovery())

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

	router.GET(basePath+"/users", func(c *gin.Context) {
		users, err := store.GetUsers()
		if err != nil {
			users = []*samlidp.User{}
		}

		c.JSON(200, users)
	})

	router.GET(basePath+"/users/create", func(c *gin.Context) {
		_, success := c.GetQuery("success")

		c.HTML(200, "create-user.html", gin.H{
			"Title":   "Create User",
			"Success": success,
		})
	})

	router.POST(basePath+"/users/create", func(c *gin.Context) {
		var errors []string

		username, ok := c.GetPostForm("username")
		if !ok || len(username) == 0 {
			errors = append(errors, "Username is required")
		}

		email, ok := c.GetPostForm("email")
		if !ok || len(email) == 0 {
			errors = append(errors, "Email is required")
		}

		firstName, ok := c.GetPostForm("first_name")
		if !ok || len(firstName) == 0 {
			errors = append(errors, "First Name is required")
		}

		lastName, ok := c.GetPostForm("last_name")
		if !ok || len(lastName) == 0 {
			errors = append(errors, "Username is required")
		}

		password, ok := c.GetPostForm("password")
		if !ok || len(password) == 0 {
			errors = append(errors, "Password is required")
		}

		if len(errors) > 0 {
			c.HTML(200, "create-user.html", gin.H{
				"Title":     "Create User",
				"Errors":    errors,
				"Username":  username,
				"Email":     email,
				"FirstName": firstName,
				"LastName":  lastName,
			})
			return
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		err := store.AddUser(&samlidp.User{
			Name:           username,
			Email:          email,
			HashedPassword: hashedPassword,
			GivenName:      firstName,
			Surname:        lastName,
		})

		if err != nil {
			errors = append(errors, err.Error())
		}

		c.Redirect(http.StatusFound, basePath+"/users/create?success")
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
			Groups:         user.Groups,
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
