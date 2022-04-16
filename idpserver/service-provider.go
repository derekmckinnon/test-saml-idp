package idpserver

import (
	"github.com/crewjam/saml"
	"net/http"
	"os"
)

func (s *IdpServer) GetServiceProvider(r *http.Request, id string) (*saml.EntityDescriptor, error) {
	return nil, os.ErrNotExist
}
