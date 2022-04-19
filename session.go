package main

import (
	"encoding/base64"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlidp"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var sessionMaxAge = time.Hour

func (s *IdpServer) GetSession(w http.ResponseWriter, r *http.Request, req *saml.IdpAuthnRequest) *saml.Session {
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

func (s *IdpServer) serveLoginPage(w http.ResponseWriter, r *http.Request, req *saml.IdpAuthnRequest, toast string) {
	data := struct {
		Username    string
		Toast       string
		Url         string
		SamlRequest string
		RelayState  string
	}{
		Username:    r.PostForm.Get("username"),
		Toast:       toast,
		Url:         req.IDP.SSOURL.String(),
		SamlRequest: base64.StdEncoding.EncodeToString(req.RequestBuffer),
		RelayState:  req.RelayState,
	}

	render := s.router.HTMLRender.Instance("login.tmpl", data)

	err := render.Render(w)
	if err != nil {
		panic(err)
	}
}
