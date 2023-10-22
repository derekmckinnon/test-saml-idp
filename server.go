package idp

import (
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlidp"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
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
		log.Fatal().Err(err).Msg("cannot parse host URL")
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
		Logger:      &zerologAdapter{},
		Certificate: options.Certificate,
		Key:         options.Key,
		MetadataURL: metadataUrl,
		SSOURL:      ssoUrl,
	}

	return idp
}

func buildRouter(host url.URL, idp *saml.IdentityProvider, store *Store) *gin.Engine {
	basePath := getBasePath(host)

	router := gin.New()
	router.Use(logger.SetLogger(
		logger.WithSkipPath([]string{healthRoute, basePath + healthRoute}),
	))
	router.Use(gin.Recovery())

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

		log.Info().Str("username", user.Username).Msg("initialized user")
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

		log.Info().Str("serviceProvider", service.EntityId).Msg("Initialized service provider")
	}

	return nil
}

func (s *Server) Run() error {
	return s.router.Run()
}
