package idp

import (
	"encoding/base64"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlidp"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"time"
)

var sessionMaxAge = time.Hour * 24 * 14  // 14 days

func init() {
	saml.MaxIssueDelay = sessionMaxAge
}

func (s *Server) GetSession(w http.ResponseWriter, r *http.Request, req *saml.IdpAuthnRequest) *saml.Session {
	if r.Method == http.MethodPost && r.PostForm.Get("username") != "" {
		user, err := s.Store.GetUser(r.PostForm.Get("username"))
		if err != nil {
			s.serveLoginPage(w, r, req, "Invalid username or password")
			return nil
		}

		if bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(r.PostForm.Get("password"))) != nil {
			s.serveLoginPage(w, r, req, "Invalid username or password")
			return nil
		}

		now := saml.TimeNow()

		session := &saml.Session{
			ID:                    uuid.NewString(),
			NameID:                user.Email,
			CreateTime:            now,
			ExpireTime:            now.Add(sessionMaxAge),
			Index:                 uuid.NewString(),
			UserName:              user.Name,
			Groups:                user.Groups[:],
			UserEmail:             user.Email,
			UserCommonName:        user.CommonName,
			UserSurname:           user.Surname,
			UserGivenName:         user.GivenName,
			UserScopedAffiliation: user.ScopedAffiliation,
		}

		if s.Store.AddSession(session) != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return nil
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    session.ID,
			MaxAge:   int(sessionMaxAge.Seconds()),
			HttpOnly: true,
			Secure:   r.URL.Scheme == "https",
			Path:     "/",
		})

		return session
	}

	if cookie, err := r.Cookie("session"); err == nil {
		session, err := s.Store.GetSession(cookie.Value)

		if err != nil {
			if err == samlidp.ErrNotFound {
				s.serveLoginPage(w, r, req, "")
				return nil
			}

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return nil
		}

		if saml.TimeNow().After(session.ExpireTime) {
			s.serveLoginPage(w, r, req, "")
			return nil
		}

		return session
	}

	s.serveLoginPage(w, r, req, "")
	return nil
}

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

	options := s.config.LoginPageOptions

	if options.Title != "" {
		data.Title = options.Title
	}

	if options.Description != "" {
		data.Description = renderMarkdown(options.Description)
	}

	if options.DumpUsers {
		data.Users = s.config.Users
	}

	render := s.router.HTMLRender.Instance("login.tmpl", data)

	err := render.Render(w)
	if err != nil {
		panic(err)
	}
}
