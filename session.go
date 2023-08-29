package idp

import (
	"errors"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlidp"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

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
		expires := now.Add(time.Duration(s.config.SessionMaxAge) * time.Minute)

		session := &saml.Session{
			ID:                    uuid.NewString(),
			NameID:                user.Email,
			CreateTime:            now,
			ExpireTime:            expires,
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
			Expires:  expires,
			HttpOnly: true,
			Secure:   r.URL.Scheme == "https",
			Path:     "/",
		})

		return session
	}

	if cookie, err := r.Cookie("session"); err == nil {
		session, err := s.Store.GetSession(cookie.Value)

		if err != nil {
			if errors.Is(err, samlidp.ErrNotFound) {
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
