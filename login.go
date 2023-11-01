package idp

import (
	"encoding/base64"
	"github.com/crewjam/saml"
	"html/template"
	"net/http"
)

type LoginPageData struct {
	Title       string
	Description template.HTML
	Users       []User
	Toast       string
	Username    string
	Url         string
	SamlRequest string
	RelayState  string
}

func (s *Server) serveLoginPage(w http.ResponseWriter, r *http.Request, req *saml.IdpAuthnRequest, toast string) {
	data := LoginPageData{
		Title:       "Login",
		Description: "",
		Toast:       toast,
		Username:    r.PostForm.Get("username"),
		Url:         req.IDP.SSOURL.String(),
		SamlRequest: base64.StdEncoding.EncodeToString(req.RequestBuffer),
		RelayState:  req.RelayState,
	}

	options := s.config.LoginPage

	if options.Title != "" {
		data.Title = options.Title
	}

	if options.Description != "" {
		data.Description = renderMarkdown(options.Description)
	}

	if options.DumpUsers {
		data.Users = s.config.Users
	}

	render := s.router.HTMLRender.Instance("login.html", data)

	err := render.Render(w)
	if err != nil {
		panic(err)
	}
}
