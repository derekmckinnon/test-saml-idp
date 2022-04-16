package idpserver

import (
	"github.com/crewjam/saml"
	"net/http"
)

func (s *IdpServer) GetSession(w http.ResponseWriter, r *http.Request, req *saml.IdpAuthnRequest) *saml.Session {
	return nil
}
